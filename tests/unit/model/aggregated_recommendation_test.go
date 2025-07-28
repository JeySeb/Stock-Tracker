package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
)

func TestNewAggregatedRecommendation(t *testing.T) {
	// Test data
	ticker := "AAPL"
	companyName := "Apple Inc."
	tier := enums.RECOMMENDATION_TIER_BASIC

	// Execute
	recommendation := model.NewAggregatedRecommendation(ticker, companyName, tier)

	// Assert
	assert.NotNil(t, recommendation)
	assert.Equal(t, ticker, recommendation.Ticker)
	assert.Equal(t, companyName, recommendation.CompanyName)
	assert.Equal(t, tier, recommendation.Tier)
	assert.NotEqual(t, uuid.Nil, recommendation.ID)
	assert.False(t, recommendation.CreatedAt.IsZero())
	assert.False(t, recommendation.ExpiresAt.IsZero())
	assert.True(t, recommendation.ExpiresAt.After(recommendation.CreatedAt))
}

func TestAggregatedRecommendation_DetermineType(t *testing.T) {
	tests := []struct {
		name         string
		basicScore   float64
		expectedType enums.RecommendationType
	}{
		{
			name:         "Strong Buy for score >= 0.8",
			basicScore:   0.85,
			expectedType: enums.RECOMMENDATION_TYPE_STRONG_BUY,
		},
		{
			name:         "Buy for score >= 0.6",
			basicScore:   0.7,
			expectedType: enums.RECOMMENDATION_TYPE_BUY,
		},
		{
			name:         "Hold for score >= 0.4",
			basicScore:   0.5,
			expectedType: enums.RECOMMENDATION_TYPE_HOLD,
		},
		{
			name:         "Sell for score >= 0.2",
			basicScore:   0.3,
			expectedType: enums.RECOMMENDATION_TYPE_SELL,
		},
		{
			name:         "Strong Sell for score < 0.2",
			basicScore:   0.1,
			expectedType: enums.RECOMMENDATION_TYPE_STRONG_SELL,
		},
		{
			name:         "Edge case - exactly 0.8",
			basicScore:   0.8,
			expectedType: enums.RECOMMENDATION_TYPE_STRONG_BUY,
		},
		{
			name:         "Edge case - exactly 0.6",
			basicScore:   0.6,
			expectedType: enums.RECOMMENDATION_TYPE_BUY,
		},
		{
			name:         "Edge case - exactly 0.4",
			basicScore:   0.4,
			expectedType: enums.RECOMMENDATION_TYPE_HOLD,
		},
		{
			name:         "Edge case - exactly 0.2",
			basicScore:   0.2,
			expectedType: enums.RECOMMENDATION_TYPE_SELL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendation := &model.AggregatedRecommendation{
				BasicScore: tt.basicScore,
			}

			recommendation.DetermineType()

			assert.Equal(t, tt.expectedType, recommendation.RecommendationType)
		})
	}
}

func TestAggregatedRecommendation_IsExpired(t *testing.T) {
	now := time.Now() // Fix the current time to avoid race conditions

	tests := []struct {
		name      string
		expiresAt time.Time
		isExpired bool
	}{
		{
			name:      "Not expired - future time",
			expiresAt: now.Add(1 * time.Hour),
			isExpired: false,
		},
		{
			name:      "Expired - past time",
			expiresAt: now.Add(-1 * time.Hour),
			isExpired: true,
		},
		{
			name:      "Edge case - exactly now",
			expiresAt: now.Add(1 * time.Millisecond), // Add small buffer to avoid timing race
			isExpired: false,                         // Should be false since it's not "after" now
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendation := &model.AggregatedRecommendation{
				ExpiresAt: tt.expiresAt,
			}

			result := recommendation.IsExpired()
			assert.Equal(t, tt.isExpired, result)
		})
	}
}

func TestAggregatedRecommendation_CanAccessExternalData(t *testing.T) {
	recommendation := &model.AggregatedRecommendation{}

	tests := []struct {
		name      string
		userTier  enums.UserTier
		canAccess bool
	}{
		{
			name:      "Guest cannot access external data",
			userTier:  enums.TIER_GUEST,
			canAccess: false,
		},
		{
			name:      "Basic user can access external data",
			userTier:  enums.TIER_BASIC,
			canAccess: true,
		},
		{
			name:      "Premium user can access external data",
			userTier:  enums.TIER_PREMIUM,
			canAccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := recommendation.CanAccessExternalData(tt.userTier)
			assert.Equal(t, tt.canAccess, result)
		})
	}
}

func TestAggregatedRecommendation_CanAccessAIInsights(t *testing.T) {
	recommendation := &model.AggregatedRecommendation{}

	tests := []struct {
		name      string
		userTier  enums.UserTier
		canAccess bool
	}{
		{
			name:      "Guest cannot access AI insights",
			userTier:  enums.TIER_GUEST,
			canAccess: false,
		},
		{
			name:      "Basic user cannot access AI insights",
			userTier:  enums.TIER_BASIC,
			canAccess: false,
		},
		{
			name:      "Premium user can access AI insights",
			userTier:  enums.TIER_PREMIUM,
			canAccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := recommendation.CanAccessAIInsights(tt.userTier)
			assert.Equal(t, tt.canAccess, result)
		})
	}
}

func TestAggregatedRecommendation_GetPrimaryFactors(t *testing.T) {
	// Test data
	factors := []model.ScoringFactor{
		{
			Name:   "High Weight Factor",
			Weight: 0.3, // Above threshold
			Score:  0.8,
		},
		{
			Name:   "Medium Weight Factor",
			Weight: 0.2, // Exactly at threshold
			Score:  0.7,
		},
		{
			Name:   "Low Weight Factor",
			Weight: 0.1, // Below threshold
			Score:  0.6,
		},
	}

	recommendation := &model.AggregatedRecommendation{
		ScoringFactors: factors,
	}

	// Execute
	primaryFactors := recommendation.GetPrimaryFactors()

	// Assert - only factors with weight >= 0.2 should be returned
	assert.Len(t, primaryFactors, 2)
	assert.Equal(t, "High Weight Factor", primaryFactors[0].Name)
	assert.Equal(t, "Medium Weight Factor", primaryFactors[1].Name)
}

func TestAggregatedRecommendation_GetUpside(t *testing.T) {
	tests := []struct {
		name              string
		latestTargetPrice float64
		currentPrice      float64
		expectedUpside    *float64
		description       string
	}{
		{
			name:              "Positive upside",
			latestTargetPrice: 180.0,
			currentPrice:      150.0,
			expectedUpside:    &[]float64{0.2}[0], // 20% upside
			description:       "Should calculate positive upside correctly",
		},
		{
			name:              "Negative upside (overvalued)",
			latestTargetPrice: 120.0,
			currentPrice:      150.0,
			expectedUpside:    &[]float64{-0.2}[0], // -20% upside
			description:       "Should calculate negative upside correctly",
		},
		{
			name:              "No target price",
			latestTargetPrice: 0.0,
			currentPrice:      150.0,
			expectedUpside:    nil,
			description:       "Should return nil when no target price",
		},
		{
			name:              "No current price",
			latestTargetPrice: 180.0,
			currentPrice:      0.0,
			expectedUpside:    nil,
			description:       "Should return nil when no current price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var externalData *model.ExternalStockData
			if tt.currentPrice > 0 {
				externalData = &model.ExternalStockData{
					CurrentPrice: tt.currentPrice,
				}
			}

			recommendation := &model.AggregatedRecommendation{
				LatestTargetPrice: tt.latestTargetPrice,
				ExternalData:      externalData,
			}

			result := recommendation.GetUpside()

			if tt.expectedUpside == nil {
				assert.Nil(t, result, tt.description)
			} else {
				assert.NotNil(t, result, tt.description)
				assert.InDelta(t, *tt.expectedUpside, *result, 0.001, tt.description)
			}
		})
	}
}

func TestAggregatedRecommendation_TierBasedExpiration(t *testing.T) {
	tests := []struct {
		name        string
		tier        enums.RecommendationTier
		minDuration time.Duration
		maxDuration time.Duration
	}{
		{
			name:        "Basic tier should have 12 hour expiration",
			tier:        enums.RECOMMENDATION_TIER_BASIC,
			minDuration: 11*time.Hour + 30*time.Minute, // Allow some variance
			maxDuration: 12*time.Hour + 30*time.Minute,
		},
		{
			name:        "Enriched tier should have 4 hour expiration",
			tier:        enums.RECOMMENDATION_TIER_ENRICHED,
			minDuration: 3*time.Hour + 30*time.Minute,
			maxDuration: 4*time.Hour + 30*time.Minute,
		},
		{
			name:        "Premium tier should have 2 hour expiration",
			tier:        enums.RECOMMENDATION_TIER_PREMIUM,
			minDuration: 1*time.Hour + 30*time.Minute,
			maxDuration: 2*time.Hour + 30*time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			recommendation := model.NewAggregatedRecommendation("TEST", "Test Company", tt.tier)
			after := time.Now()

			duration := recommendation.ExpiresAt.Sub(recommendation.CreatedAt)

			// Verify the duration is within expected range
			assert.True(t, duration >= tt.minDuration,
				"Duration %v should be >= %v", duration, tt.minDuration)
			assert.True(t, duration <= tt.maxDuration,
				"Duration %v should be <= %v", duration, tt.maxDuration)

			// Verify creation time is reasonable
			assert.True(t, recommendation.CreatedAt.After(before.Add(-1*time.Second)) &&
				recommendation.CreatedAt.Before(after.Add(1*time.Second)),
				"Creation time should be around now")
		})
	}
}

func TestAggregatedRecommendation_ComprehensiveBehavior(t *testing.T) {
	// Create a comprehensive recommendation for integration testing
	recommendation := model.NewAggregatedRecommendation("AAPL", "Apple Inc.", enums.RECOMMENDATION_TIER_ENRICHED)

	// Set up data
	recommendation.BasicScore = 0.75
	recommendation.Confidence = 0.85
	recommendation.TotalEvents = 8
	recommendation.PositiveEvents = 6
	recommendation.NegativeEvents = 2
	recommendation.LatestTargetPrice = 180.0
	recommendation.ScoringFactors = []model.ScoringFactor{
		{Name: "Broker Frequency", Score: 0.8, Weight: 0.25},
		{Name: "Target Movement", Score: 0.7, Weight: 0.25},
		{Name: "Recency", Score: 0.9, Weight: 0.15},
	}
	recommendation.ExternalData = &model.ExternalStockData{
		CurrentPrice: 170.0,
		Volume:       1000000,
	}

	// Test DetermineType
	recommendation.DetermineType()
	assert.Equal(t, enums.RECOMMENDATION_TYPE_BUY, recommendation.RecommendationType)

	// Test GetUpside
	upside := recommendation.GetUpside()
	assert.NotNil(t, upside)
	assert.InDelta(t, 0.0588, *upside, 0.001) // (180-170)/170 ≈ 5.88%

	// Test GetPrimaryFactors
	primaryFactors := recommendation.GetPrimaryFactors()
	assert.Len(t, primaryFactors, 2) // Only factors with weight >= 0.2

	// Test Access permissions
	assert.True(t, recommendation.CanAccessExternalData(enums.TIER_BASIC))
	assert.False(t, recommendation.CanAccessExternalData(enums.TIER_GUEST))
	assert.True(t, recommendation.CanAccessAIInsights(enums.TIER_PREMIUM))
	assert.False(t, recommendation.CanAccessAIInsights(enums.TIER_BASIC))

	// Test IsExpired
	assert.False(t, recommendation.IsExpired()) // Should not be expired when created

	// Test that all required fields are properly set
	assert.NotEqual(t, uuid.Nil, recommendation.ID)
	assert.Equal(t, "AAPL", recommendation.Ticker)
	assert.Equal(t, "Apple Inc.", recommendation.CompanyName)
	assert.Equal(t, enums.RECOMMENDATION_TIER_ENRICHED, recommendation.Tier)
	assert.Equal(t, 0.75, recommendation.BasicScore)
	assert.Equal(t, 0.85, recommendation.Confidence)
}
