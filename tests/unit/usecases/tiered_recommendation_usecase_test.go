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

// MockCache implementation
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

// MockLogger implementation
type MockLogger struct{}

func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {}

func TestTieredRecommendationUseCase_Initialization(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)

	// Test that we can create the use case (constructor works)
	useCase := recommendationUsecase.NewTieredRecommendationUseCase(mockRepo, mockBasicCalculator, mockExternalEnricher, mockCache, mockLogger)

	// Assert
	assert.NotNil(t, useCase)
}

func TestTieredRecommendationUseCase_CacheHit(t *testing.T) {
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

	// Create sample cached result - just the recommendations array
	cachedRecommendations := []*recommendationModel.AggregatedRecommendation{
		{
			Ticker:      "AAPL",
			CompanyName: "Apple Inc.",
			BasicScore:  0.8,
			Tier:        enums.RECOMMENDATION_TIER_BASIC,
		},
	}

	// Mock cache hit - simulate the cached recommendations being returned
	mockCache.On("Get", mock.Anything, "recommendations:basic:10", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Simulate cache hit by setting the destination
		dest := args.Get(2).(*[]*recommendationModel.AggregatedRecommendation)
		*dest = cachedRecommendations
	})

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "AAPL", result.Data[0].Ticker)

	// Verify no repository calls were made (cache hit)
	mockRepo.AssertNotCalled(t, "GetRecentByTickers")
}

func TestTieredRecommendationUseCase_CacheMiss_NilServices(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MockStockRepository)
	mockCache := new(MockCache)
	mockLogger := &MockLogger{}
	mockBasicCalculator := new(MockBasicScoringCalculator)
	mockExternalEnricher := new(MockExternalDataEnricher)

	useCase := recommendationUsecase.NewTieredRecommendationUseCase(mockRepo, mockBasicCalculator, mockExternalEnricher, mockCache, mockLogger)

	request := recommendationUsecase.RecommendationRequest{
		UserTier: enums.TIER_GUEST,
		Limit:    5,
		Filters:  recommendationUsecase.RecommendationFilters{},
	}

	// Mock cache miss (guest user with limit 5 should generate cache key "recommendations:guest:5")
	mockCache.On("Get", mock.Anything, "recommendations:guest:5", mock.Anything).Return(cache.ErrCacheMiss)

	// Mock empty ticker events to avoid calling the nil services
	emptyTickerEvents := map[string][]*stockModel.Stock{}
	mockRepo.On("GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time")).Return(emptyTickerEvents, nil)

	// Mock cache set (will be called even if no recommendations are generated)
	mockCache.On("Set", mock.Anything, "recommendations:guest:5", mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := useCase.GetRecommendations(context.Background(), request)

	// Assert - should succeed with empty results since no tickers were provided
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 0) // No recommendations since no ticker events
	assert.Equal(t, enums.TIER_GUEST, result.Meta.UserTier)
	assert.False(t, result.Meta.CacheHit)

	// Verify cache was checked and repository was called
	mockCache.AssertCalled(t, "Get", mock.Anything, "recommendations:guest:5", mock.Anything)
	mockRepo.AssertCalled(t, "GetRecentByTickers", mock.Anything, mock.AnythingOfType("time.Time"))
}
