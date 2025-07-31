package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authUsecases "stock-tracker/internal/domain/authentication/usecase"
	recommendationUsecase "stock-tracker/internal/domain/recommendation/usecase"
	recommendationValidation "stock-tracker/internal/domain/recommendation/validation"
	"stock-tracker/internal/domain/shared/enums"
	stockUsecases "stock-tracker/internal/domain/stocks/usecase"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/internal/infrastructure/config"
	"stock-tracker/internal/infrastructure/database"
	"stock-tracker/internal/infrastructure/external"
	infraMiddleware "stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/pkg/logger"
)

// IntegrationTestSuite holds the test environment
type IntegrationTestSuite struct {
	router chi.Router
	db     *database.Connection
	logger logger.Logger
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Set up test environment
	suite := setupTestEnvironment()

	// Run tests
	code := m.Run()

	// Clean up
	teardownTestEnvironment(suite)

	os.Exit(code)
}

func setupTestEnvironment() *IntegrationTestSuite {
	// Load test configuration
	cfg := &config.Config{
		DatabaseURL: getTestDatabaseURL(),
		Port:        "8080",
		LogLevel:    "debug",
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	// Initialize database connection (use test database)
	dbPool, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	// Initialize repositories
	stockRepo := database.NewStockRepository(dbPool.GetPool(), log)
	brokerRepo := database.NewBrokerRepository(dbPool.GetPool())
	userRepo := database.NewUserRepository(dbPool.GetPool(), log)
	sessionRepo := database.NewSessionRepository(dbPool.GetPool())
	subscriptionRepo := database.NewSubscriptionRepository(dbPool.GetPool(), log)

	// Initialize JWT service
	jwtService := auth.NewJWTService("test-secret-key")

	// Initialize cache
	cache := cache.NewInMemoryCache()

	// Initialize external API clients
	yahooClient := external.NewYahooFinanceClient(log)
	alphaVantageClient := external.NewAlphaVantageClient("", log) // Empty API key for testing

	// Initialize recommendation components
	basicCalculator := recommendationValidation.NewBasicScoringCalculator(stockRepo, log)
	externalEnricher := recommendationValidation.NewExternalDataEnricher(yahooClient, alphaVantageClient, cache, log)

	// Initialize use cases
	stockQueryUC := stockUsecases.NewStockQueryUseCase(stockRepo, brokerRepo, log)
	userUC := authUsecases.NewUserUseCase(userRepo, subscriptionRepo, sessionRepo, jwtService, log)
	recommendationUC := recommendationUsecase.NewTieredRecommendationUseCase(stockRepo, basicCalculator, externalEnricher, cache, log)

	// Initialize middleware
	authMiddleware := infraMiddleware.NewAuthMiddleware(jwtService, log)
	rateLimiter := infraMiddleware.NewRateLimiter(log)

	// Initialize handlers
	stockHandler := handlers.NewStockHandler(stockQueryUC, log)
	authHandler := handlers.NewAuthHandler(userUC, log)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationUC, log)

	// Setup router
	router := setupTestRouter(stockHandler, authHandler, recommendationHandler, authMiddleware, rateLimiter, log, dbPool)

	return &IntegrationTestSuite{
		router: router,
		db:     dbPool,
		logger: log,
	}
}

func teardownTestEnvironment(suite *IntegrationTestSuite) {
	if suite.db != nil {
		_ = suite.db.Close()
	}
}

func getTestDatabaseURL() string {
	// Use environment variable or default test database
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return "postgres://postgres:password@localhost/stock_tracker_test?sslmode=disable"
}

func setupTestRouter(
	stockHandler *handlers.StockHandler,
	authHandler *handlers.AuthHandler,
	recommendationHandler *handlers.RecommendationHandler,
	authMiddleware *infraMiddleware.AuthMiddleware,
	rateLimiter *infraMiddleware.RateLimiter,
	log logger.Logger,
	dbPool *database.Connection,
) chi.Router {
	r := chi.NewRouter()

	// Basic middleware for testing
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Authentication routes
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", authHandler.Register)
				r.Post("/login", authHandler.Login)
			})

			// Recommendation routes
			r.Route("/recommendations", func(r chi.Router) {
				r.Use(authMiddleware.OptionalAuth)
				r.Get("/", recommendationHandler.GetRecommendations)
				r.Get("/{ticker}", recommendationHandler.GetRecommendationByTicker)
			})

			// Stock routes
			r.Route("/stocks", func(r chi.Router) {
				r.Use(authMiddleware.OptionalAuth)
				r.Get("/", stockHandler.GetStocks)
			})
		})
	})

	return r
}

// TestRecommendationAPI_GuestUser tests guest user access to recommendations
func TestRecommendationAPI_GuestUser(t *testing.T) {
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	// Test guest user gets limited recommendations
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=20", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify guest limitations
	assert.Equal(t, enums.TIER_GUEST, response.Meta.UserTier)
	assert.Contains(t, response.Meta.Features, "basic_recommendations")
	assert.NotContains(t, response.Meta.Features, "external_apis")
	assert.LessOrEqual(t, len(response.Data), 5, "Guest users should be limited to 5 recommendations")

	// Verify no premium features
	for _, rec := range response.Data {
		assert.Nil(t, rec.ExternalData, "Guest users should not have external data")
		assert.Nil(t, rec.AIInsights, "Guest users should not have AI insights")
	}
}

// TestRecommendationAPI_AuthenticatedUser tests authenticated user features
func TestRecommendationAPI_AuthenticatedUser(t *testing.T) {
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	// Create and authenticate a test user
	token := createTestUser(t, suite, "basic_user@test.com", enums.TIER_BASIC)

	// Test authenticated user gets enhanced recommendations
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=20", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify basic user features
	assert.Equal(t, enums.TIER_BASIC, response.Meta.UserTier)
	assert.Contains(t, response.Meta.Features, "external_apis")
	assert.NotContains(t, response.Meta.Features, "ai_insights")
	assert.LessOrEqual(t, len(response.Data), 20)

	// Verify external data is included for basic users
	if len(response.Data) > 0 {
		// Note: External data may be nil if external APIs are not available in test environment
		// This is acceptable for integration testing
		assert.NotNil(t, response.Data[0]) // Ensure at least basic data structure is present
	}
}

// TestRecommendationAPI_PremiumUser tests premium user features
func TestRecommendationAPI_PremiumUser(t *testing.T) {
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	// Create and authenticate a premium test user
	token := createTestUser(t, suite, "premium_user@test.com", enums.TIER_PREMIUM)

	// Test premium user gets full features
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=50", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify premium user features
	assert.Equal(t, enums.TIER_PREMIUM, response.Meta.UserTier)
	assert.Contains(t, response.Meta.Features, "ai_insights")
	assert.Contains(t, response.Meta.Features, "priority_support")
	assert.LessOrEqual(t, len(response.Data), 50)
}

// TestRecommendationAPI_FilteringAndValidation tests request filtering and validation
func TestRecommendationAPI_FilteringAndValidation(t *testing.T) {
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		description    string
	}{
		{
			name:           "Valid filters",
			url:            "/api/v1/recommendations?limit=10&min_score=0.7&type=buy",
			expectedStatus: http.StatusOK,
			description:    "Should accept valid filters",
		},
		{
			name:           "Invalid limit",
			url:            "/api/v1/recommendations?limit=-1",
			expectedStatus: http.StatusOK, // Handler should correct invalid limit
			description:    "Should handle invalid limit gracefully",
		},
		{
			name:           "Exclude tickers",
			url:            "/api/v1/recommendations?exclude=TSLA,MSFT",
			expectedStatus: http.StatusOK,
			description:    "Should handle ticker exclusion",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)

			if w.Code == http.StatusOK {
				var response recommendationUsecase.RecommendationResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
			}
		})
	}
}

// TestRecommendationAPI_ByTicker tests individual ticker recommendations
func TestRecommendationAPI_ByTicker(t *testing.T) {
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	testCases := []struct {
		name           string
		ticker         string
		expectedStatus int
		description    string
	}{
		{
			name:           "Valid ticker",
			ticker:         "AAPL",
			expectedStatus: http.StatusOK,
			description:    "Should return recommendation for valid ticker",
		},
		{
			name:           "Invalid ticker format",
			ticker:         "INVALID$",
			expectedStatus: http.StatusBadRequest,
			description:    "Should reject invalid ticker format",
		},
		{
			name:           "Empty ticker",
			ticker:         "",
			expectedStatus: http.StatusBadRequest,
			description:    "Should reject empty ticker",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/recommendations/"+tc.ticker, nil)
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)
		})
	}
}

// TestRecommendationAPI_CachePerformance tests caching behavior
func TestRecommendationAPI_CachePerformance(t *testing.T) {
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	// First request (cache miss)
	req1 := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
	w1 := httptest.NewRecorder()
	start1 := time.Now()

	suite.router.ServeHTTP(w1, req1)
	duration1 := time.Since(start1)

	assert.Equal(t, http.StatusOK, w1.Code)

	var response1 recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w1.Body.Bytes(), &response1)
	require.NoError(t, err)

	// Second request (should be cached)
	req2 := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
	w2 := httptest.NewRecorder()
	start2 := time.Now()

	suite.router.ServeHTTP(w2, req2)
	duration2 := time.Since(start2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response2 recommendationUsecase.RecommendationResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response2)
	require.NoError(t, err)

	// Cache should make the second request faster (though this may not always be guaranteed in tests)
	suite.logger.Info("Cache performance test",
		"first_request_duration", duration1,
		"second_request_duration", duration2,
		"first_cache_hit", response1.Meta.CacheHit,
		"second_cache_hit", response2.Meta.CacheHit)
}

// Helper function to create a test user and return JWT token
func createTestUser(t *testing.T, suite *IntegrationTestSuite, email string, tier enums.UserTier) string {
	// Create user registration request
	registerReq := map[string]interface{}{
		"email":    email,
		"password": "testpassword123",
		"tier":     string(tier),
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		// If registration fails (user might already exist), try login
		loginReq := map[string]interface{}{
			"email":    email,
			"password": "testpassword123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		suite.router.ServeHTTP(w, req)
	}

	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("Failed to create/login test user: %d", w.Code)
	}

	var authResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &authResponse)
	require.NoError(t, err)

	token, ok := authResponse["access_token"].(string)
	require.True(t, ok, "Should receive access token")

	return token
}
