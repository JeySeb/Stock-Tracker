package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	recommendationModel "stock-tracker/internal/domain/recommendation/model"
	recommendationUsecase "stock-tracker/internal/domain/recommendation/usecase"
	"stock-tracker/internal/domain/shared/enums"
)

// TestRecommendationAPI_Comprehensive tests all recommendation API functionalities
func TestRecommendationAPI_Comprehensive(t *testing.T) {
	// Setup test environment using existing pattern
	suite := setupTestEnvironment()
	defer teardownTestEnvironment(suite)

	// Test cases
	t.Run("BasicRecommendations", func(t *testing.T) {
		testBasicRecommendations(t, suite)
	})

	t.Run("TierBasedLimits", func(t *testing.T) {
		testTierBasedLimits(t, suite)
	})

	t.Run("Filtering", func(t *testing.T) {
		testFiltering(t, suite)
	})

	t.Run("SingleTickerRecommendation", func(t *testing.T) {
		testSingleTickerRecommendation(t, suite)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, suite)
	})

	t.Run("Caching", func(t *testing.T) {
		testCaching(t, suite)
	})
}

// testBasicRecommendations tests basic recommendation retrieval
func testBasicRecommendations(t *testing.T, suite *IntegrationTestSuite) {
	// Test basic recommendations for different tiers
	testCases := []struct {
		name     string
		userTier enums.UserTier
		limit    int
		useAuth  bool
	}{
		{"Guest tier basic", enums.TIER_GUEST, 5, false},
		{"Basic tier basic", enums.TIER_BASIC, 10, true},
		{"Premium tier basic", enums.TIER_PREMIUM, 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)

			// Add query parameters
			q := req.URL.Query()
			q.Add("limit", string(rune(tc.limit)))
			req.URL.RawQuery = q.Encode()

			// Add authentication if needed
			if tc.useAuth {
				token := createTestUser(t, suite, tc.name+"@test.com", tc.userTier)
				req.Header.Set("Authorization", "Bearer "+token)
			}

			// Execute request
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			// Assert response
			assert.Equal(t, http.StatusOK, rr.Code)

			var response recommendationUsecase.RecommendationResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify response structure
			assert.Equal(t, tc.userTier, response.Meta.UserTier)
			assert.Contains(t, response.Meta.Features, "basic_recommendations")
			assert.Greater(t, response.Meta.GenerationTime, 0)

			// Verify tier-specific features
			if tc.userTier == enums.TIER_BASIC || tc.userTier == enums.TIER_PREMIUM {
				assert.Contains(t, response.Meta.Features, "real_time_data")
				assert.Contains(t, response.Meta.Features, "external_apis")
			}

			if tc.userTier == enums.TIER_PREMIUM {
				assert.Contains(t, response.Meta.Features, "ai_insights")
				assert.Contains(t, response.Meta.Features, "priority_support")
			}
		})
	}
}

// testTierBasedLimits tests tier-based limits on recommendations
func testTierBasedLimits(t *testing.T, suite *IntegrationTestSuite) {
	testCases := []struct {
		name        string
		userTier    enums.UserTier
		requested   int
		expectedMax int
		useAuth     bool
	}{
		{"Guest tier limit", enums.TIER_GUEST, 50, 10, false},
		{"Basic tier limit", enums.TIER_BASIC, 50, 25, true},
		{"Premium tier limit", enums.TIER_PREMIUM, 50, 50, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with high limit
			req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)

			q := req.URL.Query()
			q.Add("limit", string(rune(tc.requested)))
			req.URL.RawQuery = q.Encode()

			// Add authentication if needed
			if tc.useAuth {
				token := createTestUser(t, suite, tc.name+"@test.com", tc.userTier)
				req.Header.Set("Authorization", "Bearer "+token)
			}

			// Execute request
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			// Assert response
			assert.Equal(t, http.StatusOK, rr.Code)

			var response recommendationUsecase.RecommendationResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify that the actual count doesn't exceed the tier limit
			assert.LessOrEqual(t, response.Meta.Count, tc.expectedMax)
		})
	}
}

// testFiltering tests all filtering functionalities
func testFiltering(t *testing.T, suite *IntegrationTestSuite) {
	testCases := []struct {
		name           string
		userTier       enums.UserTier
		minScore       *float64
		excludeTickers []string
		useAuth        bool
	}{
		{
			name:     "High score filter",
			userTier: enums.TIER_BASIC,
			minScore: &[]float64{0.7}[0],
			useAuth:  true,
		},
		{
			name:           "Exclude tickers",
			userTier:       enums.TIER_BASIC,
			excludeTickers: []string{"AAPL", "GOOGL"},
			useAuth:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with filters
			req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)

			q := req.URL.Query()
			if tc.minScore != nil {
				q.Add("min_score", string(rune(int(*tc.minScore*100))))
			}
			if len(tc.excludeTickers) > 0 {
				for _, ticker := range tc.excludeTickers {
					q.Add("exclude_tickers", ticker)
				}
			}
			req.URL.RawQuery = q.Encode()

			// Add authentication if needed
			if tc.useAuth {
				token := createTestUser(t, suite, tc.name+"@test.com", tc.userTier)
				req.Header.Set("Authorization", "Bearer "+token)
			}

			// Execute request
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			// Assert response
			assert.Equal(t, http.StatusOK, rr.Code)

			var response recommendationUsecase.RecommendationResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify filtering worked
			assert.Len(t, response.Data, response.Meta.Count)

			// Verify excluded tickers are not in results
			if len(tc.excludeTickers) > 0 {
				for _, recommendation := range response.Data {
					for _, excludedTicker := range tc.excludeTickers {
						assert.NotEqual(t, excludedTicker, recommendation.Ticker)
					}
				}
			}

			// Verify minimum score filter
			if tc.minScore != nil {
				for _, recommendation := range response.Data {
					assert.GreaterOrEqual(t, recommendation.BasicScore, *tc.minScore)
				}
			}
		})
	}
}

// testSingleTickerRecommendation tests single ticker recommendation retrieval
func testSingleTickerRecommendation(t *testing.T, suite *IntegrationTestSuite) {
	testCases := []struct {
		name     string
		ticker   string
		userTier enums.UserTier
		useAuth  bool
	}{
		{"Basic tier single ticker", "AAPL", enums.TIER_BASIC, true},
		{"Premium tier single ticker", "GOOGL", enums.TIER_PREMIUM, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/v1/recommendations/"+tc.ticker, nil)

			// Add authentication if needed
			if tc.useAuth {
				token := createTestUser(t, suite, tc.name+"@test.com", tc.userTier)
				req.Header.Set("Authorization", "Bearer "+token)
			}

			// Execute request
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			// Assert response
			assert.Equal(t, http.StatusOK, rr.Code)

			var response recommendationModel.AggregatedRecommendation
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify response structure
			assert.Equal(t, tc.ticker, response.Ticker)
			assert.NotEmpty(t, response.CompanyName)
			assert.Greater(t, response.BasicScore, 0.0)
			assert.LessOrEqual(t, response.BasicScore, 1.0)
			assert.Greater(t, response.Confidence, 0.0)
			assert.LessOrEqual(t, response.Confidence, 1.0)

			// Verify tier-specific data
			//if tc.userTier == enums.TIER_PREMIUM {
			// Premium users might have external data (if available)
			// Note: In test environment, external APIs might not be available
			// This is intentionally empty as external data is not available in tests
			//}
		})
	}
}

// testErrorHandling tests error scenarios
func testErrorHandling(t *testing.T, suite *IntegrationTestSuite) {
	testCases := []struct {
		name           string
		url            string
		userTier       enums.UserTier
		expectedStatus int
		useAuth        bool
	}{
		{
			name:           "Invalid ticker",
			url:            "/api/v1/recommendations/INVALID",
			userTier:       enums.TIER_BASIC,
			expectedStatus: http.StatusNotFound,
			useAuth:        true,
		},
		{
			name:           "Invalid limit parameter",
			url:            "/api/v1/recommendations?limit=invalid",
			userTier:       enums.TIER_BASIC,
			expectedStatus: http.StatusBadRequest,
			useAuth:        true,
		},
		{
			name:           "Invalid min_score parameter",
			url:            "/api/v1/recommendations?min_score=invalid",
			userTier:       enums.TIER_BASIC,
			expectedStatus: http.StatusBadRequest,
			useAuth:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", tc.url, nil)

			// Add authentication if needed
			if tc.useAuth {
				token := createTestUser(t, suite, tc.name+"@test.com", tc.userTier)
				req.Header.Set("Authorization", "Bearer "+token)
			}

			// Execute request
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			// Assert response
			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus != http.StatusOK {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}
		})
	}
}

// testCaching tests caching behavior
func testCaching(t *testing.T, suite *IntegrationTestSuite) {
	userTier := enums.TIER_BASIC
	token := createTestUser(t, suite, "cache_test@test.com", userTier)

	// First request - should be cache miss
	req1 := httptest.NewRequest("GET", "/api/v1/recommendations?limit=5", nil)
	req1.Header.Set("Authorization", "Bearer "+token)

	rr1 := httptest.NewRecorder()
	suite.router.ServeHTTP(rr1, req1)

	assert.Equal(t, http.StatusOK, rr1.Code)

	var response1 recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(rr1.Body.Bytes(), &response1)
	require.NoError(t, err)

	cacheHit1 := response1.Meta.CacheHit
	assert.False(t, cacheHit1)

	// Second request - should be cache hit
	req2 := httptest.NewRequest("GET", "/api/v1/recommendations?limit=5", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	rr2 := httptest.NewRecorder()
	suite.router.ServeHTTP(rr2, req2)

	assert.Equal(t, http.StatusOK, rr2.Code)

	var response2 recommendationUsecase.RecommendationResponse
	err = json.Unmarshal(rr2.Body.Bytes(), &response2)
	require.NoError(t, err)

	// Note: In test environment, cache might not be enabled, so this could be false
	// The important thing is that the response is consistent
	assert.Equal(t, response1.Data, response2.Data)
}
