package recommendation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	recommendationModel "stock-tracker/internal/domain/recommendation/model"
	recommendationUsecase "stock-tracker/internal/domain/recommendation/usecase"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/tests/mocks"
)

// MockBasicScoringCalculator mock implementation
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

// MockExternalDataEnricher mock implementation
type MockExternalDataEnricher struct {
	mock.Mock
}

func (m *MockExternalDataEnricher) EnrichRecommendation(ctx context.Context, recommendation *recommendationModel.AggregatedRecommendation) (*recommendationModel.AggregatedRecommendation, error) {
	args := m.Called(ctx, recommendation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*recommendationModel.AggregatedRecommendation), args.Error(1)
}

// MockCache mock implementation for recommendation usecase
type MockRecommendationCache struct {
	mock.Mock
}

func (m *MockRecommendationCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockRecommendationCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockRecommendationCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRecommendationCache) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestTieredRecommendationUseCase_GetRecommendations_BasicFunctionality tests basic recommendation retrieval
func TestTieredRecommendationUseCase_GetRecommendations_BasicFunctionality(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockStockRepository)
	mockCalculator := new(MockBasicScoringCalculator)
	mockEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockRecommendationCache)
	mockLogger := &MockLogger{}

	usecase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo, mockCalculator, mockEnricher, mockCache, mockLogger,
	)

	ctx := context.Background()
	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_BASIC,
		Limit:    10,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss
	mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*[]*model.AggregatedRecommendation")).Return(errors.New("cache miss"))

	// Mock stock repository with a single ticker to simplify
	testTickerEvents := map[string][]*stockModel.Stock{
		"AAPL": {
			{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				Brokerage:  "Goldman Sachs",
				Action:     "upgraded by",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				TargetFrom: 150.0,
				TargetTo:   160.0,
				EventTime:  time.Now().AddDate(0, 0, -1),
			},
		},
	}

	mockRepo.On("GetRecentByTickers", ctx, mock.AnythingOfType("time.Time")).Return(testTickerEvents, nil)

	// Mock calculator for the ticker
	basicRecommendation := &recommendationModel.AggregatedRecommendation{
		Ticker:            "AAPL",
		CompanyName:       "Apple Inc.",
		BasicScore:        0.75,
		Confidence:        0.8,
		Tier:              enums.RECOMMENDATION_TIER_BASIC,
		TotalEvents:       1,
		LatestTargetPrice: 160.0,
		ScoringFactors:    []recommendationModel.ScoringFactor{
			{Name: "broker_credibility", Score: 0.8, Weight: 0.3, Explanation: "High broker credibility"},
			{Name: "target_movement", Score: 0.7, Weight: 0.4, Explanation: "Positive target movement"},
		},
	}

	mockCalculator.On("CalculateAggregatedRecommendationFromEvents", ctx, "AAPL", testTickerEvents["AAPL"]).Return(basicRecommendation, nil)

	// Mock cache set
	mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]*model.AggregatedRecommendation"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := usecase.GetRecommendations(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, enums.TIER_BASIC, result.Meta.UserTier)
	assert.False(t, result.Meta.CacheHit)
	assert.Greater(t, result.Meta.GenerationTime, time.Duration(0))
	assert.Contains(t, result.Meta.Features, "basic_recommendations")
	assert.Contains(t, result.Meta.Features, "real_time_data")

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockCalculator.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestTieredRecommendationUseCase_GetRecommendations_CacheHit tests cache hit scenario
func TestTieredRecommendationUseCase_GetRecommendations_CacheHit(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockStockRepository)
	mockCalculator := new(MockBasicScoringCalculator)
	mockEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockRecommendationCache)
	mockLogger := &MockLogger{}

	usecase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo, mockCalculator, mockEnricher, mockCache, mockLogger,
	)

	ctx := context.Background()
	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_PREMIUM,
		Limit:    20,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache hit
	cachedRecommendations := []*recommendationModel.AggregatedRecommendation{
		{
			Ticker:            "AAPL",
			CompanyName:       "Apple Inc.",
			BasicScore:        0.85,
			Confidence:        0.9,
			Tier:              enums.RECOMMENDATION_TIER_PREMIUM,
			TotalEvents:       3,
			LatestTargetPrice: 180.0,
		},
	}

	mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*[]*model.AggregatedRecommendation")).
		Run(func(args mock.Arguments) {
			dest := args.Get(1).(*[]*recommendationModel.AggregatedRecommendation)
			*dest = cachedRecommendations
		}).Return(nil)

	// Execute
	result, err := usecase.GetRecommendations(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, enums.TIER_PREMIUM, result.Meta.UserTier)
	assert.True(t, result.Meta.CacheHit)
	assert.Equal(t, "AAPL", result.Data[0].Ticker)
	assert.Contains(t, result.Meta.Features, "ai_insights")
	assert.Contains(t, result.Meta.Features, "priority_support")

	// Verify that repository and calculator were not called due to cache hit
	mockRepo.AssertNotCalled(t, "GetRecentByTickers")
	mockCalculator.AssertNotCalled(t, "CalculateAggregatedRecommendationFromEvents")
}

// TestTieredRecommendationUseCase_GetRecommendations_EmptyResults tests handling of empty results
func TestTieredRecommendationUseCase_GetRecommendations_EmptyResults(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockStockRepository)
	mockCalculator := new(MockBasicScoringCalculator)
	mockEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockRecommendationCache)
	mockLogger := &MockLogger{}

	usecase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo, mockCalculator, mockEnricher, mockCache, mockLogger,
	)

	ctx := context.Background()
	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_BASIC,
		Limit:    10,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss
	mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*[]*model.AggregatedRecommendation")).Return(errors.New("cache miss"))

	// Mock empty repository
	mockRepo.On("GetRecentByTickers", ctx, mock.AnythingOfType("time.Time")).Return(map[string][]*stockModel.Stock{}, nil)

	// Mock cache set
	mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]*model.AggregatedRecommendation"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := usecase.GetRecommendations(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 0)
	assert.Equal(t, 0, result.Meta.Count)
	assert.False(t, result.Meta.CacheHit)
	assert.Greater(t, result.Meta.GenerationTime, time.Duration(0))
}

// TestTieredRecommendationUseCase_GetRecommendations_ErrorHandling tests error scenarios
func TestTieredRecommendationUseCase_GetRecommendations_ErrorHandling(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockStockRepository)
	mockCalculator := new(MockBasicScoringCalculator)
	mockEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockRecommendationCache)
	mockLogger := &MockLogger{}

	usecase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo, mockCalculator, mockEnricher, mockCache, mockLogger,
	)

	ctx := context.Background()
	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_BASIC,
		Limit:    10,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	t.Run("Repository error", func(t *testing.T) {
		// Mock cache miss
		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*[]*model.AggregatedRecommendation")).Return(errors.New("cache miss"))

		// Mock repository error
		mockRepo.On("GetRecentByTickers", ctx, mock.AnythingOfType("time.Time")).Return(nil, errors.New("database error"))

		// Execute
		result, err := usecase.GetRecommendations(ctx, request)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to generate recommendations")
	})

	t.Run("Cache set error", func(t *testing.T) {
		// Mock cache miss
		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*[]*model.AggregatedRecommendation")).Return(errors.New("cache miss"))

		// Mock successful repository call
		testTickerEvents := map[string][]*stockModel.Stock{
			"AAPL": {
				{
					Ticker:     "AAPL",
					Company:    "Apple Inc.",
					Brokerage:  "Goldman Sachs",
					Action:     "upgraded by",
					RatingFrom: "Hold",
					RatingTo:   "Buy",
					TargetFrom: 150.0,
					TargetTo:   160.0,
					EventTime:  time.Now().AddDate(0, 0, -1),
				},
			},
		}
		mockRepo.On("GetRecentByTickers", ctx, mock.AnythingOfType("time.Time")).Return(testTickerEvents, nil)

		// Mock successful calculator call
		basicRecommendation := &recommendationModel.AggregatedRecommendation{
			Ticker:            "AAPL",
			CompanyName:       "Apple Inc.",
			BasicScore:        0.75,
			Confidence:        0.8,
			Tier:              enums.RECOMMENDATION_TIER_BASIC,
			TotalEvents:       1,
			LatestTargetPrice: 160.0,
		}
		mockCalculator.On("CalculateAggregatedRecommendationFromEvents", ctx, "AAPL", testTickerEvents["AAPL"]).Return(basicRecommendation, nil)

		// Mock cache set error
		mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]*model.AggregatedRecommendation"), mock.AnythingOfType("time.Duration")).Return(errors.New("cache error"))

		// Execute
		result, err := usecase.GetRecommendations(ctx, request)

		// Assert - should succeed even with cache error
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
	})
}

// TestTieredRecommendationUseCase_GetRecommendationByTicker tests single ticker recommendation
func TestTieredRecommendationUseCase_GetRecommendationByTicker(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockStockRepository)
	mockCalculator := new(MockBasicScoringCalculator)
	mockEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockRecommendationCache)
	mockLogger := &MockLogger{}

	usecase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo, mockCalculator, mockEnricher, mockCache, mockLogger,
	)

	ctx := context.Background()
	ticker := "AAPL"
	userTier := enums.TIER_PREMIUM

	t.Run("Cache hit", func(t *testing.T) {
		// Mock cache hit
		cachedRecommendation := &recommendationModel.AggregatedRecommendation{
			Ticker:            "AAPL",
			CompanyName:       "Apple Inc.",
			BasicScore:        0.85,
			Confidence:        0.9,
			Tier:              enums.RECOMMENDATION_TIER_PREMIUM,
			TotalEvents:       3,
			LatestTargetPrice: 180.0,
		}

		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*model.AggregatedRecommendation")).
			Run(func(args mock.Arguments) {
				dest := args.Get(1).(*recommendationModel.AggregatedRecommendation)
				*dest = *cachedRecommendation
			}).Return(nil)

		// Execute
		result, err := usecase.GetRecommendationByTicker(ctx, ticker, userTier)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "AAPL", result.Ticker)
		assert.Equal(t, enums.RECOMMENDATION_TIER_PREMIUM, result.Tier)

		// Verify calculator was not called due to cache hit
		mockCalculator.AssertNotCalled(t, "CalculateAggregatedRecommendation")
	})

	t.Run("Cache miss - fresh calculation", func(t *testing.T) {
		// Mock cache miss
		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*model.AggregatedRecommendation")).Return(errors.New("cache miss"))

		// Mock calculator
		basicRecommendation := &recommendationModel.AggregatedRecommendation{
			Ticker:            "AAPL",
			CompanyName:       "Apple Inc.",
			BasicScore:        0.75,
			Confidence:        0.8,
			Tier:              enums.RECOMMENDATION_TIER_BASIC,
			TotalEvents:       1,
			LatestTargetPrice: 160.0,
		}
		mockCalculator.On("CalculateAggregatedRecommendation", ctx, ticker).Return(basicRecommendation, nil)

		// Mock enricher for premium tier
		enrichedRecommendation := &recommendationModel.AggregatedRecommendation{
			Ticker:            "AAPL",
			CompanyName:       "Apple Inc.",
			BasicScore:        0.85,
			Confidence:        0.9,
			Tier:              enums.RECOMMENDATION_TIER_PREMIUM,
			TotalEvents:       1,
			LatestTargetPrice: 160.0,
		}
		mockEnricher.On("EnrichRecommendation", ctx, basicRecommendation).Return(enrichedRecommendation, nil)

		// Mock cache set
		mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*model.AggregatedRecommendation"), mock.AnythingOfType("time.Duration")).Return(nil)

		// Execute
		result, err := usecase.GetRecommendationByTicker(ctx, ticker, userTier)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "AAPL", result.Ticker)
		assert.Equal(t, enums.RECOMMENDATION_TIER_PREMIUM, result.Tier)
		assert.Equal(t, 0.85, result.BasicScore)
	})
}

// TestTieredRecommendationUseCase_TierLimits tests tier-based limits
func TestTieredRecommendationUseCase_TierLimits(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockStockRepository)
	mockCalculator := new(MockBasicScoringCalculator)
	mockEnricher := new(MockExternalDataEnricher)
	mockCache := new(MockRecommendationCache)
	mockLogger := &MockLogger{}

	usecase := recommendationUsecase.NewTieredRecommendationUseCase(
		mockRepo, mockCalculator, mockEnricher, mockCache, mockLogger,
	)

	ctx := context.Background()

	testCases := []struct {
		name      string
		userTier  enums.UserTier
		requested int
		expected  int
	}{
		{"Guest tier limit", enums.TIER_GUEST, 50, 10},
		{"Basic tier limit", enums.TIER_BASIC, 50, 25},
		{"Premium tier limit", enums.TIER_PREMIUM, 50, 50},
		{"Guest tier within limit", enums.TIER_GUEST, 5, 5},
		{"Basic tier within limit", enums.TIER_BASIC, 15, 15},
		{"Premium tier within limit", enums.TIER_PREMIUM, 30, 30},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := recommendationUsecase.RecommendationRequest{
				UserTier: tc.userTier,
				Limit:    tc.requested,
				Filters:  recommendationUsecase.RecommendationFilters{},
			}

			// Mock cache miss
			mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*[]*model.AggregatedRecommendation")).Return(errors.New("cache miss"))

			// Mock empty stock repository to test limits
			mockRepo.On("GetRecentByTickers", ctx, mock.AnythingOfType("time.Time")).Return(map[string][]*stockModel.Stock{}, nil)

			// Mock cache set
			mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]*model.AggregatedRecommendation"), mock.AnythingOfType("time.Duration")).Return(nil)

			// Execute
			result, err := usecase.GetRecommendations(ctx, request)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.userTier, result.Meta.UserTier)
			assert.Len(t, result.Data, 0) // Empty because no stock data

			// Verify cache key was built with correct limit
			mockCache.AssertExpectations(t)
		})
	}
}