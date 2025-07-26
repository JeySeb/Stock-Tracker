package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/application/usecases"
	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/internal/presentation/handlers"
)

// MockTieredRecommendationUseCase mock implementation
type MockTieredRecommendationUseCase struct {
	mock.Mock
}

func (m *MockTieredRecommendationUseCase) GetRecommendations(ctx context.Context, request usecases.RecommendationRequest) (*usecases.RecommendationResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.RecommendationResponse), args.Error(1)
}

func (m *MockTieredRecommendationUseCase) GetRecommendationByTicker(ctx context.Context, ticker string, userTier enums.UserTier) (*model.AggregatedRecommendation, error) {
	args := m.Called(ctx, ticker, userTier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AggregatedRecommendation), args.Error(1)
}

// MockLogger implementation
type MockLogger struct{}

func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {}

func TestRecommendationHandler_GetRecommendations_GuestUser(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data
	expectedResponse := &usecases.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{
				Ticker:      "AAPL",
				CompanyName: "Apple Inc.",
				BasicScore:  0.8,
				Tier:        enums.RECOMMENDATION_TIER_BASIC,
			},
		},
		Meta: usecases.RecommendationMeta{
			Count:          1,
			UserTier:       enums.TIER_GUEST,
			Features:       []string{"basic_recommendations", "market_analytics"},
			CacheHit:       false,
			GenerationTime: 100 * time.Millisecond,
		},
	}

	// Mock expectations
	mockUseCase.On("GetRecommendations", mock.Anything, mock.MatchedBy(func(req usecases.RecommendationRequest) bool {
		return req.UserTier == enums.TIER_GUEST && req.Limit == 10 // Default limit for guests
	})).Return(expectedResponse, nil)

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=5", nil)
	// Context without user tier should default to guest
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Meta.Count)
	assert.Equal(t, enums.TIER_GUEST, response.Meta.UserTier)
	assert.Contains(t, response.Meta.Features, "basic_recommendations")
	assert.NotContains(t, response.Meta.Features, "real_time_data")

	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_GetRecommendations_BasicUser(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data
	expectedResponse := &usecases.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{
				Ticker:      "AAPL",
				CompanyName: "Apple Inc.",
				BasicScore:  0.85,
				Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
				ExternalData: &model.ExternalStockData{
					CurrentPrice: 170.0,
				},
			},
		},
		Meta: usecases.RecommendationMeta{
			Count:    1,
			UserTier: enums.TIER_BASIC,
			Features: []string{"basic_recommendations", "market_analytics", "real_time_data", "external_apis", "advanced_analytics"},
			CacheHit: false,
		},
	}

	// Mock expectations
	mockUseCase.On("GetRecommendations", mock.Anything, mock.MatchedBy(func(req usecases.RecommendationRequest) bool {
		return req.UserTier == enums.TIER_BASIC && req.Limit == 20
	})).Return(expectedResponse, nil)

	// Create request with basic user context
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=20", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, enums.TIER_BASIC, response.Meta.UserTier)
	assert.Contains(t, response.Meta.Features, "real_time_data")
	assert.NotNil(t, response.Data[0].ExternalData)

	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_GetRecommendations_WithFilters(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data
	expectedResponse := &usecases.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{
				Ticker:             "AAPL",
				CompanyName:        "Apple Inc.",
				BasicScore:         0.8,
				RecommendationType: enums.RECOMMENDATION_TYPE_BUY,
			},
		},
		Meta: usecases.RecommendationMeta{
			Count:    1,
			UserTier: enums.TIER_BASIC,
		},
	}

	// Mock expectations - verify filters are passed correctly
	mockUseCase.On("GetRecommendations", mock.Anything, mock.MatchedBy(func(req usecases.RecommendationRequest) bool {
		return req.UserTier == enums.TIER_BASIC &&
			req.Limit == 15 &&
			req.Filters.MinScore != nil && *req.Filters.MinScore == 0.7 &&
			req.Filters.RecommendationType != nil && *req.Filters.RecommendationType == enums.RECOMMENDATION_TYPE_BUY &&
			len(req.Filters.ExcludeTickers) == 2 &&
			req.Filters.ExcludeTickers[0] == "TSLA" &&
			req.Filters.ExcludeTickers[1] == "MSFT"
	})).Return(expectedResponse, nil)

	// Create request with filters
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=15&min_score=0.7&type=buy&exclude=TSLA,MSFT", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_GetRecommendations_WithRateLimit(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data
	expectedResponse := &usecases.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{Ticker: "AAPL", BasicScore: 0.8},
		},
		Meta: usecases.RecommendationMeta{
			Count:    1,
			UserTier: enums.TIER_BASIC,
		},
	}

	mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Create request with rate limit info in context
	req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	ctx = context.WithValue(ctx, middleware.RateLimitRemainingKey, 450)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 450, response.Meta.RateLimitRemaining)

	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_GetRecommendationByTicker_Success(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data
	expectedRecommendation := &model.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.85,
		Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
		ExternalData: &model.ExternalStockData{
			CurrentPrice: 170.0,
		},
	}

	// Mock expectations
	mockUseCase.On("GetRecommendationByTicker", mock.Anything, "AAPL", enums.TIER_BASIC).Return(expectedRecommendation, nil)

	// Create request with URL parameter
	req := httptest.NewRequest("GET", "/api/v1/recommendations/AAPL", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ticker", "AAPL")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendationByTicker(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["meta"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "AAPL", data["ticker"])
	assert.Equal(t, "Apple Inc.", data["company_name"])

	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_GetRecommendationByTicker_InvalidTicker(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	tests := []struct {
		name   string
		ticker string
	}{
		{"Empty ticker", ""},
		{"Too long ticker", "VERYLONGTICKER"},
		{"Invalid characters", "A@PL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/recommendations/"+tt.ticker, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ticker", tt.ticker)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.GetRecommendationByTicker(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], "Ticker")
		})
	}
}

func TestRecommendationHandler_GetRecommendationByTicker_NotFound(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Mock expectations - return nil (not found)
	mockUseCase.On("GetRecommendationByTicker", mock.Anything, "NOTFOUND", enums.TIER_GUEST).Return(nil, nil)

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/recommendations/NOTFOUND", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ticker", "NOTFOUND")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendationByTicker(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "not found")

	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_GetRecommendationPreview_RequiresAuth(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Create request without auth (guest user)
	req := httptest.NewRequest("GET", "/api/v1/recommendations/preview/AAPL", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ticker", "AAPL")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendationPreview(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Authentication required")
}

func TestRecommendationHandler_GetRecommendationPreview_Success(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data
	basicRecommendation := &model.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.8,
		Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
	}

	premiumRecommendation := &model.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.9,
		Tier:        enums.RECOMMENDATION_TIER_PREMIUM,
		ExternalData: &model.ExternalStockData{
			CurrentPrice: 170.0,
		},
		// AI insights would be here in Phase 6
	}

	// Mock expectations
	mockUseCase.On("GetRecommendationByTicker", mock.Anything, "AAPL", enums.TIER_BASIC).Return(basicRecommendation, nil)
	mockUseCase.On("GetRecommendationByTicker", mock.Anything, "AAPL", enums.TIER_PREMIUM).Return(premiumRecommendation, nil)

	// Create request with basic user auth
	req := httptest.NewRequest("GET", "/api/v1/recommendations/preview/AAPL", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ticker", "AAPL")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendationPreview(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["current_tier"])
	assert.NotNil(t, response["premium_preview"])
	assert.NotNil(t, response["upgrade_benefits"])

	mockUseCase.AssertExpectations(t)
}

func TestRecommendationHandler_ParameterParsing(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	tests := []struct {
		name       string
		url        string
		setupMocks func()
		validate   func(req usecases.RecommendationRequest) bool
	}{
		{
			name: "Parse min_score filter",
			url:  "/api/v1/recommendations?min_score=0.75",
			setupMocks: func() {
				mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(&usecases.RecommendationResponse{}, nil)
			},
			validate: func(req usecases.RecommendationRequest) bool {
				return req.Filters.MinScore != nil && *req.Filters.MinScore == 0.75
			},
		},
		{
			name: "Parse recommendation type filter",
			url:  "/api/v1/recommendations?type=strong_buy",
			setupMocks: func() {
				mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(&usecases.RecommendationResponse{}, nil)
			},
			validate: func(req usecases.RecommendationRequest) bool {
				return req.Filters.RecommendationType != nil && *req.Filters.RecommendationType == enums.RECOMMENDATION_TYPE_STRONG_BUY
			},
		},
		{
			name: "Parse exclude tickers",
			url:  "/api/v1/recommendations?exclude=TSLA,MSFT,GOOGL",
			setupMocks: func() {
				mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(&usecases.RecommendationResponse{}, nil)
			},
			validate: func(req usecases.RecommendationRequest) bool {
				return len(req.Filters.ExcludeTickers) == 3 &&
					req.Filters.ExcludeTickers[0] == "TSLA" &&
					req.Filters.ExcludeTickers[1] == "MSFT" &&
					req.Filters.ExcludeTickers[2] == "GOOGL"
			},
		},
		{
			name: "Parse limit parameter",
			url:  "/api/v1/recommendations?limit=25",
			setupMocks: func() {
				mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(&usecases.RecommendationResponse{}, nil)
			},
			validate: func(req usecases.RecommendationRequest) bool {
				return req.Limit == 25
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUseCase.ExpectedCalls = nil
			mockUseCase.Calls = nil

			tt.setupMocks()

			// Mock with validation
			mockUseCase.On("GetRecommendations", mock.Anything, mock.MatchedBy(tt.validate)).Return(&usecases.RecommendationResponse{}, nil)

			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			handler.GetRecommendations(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestRecommendationHandler_ErrorHandling(t *testing.T) {
	// Setup
	mockUseCase := new(MockTieredRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	tests := []struct {
		name           string
		setupMocks     func()
		expectedStatus int
		errorSubstring string
	}{
		{
			name: "UseCase error should return 500",
			setupMocks: func() {
				mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			errorSubstring: "Failed to retrieve",
		},
		{
			name: "Single ticker error should return 500",
			setupMocks: func() {
				mockUseCase.On("GetRecommendationByTicker", mock.Anything, "ERROR", enums.TIER_GUEST).Return(nil, errors.New("calculation error"))
			},
			expectedStatus: http.StatusInternalServerError,
			errorSubstring: "Failed to retrieve",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUseCase.ExpectedCalls = nil
			mockUseCase.Calls = nil

			tt.setupMocks()

			var req *http.Request
			if strings.Contains(tt.name, "Single ticker") {
				req = httptest.NewRequest("GET", "/api/v1/recommendations/ERROR", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("ticker", "ERROR")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			} else {
				req = httptest.NewRequest("GET", "/api/v1/recommendations", nil)
			}

			w := httptest.NewRecorder()

			if strings.Contains(tt.name, "Single ticker") {
				handler.GetRecommendationByTicker(w, req)
			} else {
				handler.GetRecommendations(w, req)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], tt.errorSubstring)
		})
	}
}
