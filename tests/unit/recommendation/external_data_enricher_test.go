package recommendation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/application/recommendation"
	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/internal/infrastructure/external"
)

// MockYahooFinanceClient mock implementation
type MockYahooFinanceClient struct {
	mock.Mock
}

func (m *MockYahooFinanceClient) GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExternalStockData), args.Error(1)
}

func (m *MockYahooFinanceClient) GetHistoricalData(ctx context.Context, symbol string, period string) ([]external.HistoricalDataPoint, error) {
	args := m.Called(ctx, symbol, period)
	return args.Get(0).([]external.HistoricalDataPoint), args.Error(1)
}

// MockAlphaVantageClient mock implementation
type MockAlphaVantageClient struct {
	mock.Mock
}

func (m *MockAlphaVantageClient) GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExternalStockData), args.Error(1)
}

func (m *MockAlphaVantageClient) GetCompanyOverview(ctx context.Context, symbol string) (*external.CompanyOverview, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.CompanyOverview), args.Error(1)
}

// MockCache mock implementation
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestExternalDataEnricher_EnrichRecommendation_Success(t *testing.T) {
	// Setup mocks
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	// Test data
	baseRecommendation := &model.AggregatedRecommendation{
		Ticker:            "AAPL",
		CompanyName:       "Apple Inc.",
		LatestTargetPrice: 180.0,
		BasicScore:        0.7,
		Tier:              enums.RECOMMENDATION_TIER_BASIC,
	}

	externalData := &model.ExternalStockData{
		CurrentPrice:     170.0,
		DayChange:        2.5,
		DayChangePercent: 1.5,
		Volume:           1000000,
		MarketCap:        2800000000000,
		PERatio:          &[]float64{28.5}[0],
		Week52High:       &[]float64{190.0}[0],
		Week52Low:        &[]float64{140.0}[0],
		AvgVolume:        &[]int64{800000}[0],
		LastUpdated:      time.Now(),
	}

	// Mock expectations
	mockCache.On("Get", mock.Anything, "external_data:AAPL", mock.Anything).Return(cache.ErrCacheMiss)
	mockYahoo.On("GetQuote", mock.Anything, "AAPL").Return(externalData, nil)
	mockCache.On("Set", mock.Anything, "external_data:AAPL", externalData, 5*time.Minute).Return(nil)

	// Execute
	result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, enums.RECOMMENDATION_TIER_ENRICHED, result.Tier)
	assert.NotNil(t, result.ExternalData)
	assert.Equal(t, 170.0, result.ExternalData.CurrentPrice)
	assert.NotEqual(t, baseRecommendation.BasicScore, result.BasicScore)                  // Score should be enriched
	assert.Greater(t, len(result.ScoringFactors), len(baseRecommendation.ScoringFactors)) // More factors

	// Verify mocks
	mockYahoo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestExternalDataEnricher_EnrichRecommendation_CacheHit(t *testing.T) {
	// Setup mocks
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	baseRecommendation := &model.AggregatedRecommendation{
		Ticker:            "AAPL",
		CompanyName:       "Apple Inc.",
		LatestTargetPrice: 180.0,
		BasicScore:        0.7,
		Tier:              enums.RECOMMENDATION_TIER_BASIC,
	}

	// Mock cache hit - should not call external APIs
	mockCache.On("Get", mock.Anything, "external_data:AAPL", mock.Anything).Return(nil)

	// Execute
	result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify Yahoo was NOT called (cache hit)
	mockYahoo.AssertNotCalled(t, "GetQuote")
	mockCache.AssertExpectations(t)
}

func TestExternalDataEnricher_EnrichRecommendation_YahooFailsAlphaVantageWorks(t *testing.T) {
	// Setup mocks
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	baseRecommendation := &model.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.7,
		Tier:        enums.RECOMMENDATION_TIER_BASIC,
	}

	externalData := &model.ExternalStockData{
		CurrentPrice: 170.0,
		LastUpdated:  time.Now(),
	}

	// Mock expectations
	mockCache.On("Get", mock.Anything, "external_data:AAPL", mock.Anything).Return(cache.ErrCacheMiss)
	mockYahoo.On("GetQuote", mock.Anything, "AAPL").Return(nil, errors.New("Yahoo API failed"))
	mockAlpha.On("GetQuote", mock.Anything, "AAPL").Return(externalData, nil)
	mockCache.On("Set", mock.Anything, "external_data:AAPL", externalData, 5*time.Minute).Return(nil)

	// Execute
	result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, enums.RECOMMENDATION_TIER_ENRICHED, result.Tier)
	assert.NotNil(t, result.ExternalData)

	// Verify fallback was used
	mockYahoo.AssertExpectations(t)
	mockAlpha.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestExternalDataEnricher_EnrichRecommendation_GracefulDegradation(t *testing.T) {
	// Setup mocks
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	baseRecommendation := &model.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.7,
		Tier:        enums.RECOMMENDATION_TIER_BASIC,
	}

	// Mock expectations - both APIs fail
	mockCache.On("Get", mock.Anything, "external_data:AAPL", mock.Anything).Return(cache.ErrCacheMiss)
	mockYahoo.On("GetQuote", mock.Anything, "AAPL").Return(nil, errors.New("Yahoo API failed"))
	mockAlpha.On("GetQuote", mock.Anything, "AAPL").Return(nil, errors.New("Alpha Vantage API failed"))

	// Execute
	result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

	// Assert - should return original recommendation (graceful degradation)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, baseRecommendation, result) // Should be unchanged

	// Verify both APIs were tried
	mockYahoo.AssertExpectations(t)
	mockAlpha.AssertExpectations(t)
}

func TestExternalDataEnricher_UpsidePotentialCalculation(t *testing.T) {
	// Setup
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	tests := []struct {
		name               string
		currentPrice       float64
		targetPrice        float64
		expectedUpsideSign string // "positive", "negative", "neutral"
	}{
		{
			name:               "Strong upside potential",
			currentPrice:       100.0,
			targetPrice:        130.0, // 30% upside
			expectedUpsideSign: "positive",
		},
		{
			name:               "Negative upside (overvalued)",
			currentPrice:       150.0,
			targetPrice:        120.0, // -20% upside
			expectedUpsideSign: "negative",
		},
		{
			name:               "Neutral position",
			currentPrice:       100.0,
			targetPrice:        100.0, // 0% upside
			expectedUpsideSign: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseRecommendation := &model.AggregatedRecommendation{
				Ticker:            "TEST",
				LatestTargetPrice: tt.targetPrice,
				BasicScore:        0.5,
			}

			externalData := &model.ExternalStockData{
				CurrentPrice: tt.currentPrice,
			}

			// Mock cache miss and successful API call
			mockCache.On("Get", mock.Anything, "external_data:TEST", mock.Anything).Return(cache.ErrCacheMiss)
			mockYahoo.On("GetQuote", mock.Anything, "TEST").Return(externalData, nil)
			mockCache.On("Set", mock.Anything, "external_data:TEST", externalData, 5*time.Minute).Return(nil)

			result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Find upside potential factor
			var upsideFactor *model.ScoringFactor
			for _, factor := range result.ScoringFactors {
				if factor.Name == "Real-time Upside Potential" {
					upsideFactor = &factor
					break
				}
			}

			assert.NotNil(t, upsideFactor, "Upside potential factor should be present")

			switch tt.expectedUpsideSign {
			case "positive":
				assert.Greater(t, upsideFactor.Score, 0.6, "Strong upside should have high score")
			case "negative":
				assert.Less(t, upsideFactor.Score, 0.4, "Negative upside should have low score")
			case "neutral":
				assert.InDelta(t, 0.5, upsideFactor.Score, 0.1, "Neutral upside should have neutral score")
			}

			// Clear mock calls for next iteration
			mockCache.ExpectedCalls = nil
			mockCache.Calls = nil
			mockYahoo.ExpectedCalls = nil
			mockYahoo.Calls = nil
		})
	}
}

func TestExternalDataEnricher_VolumeActivityScoring(t *testing.T) {
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	tests := []struct {
		name          string
		currentVolume int64
		avgVolume     int64
		expectedScore string
	}{
		{
			name:          "Very high volume activity",
			currentVolume: 3000000,
			avgVolume:     1000000, // 3x average
			expectedScore: "high",
		},
		{
			name:          "Normal volume activity",
			currentVolume: 1000000,
			avgVolume:     1000000, // 1x average
			expectedScore: "medium",
		},
		{
			name:          "Low volume activity",
			currentVolume: 250000,
			avgVolume:     1000000, // 0.25x average
			expectedScore: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseRecommendation := &model.AggregatedRecommendation{
				Ticker:     "TEST",
				BasicScore: 0.5,
			}

			externalData := &model.ExternalStockData{
				CurrentPrice: 100.0,
				Volume:       tt.currentVolume,
				AvgVolume:    &tt.avgVolume,
			}

			mockCache.On("Get", mock.Anything, "external_data:TEST", mock.Anything).Return(cache.ErrCacheMiss)
			mockYahoo.On("GetQuote", mock.Anything, "TEST").Return(externalData, nil)
			mockCache.On("Set", mock.Anything, "external_data:TEST", externalData, 5*time.Minute).Return(nil)

			result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

			assert.NoError(t, err)

			// Find volume activity factor
			var volumeFactor *model.ScoringFactor
			for _, factor := range result.ScoringFactors {
				if factor.Name == "Volume Activity" {
					volumeFactor = &factor
					break
				}
			}

			assert.NotNil(t, volumeFactor)

			switch tt.expectedScore {
			case "high":
				assert.Greater(t, volumeFactor.Score, 0.7, "High volume should have high score")
			case "medium":
				assert.InDelta(t, 0.5, volumeFactor.Score, 0.2, "Normal volume should have medium score")
			case "low":
				assert.Less(t, volumeFactor.Score, 0.4, "Low volume should have low score")
			}

			// Clear mock calls
			mockCache.ExpectedCalls = nil
			mockCache.Calls = nil
			mockYahoo.ExpectedCalls = nil
			mockYahoo.Calls = nil
		})
	}
}

func TestExternalDataEnricher_GetEnrichmentPreview(t *testing.T) {
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	externalData := &model.ExternalStockData{
		CurrentPrice: 170.0,
		Volume:       1000000,
		LastUpdated:  time.Now(),
	}

	// Mock expectations
	mockCache.On("Get", mock.Anything, "external_data:AAPL", mock.Anything).Return(cache.ErrCacheMiss)
	mockYahoo.On("GetQuote", mock.Anything, "AAPL").Return(externalData, nil)
	mockCache.On("Set", mock.Anything, "external_data:AAPL", externalData, 5*time.Minute).Return(nil)

	// Execute
	result, err := enricher.GetEnrichmentPreview(context.Background(), "AAPL")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 170.0, result.CurrentPrice)
	assert.Equal(t, int64(1000000), result.Volume)

	mockYahoo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestExternalDataEnricher_ErrorHandling(t *testing.T) {
	mockYahoo := new(MockYahooFinanceClient)
	mockAlpha := new(MockAlphaVantageClient)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	enricher := recommendation.NewExternalDataEnricher(mockYahoo, mockAlpha, mockCache, mockLogger)

	tests := []struct {
		name            string
		setupMocks      func()
		expectedFailure bool
	}{
		{
			name: "Cache error should not affect enrichment",
			setupMocks: func() {
				externalData := &model.ExternalStockData{CurrentPrice: 100.0}
				mockCache.On("Get", mock.Anything, "external_data:TEST", mock.Anything).Return(cache.ErrCacheMiss)
				mockYahoo.On("GetQuote", mock.Anything, "TEST").Return(externalData, nil)
				mockCache.On("Set", mock.Anything, "external_data:TEST", externalData, 5*time.Minute).Return(errors.New("cache error"))
			},
			expectedFailure: false, // Should still work
		},
		{
			name: "Both APIs fail - graceful degradation",
			setupMocks: func() {
				mockCache.On("Get", mock.Anything, "external_data:TEST", mock.Anything).Return(cache.ErrCacheMiss)
				mockYahoo.On("GetQuote", mock.Anything, "TEST").Return(nil, errors.New("yahoo failed"))
				mockAlpha.On("GetQuote", mock.Anything, "TEST").Return(nil, errors.New("alpha failed"))
			},
			expectedFailure: false, // Should return original recommendation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCache.ExpectedCalls = nil
			mockCache.Calls = nil
			mockYahoo.ExpectedCalls = nil
			mockYahoo.Calls = nil
			mockAlpha.ExpectedCalls = nil
			mockAlpha.Calls = nil

			tt.setupMocks()

			baseRecommendation := &model.AggregatedRecommendation{
				Ticker:     "TEST",
				BasicScore: 0.5,
			}

			result, err := enricher.EnrichRecommendation(context.Background(), baseRecommendation)

			if tt.expectedFailure {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
