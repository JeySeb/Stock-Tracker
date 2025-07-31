package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	recommendationModel "stock-tracker/internal/domain/recommendation/model"
	recommendationUsecase "stock-tracker/internal/domain/recommendation/usecase"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/tests/mocks"
)

func TestTieredRecommendationUseCase_GuestUserLimitations(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(mockRepo, mockBasicCalculator, mockExternalEnricher, mockCache, mockLogger)

	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_GUEST,
		Limit:    20, // Guest requests 20 but should be limited to 10
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss
	mockCache.On("Get", mock.Anything, "recommendations:guest:10", mock.Anything).Return(cache.ErrCacheMiss)

	// Create sample ticker events
	tickerEvents := map[string][]*stockModel.Stock{
		"AAPL": {
			&stockModel.Stock{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				EventTime:  time.Now().AddDate(0, 0, -1),
			},
		},
	}

	mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(tickerEvents, nil)

	// Mock the basic calculator to return a sample recommendation
	sampleRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.75,
		Tier:        enums.RECOMMENDATION_TIER_BASIC,
	}
	mockBasicCalculator.On("CalculateAggregatedRecommendationFromEvents", mock.Anything, "AAPL", mock.Anything).Return(sampleRecommendation, nil)

	mockCache.On("Set", mock.Anything, "recommendations:guest:10", mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, enums.TIER_GUEST, result.Meta.UserTier)
	assert.Contains(t, result.Meta.Features, "basic_recommendations")
	assert.NotContains(t, result.Meta.Features, "external_apis")
	assert.NotContains(t, result.Meta.Features, "ai_insights")
	assert.False(t, result.Meta.CacheHit)

	// Verify guest user limitations are enforced
	assert.LessOrEqual(t, len(result.Data), 5, "Guest users should be limited to 5 recommendations")

	// Verify no enriched data for guest users
	for _, rec := range result.Data {
		assert.Equal(t, enums.RECOMMENDATION_TIER_BASIC, rec.Tier, "Guest users should only get basic tier recommendations")
		assert.Nil(t, rec.ExternalData, "Guest users should not have external data")
		assert.Nil(t, rec.AIInsights, "Guest users should not have AI insights")
	}

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestTieredRecommendationUseCase_BasicUserEnrichment(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo,
		mockBasicCalculator,
		mockExternalEnricher,
		mockCache,
		mockLogger,
	)

	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_BASIC,
		Limit:    20,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss
	mockCache.On("Get", mock.Anything, "recommendations:basic:20", mock.Anything).Return(cache.ErrCacheMiss)

	// Create sample ticker events
	tickerEvents := map[string][]*stockModel.Stock{
		"AAPL": {
			&stockModel.Stock{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				EventTime:  time.Now().AddDate(0, 0, -1),
			},
		},
	}

	basicRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.75,
		Tier:        enums.RECOMMENDATION_TIER_BASIC,
	}

	enrichedRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.82,
		Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
		ExternalData: &recommendationModel.ExternalStockData{
			CurrentPrice: 175.50,
			Volume:       45000000,
		},
	}

	mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(tickerEvents, nil)
	mockBasicCalculator.On("CalculateAggregatedRecommendationFromEvents", mock.Anything, "AAPL", mock.Anything).Return(basicRecommendation, nil)
	mockExternalEnricher.On("EnrichRecommendation", mock.Anything, basicRecommendation).Return(enrichedRecommendation, nil)
	mockCache.On("Set", mock.Anything, "recommendations:basic:20", mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, enums.TIER_BASIC, result.Meta.UserTier)
	assert.Contains(t, result.Meta.Features, "external_apis")
	assert.NotContains(t, result.Meta.Features, "ai_insights")
	assert.LessOrEqual(t, len(result.Data), 20)

	// Verify enrichment for basic users
	if len(result.Data) > 0 {
		assert.Equal(t, enums.RECOMMENDATION_TIER_ENRICHED, result.Data[0].Tier)
		assert.NotNil(t, result.Data[0].ExternalData, "Basic users should have external data")
		assert.Nil(t, result.Data[0].AIInsights, "Basic users should not have AI insights")
	}

	mockRepo.AssertExpectations(t)
	mockBasicCalculator.AssertExpectations(t)
	mockExternalEnricher.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestTieredRecommendationUseCase_PremiumUserFullFeatures(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo,
		mockBasicCalculator,
		mockExternalEnricher,
		mockCache,
		mockLogger,
	)

	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_PREMIUM,
		Limit:    50,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss
	mockCache.On("Get", mock.Anything, "recommendations:premium:50", mock.Anything).Return(cache.ErrCacheMiss)

	// Create sample ticker events
	tickerEvents := map[string][]*stockModel.Stock{
		"AAPL": {
			&stockModel.Stock{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				EventTime:  time.Now().AddDate(0, 0, -1),
			},
		},
	}

	basicRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.75,
		Tier:        enums.RECOMMENDATION_TIER_BASIC,
	}

	premiumRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		BasicScore:  0.92,
		Tier:        enums.RECOMMENDATION_TIER_PREMIUM,
		ExternalData: &recommendationModel.ExternalStockData{
			CurrentPrice: 175.50,
			Volume:       45000000,
		},
		AIInsights: &recommendationModel.AIGeneratedInsights{
			MarketSentiment: "Bullish",
			SentimentScore:  0.89,
			RiskAssessment:  "Medium",
			KeyDrivers:      []string{"Earnings growth", "Market expansion"},
			GeneratedAt:     time.Now(),
		},
	}

	mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(tickerEvents, nil)
	mockBasicCalculator.On("CalculateAggregatedRecommendationFromEvents", mock.Anything, "AAPL", mock.Anything).Return(basicRecommendation, nil)
	mockExternalEnricher.On("EnrichRecommendation", mock.Anything, basicRecommendation).Return(premiumRecommendation, nil)
	mockCache.On("Set", mock.Anything, "recommendations:premium:50", mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, enums.TIER_PREMIUM, result.Meta.UserTier)
	assert.Contains(t, result.Meta.Features, "ai_insights")
	assert.Contains(t, result.Meta.Features, "priority_support")
	assert.LessOrEqual(t, len(result.Data), 50)

	// Verify premium features
	if len(result.Data) > 0 {
		assert.Equal(t, enums.RECOMMENDATION_TIER_PREMIUM, result.Data[0].Tier)
		assert.NotNil(t, result.Data[0].ExternalData, "Premium users should have external data")
		assert.NotNil(t, result.Data[0].AIInsights, "Premium users should have AI insights")
		assert.Equal(t, "Bullish", result.Data[0].AIInsights.MarketSentiment)
	}

	mockRepo.AssertExpectations(t)
	mockBasicCalculator.AssertExpectations(t)
	mockExternalEnricher.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestTieredRecommendationUseCase_CacheStrategyPerTier(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(mockRepo, mockBasicCalculator, mockExternalEnricher, mockCache, mockLogger)

	testCases := []struct {
		name        string
		userTier    enums.UserTier
		expectedTTL time.Duration
		cacheKey    string
	}{
		{
			name:        "Guest user - long cache",
			userTier:    enums.TIER_GUEST,
			expectedTTL: 12 * time.Hour,
			cacheKey:    "recommendations:guest:10",
		},
		{
			name:        "Basic user - medium cache",
			userTier:    enums.TIER_BASIC,
			expectedTTL: 4 * time.Hour,
			cacheKey:    "recommendations:basic:10",
		},
		{
			name:        "Premium user - short cache for real-time data",
			userTier:    enums.TIER_PREMIUM,
			expectedTTL: 2 * time.Hour,
			cacheKey:    "recommendations:premium:10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := recommendationUsecase.RecommendationRequest{
				UserTier: tc.userTier,
				Limit:    10,
				Filters:  recommendationUsecase.RecommendationFilters{},
			}

			// Mock cache miss, then verify Set is called with correct TTL
			mockCache.On("Get", mock.Anything, tc.cacheKey, mock.Anything).Return(cache.ErrCacheMiss).Once()
			mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(map[string][]*stockModel.Stock{}, nil).Once()
			mockCache.On("Set", mock.Anything, tc.cacheKey, mock.Anything, tc.expectedTTL).Return(nil).Once()

			// Execute
			result, err := useCase.GetRecommendations(context.Background(), request)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.userTier, result.Meta.UserTier)
			assert.False(t, result.Meta.CacheHit)
		})
	}

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestTieredRecommendationUseCase_FilterApplication(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}
	mockExternalEnricher := new(MockExternalDataEnricher)

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(mockRepo, mockBasicCalculator, mockExternalEnricher, mockCache, mockLogger)

	minScore := 0.8
	recType := enums.RECOMMENDATION_TYPE_BUY
	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_BASIC,
		Limit:    10,
		Filters: recommendationUsecase.RecommendationFilters{
			MinScore:           &minScore,
			RecommendationType: &recType,
			ExcludeTickers:     []string{"TSLA", "MSFT"},
		},
	}

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cache.ErrCacheMiss)

	// Create sample ticker events including excluded tickers
	tickerEvents := map[string][]*stockModel.Stock{
		"AAPL": {
			&stockModel.Stock{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				EventTime:  time.Now().AddDate(0, 0, -1),
			},
		},
		"TSLA": { // This should be excluded
			&stockModel.Stock{
				Ticker:     "TSLA",
				Company:    "Tesla Inc.",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				EventTime:  time.Now().AddDate(0, 0, -1),
			},
		},
	}

	aaplRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:             "AAPL",
		CompanyName:        "Apple Inc.",
		BasicScore:         0.85, // Above min score
		RecommendationType: enums.RECOMMENDATION_TYPE_BUY,
		Tier:               enums.RECOMMENDATION_TIER_BASIC,
	}

	mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(tickerEvents, nil)
	mockBasicCalculator.On("CalculateAggregatedRecommendationFromEvents", mock.Anything, "AAPL", mock.Anything).Return(aaplRecommendation, nil)
	mockExternalEnricher.On("EnrichRecommendation", mock.Anything, aaplRecommendation).Return(aaplRecommendation, nil)
	// TSLA should not be calculated due to exclusion filter
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify filtering worked
	for _, rec := range result.Data {
		assert.NotContains(t, []string{"TSLA", "MSFT"}, rec.Ticker, "Excluded tickers should not be in results")
		assert.GreaterOrEqual(t, rec.BasicScore, 0.8, "All recommendations should meet minimum score")
		assert.Equal(t, enums.RECOMMENDATION_TYPE_BUY, rec.RecommendationType, "All recommendations should match type filter")
	}

	mockRepo.AssertExpectations(t)
	mockBasicCalculator.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestTieredRecommendationUseCase_ErrorHandling(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(mockRepo, mockBasicCalculator, mockExternalEnricher, mockCache, mockLogger)

	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_BASIC,
		Limit:    10,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cache.ErrCacheMiss)

	// Mock repository error
	mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(map[string][]*stockModel.Stock{}, assert.AnError)

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate recommendations")

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// Mock types for comprehensive testing
type MockBasicScoringCalculator struct {
	mock.Mock
}

func (m *MockBasicScoringCalculator) CalculateAggregatedRecommendation(ctx context.Context, ticker string) (*recommendationModel.AggregatedRecommendation, error) {
	args := m.Called(ctx, ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*recommendationModel.AggregatedRecommendation), args.Error(1)
}

func (m *MockBasicScoringCalculator) CalculateAggregatedRecommendationFromEvents(ctx context.Context, ticker string, events []*stockModel.Stock) (*recommendationModel.AggregatedRecommendation, error) {
	args := m.Called(ctx, ticker, events)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*recommendationModel.AggregatedRecommendation), args.Error(1)
}

type MockExternalDataEnricher struct {
	mock.Mock
}

func (m *MockExternalDataEnricher) EnrichRecommendation(ctx context.Context, baseRecommendation *recommendationModel.AggregatedRecommendation) (*recommendationModel.AggregatedRecommendation, error) {
	args := m.Called(ctx, baseRecommendation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*recommendationModel.AggregatedRecommendation), args.Error(1)
}
