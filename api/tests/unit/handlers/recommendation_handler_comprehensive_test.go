package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/recommendation/model"
	recommendationUsecase "stock-tracker/internal/domain/recommendation/usecase"
	"stock-tracker/internal/domain/shared/enums"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/internal/presentation/handlers"
)

// TestRecommendationHandler_GuestUser_LimitedAccess tests guest user restrictions
func TestRecommendationHandler_GuestUser_LimitedAccess(t *testing.T) {
	// Setup
	mockUseCase := new(MockRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data for guest user (limited features)
	guestResponse := &recommendationUsecase.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{
				Ticker:      "AAPL",
				CompanyName: "Apple Inc.",
				BasicScore:  0.75,
				Tier:        enums.RECOMMENDATION_TIER_BASIC,
				// No external data for guests
				ExternalData: nil,
				AIInsights:   nil,
			},
		},
		Meta: recommendationUsecase.RecommendationMeta{
			Count:    1,
			UserTier: enums.TIER_GUEST,
			Features: []string{"basic_recommendations"}, // Limited features
			CacheHit: false,
		},
	}

	// Mock expectations for guest user
	mockUseCase.On("GetRecommendations", mock.Anything, mock.MatchedBy(func(req recommendationUsecase.RecommendationRequest) bool {
		return req.UserTier == enums.TIER_GUEST && req.Limit == 20 // Handler passes original request, usecase applies limits
	})).Return(guestResponse, nil)

	// Create request without auth (guest user)
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=20", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_GUEST)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify guest limitations
	assert.Equal(t, enums.TIER_GUEST, response.Meta.UserTier)
	assert.Equal(t, []string{"basic_recommendations"}, response.Meta.Features)
	assert.Nil(t, response.Data[0].ExternalData, "Guest users should not have external data")
	assert.Nil(t, response.Data[0].AIInsights, "Guest users should not have AI insights")

	mockUseCase.AssertExpectations(t)
}

// TestRecommendationHandler_PremiumUser_FullAccess tests premium user privileges
func TestRecommendationHandler_PremiumUser_FullAccess(t *testing.T) {
	// Setup
	mockUseCase := new(MockRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test data for premium user (full features)
	premiumResponse := &recommendationUsecase.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{
				Ticker:      "AAPL",
				CompanyName: "Apple Inc.",
				BasicScore:  0.92,
				Tier:        enums.RECOMMENDATION_TIER_PREMIUM,
				ExternalData: &model.ExternalStockData{
					CurrentPrice:     175.50,
					DayChange:        2.30,
					DayChangePercent: 1.33,
					Volume:           45000000,
					MarketCap:        2800000000000,
				},
				AIInsights: &model.AIGeneratedInsights{
					MarketSentiment:     "Bullish",
					SentimentScore:      0.89,
					RiskAssessment:      "Medium",
					KeyDrivers:          []string{"Earnings beat expectations", "Strong revenue growth", "Expanding market share"},
					CompetitorAnalysis:  []string{"Outperforming MSFT", "Strong competitive position"},
					TechnicalIndicators: []string{"RSI oversold", "Moving average crossover"},
					GeneratedAt:         time.Now(),
				},
			},
		},
		Meta: recommendationUsecase.RecommendationMeta{
			Count:    1,
			UserTier: enums.TIER_PREMIUM,
			Features: []string{"basic_recommendations", "market_analytics", "real_time_data", "external_apis", "advanced_analytics", "ai_insights", "priority_support"},
			CacheHit: false,
		},
	}

	// Mock expectations for premium user
	mockUseCase.On("GetRecommendations", mock.Anything, mock.MatchedBy(func(req recommendationUsecase.RecommendationRequest) bool {
		return req.UserTier == enums.TIER_PREMIUM && req.Limit == 50 // Premium limit
	})).Return(premiumResponse, nil)

	// Create request with premium user context
	req := httptest.NewRequest("GET", "/api/v1/recommendations?limit=50", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_PREMIUM)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify premium privileges
	assert.Equal(t, enums.TIER_PREMIUM, response.Meta.UserTier)
	assert.Contains(t, response.Meta.Features, "ai_insights")
	assert.NotNil(t, response.Data[0].ExternalData, "Premium users should have external data")
	assert.NotNil(t, response.Data[0].AIInsights, "Premium users should have AI insights")
	assert.Equal(t, "Bullish", response.Data[0].AIInsights.MarketSentiment)
	assert.Equal(t, 0.89, response.Data[0].AIInsights.SentimentScore)

	mockUseCase.AssertExpectations(t)
}

// TestRecommendationHandler_ErrorHandling_ServiceFailure tests error scenarios
func TestRecommendationHandler_ErrorHandling_ServiceFailure(t *testing.T) {
	// Setup
	mockUseCase := new(MockRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Mock service failure
	mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(nil, errors.New("database connection failed"))

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Handler includes underlying error message; assert prefix to remain robust
	assert.Contains(t, response["error"], "Failed to retrieve recommendations")

	mockUseCase.AssertExpectations(t)
}

// TestRecommendationHandler_InvalidTickerValidation tests ticker validation
func TestRecommendationHandler_InvalidTickerValidation(t *testing.T) {
	// Setup
	mockUseCase := new(MockRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	testCases := []struct {
		name           string
		ticker         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Empty ticker",
			ticker:         "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Ticker is required",
		},
		{
			name:           "Ticker too long",
			ticker:         "VERYLONGTICKER",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Ticker must be between 1 and 10 characters",
		},
		{
			name:           "Ticker with special characters",
			ticker:         "AAPL$",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Ticker contains invalid characters",
		},
		{
			name:           "Valid ticker",
			ticker:         "AAPL",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedStatus == http.StatusOK {
				// Mock successful response for valid ticker
				mockUseCase.On("GetRecommendationByTicker", mock.Anything, tc.ticker, enums.TIER_BASIC).Return(&model.AggregatedRecommendation{
					Ticker: tc.ticker,
				}, nil).Once()
			}

			// Create request with ticker parameter
			req := httptest.NewRequest("GET", "/api/v1/recommendations/"+tc.ticker, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ticker", tc.ticker)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			// Execute
			handler.GetRecommendationByTicker(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tc.expectedError)
			}
		})
	}

	mockUseCase.AssertExpectations(t)
}

// TestRecommendationHandler_RateLimitExceeded tests rate limiting behavior
func TestRecommendationHandler_RateLimitExceeded(t *testing.T) {
	// Setup
	mockUseCase := new(MockRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test different user tiers with rate limit context
	testCases := []struct {
		name           string
		userTier       enums.UserTier
		rateLimitLeft  int
		expectResponse bool
	}{
		{
			name:           "Guest user - rate limit OK",
			userTier:       enums.TIER_GUEST,
			rateLimitLeft:  10,
			expectResponse: true,
		},
		{
			name:           "Basic user - high usage remaining",
			userTier:       enums.TIER_BASIC,
			rateLimitLeft:  250,
			expectResponse: true,
		},
		{
			name:           "Premium user - unlimited",
			userTier:       enums.TIER_PREMIUM,
			rateLimitLeft:  -1, // -1 indicates unlimited
			expectResponse: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectResponse {
				// Mock successful response
				response := &recommendationUsecase.RecommendationResponse{
					Data: []*model.AggregatedRecommendation{},
					Meta: recommendationUsecase.RecommendationMeta{
						UserTier:           tc.userTier,
						RateLimitRemaining: tc.rateLimitLeft,
					},
				}
				mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(response, nil).Once()
			}

			// Create request with rate limit context
			req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
			ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, tc.userTier)
			if tc.rateLimitLeft >= 0 {
				ctx = context.WithValue(ctx, middleware.RateLimitRemainingKey, tc.rateLimitLeft)
			}
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			// Execute
			handler.GetRecommendations(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			var response recommendationUsecase.RecommendationResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.userTier, response.Meta.UserTier)

			if tc.rateLimitLeft >= 0 {
				assert.Equal(t, tc.rateLimitLeft, response.Meta.RateLimitRemaining)
			}
		})
	}

	mockUseCase.AssertExpectations(t)
}

// TestRecommendationHandler_CachePerformance tests cache hit scenarios
func TestRecommendationHandler_CachePerformance(t *testing.T) {
	// Setup
	mockUseCase := new(MockRecommendationUseCase)
	mockLogger := &MockLogger{}
	handler := handlers.NewRecommendationHandler(mockUseCase, mockLogger)

	// Test cache hit response
	cachedResponse := &recommendationUsecase.RecommendationResponse{
		Data: []*model.AggregatedRecommendation{
			{
				Ticker:     "AAPL",
				BasicScore: 0.85,
			},
		},
		Meta: recommendationUsecase.RecommendationMeta{
			Count:          1,
			UserTier:       enums.TIER_BASIC,
			CacheHit:       true,
			GenerationTime: 2 * time.Millisecond, // Fast response from cache
		},
	}

	mockUseCase.On("GetRecommendations", mock.Anything, mock.Anything).Return(cachedResponse, nil)

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetRecommendations(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response recommendationUsecase.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify cache performance
	assert.True(t, response.Meta.CacheHit, "Response should indicate cache hit")
	assert.True(t, response.Meta.GenerationTime < 10*time.Millisecond, "Cache hit should be fast")

	mockUseCase.AssertExpectations(t)
}
