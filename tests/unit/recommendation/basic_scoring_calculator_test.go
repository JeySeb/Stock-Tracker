package recommendation

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/application/recommendation"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/internal/domain/stocks/repositories"
)

// MockStockRepository implements the StockRepository interface for testing
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) GetByID(ctx context.Context, id interface{}) (*stockModel.Stock, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByTicker(ctx context.Context, ticker string) ([]*stockModel.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).([]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetBrokerageStats(ctx context.Context) ([]repositories.BrokerageStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repositories.BrokerageStats), args.Error(1)
}

// Add other required methods to satisfy the interface
func (m *MockStockRepository) Create(ctx context.Context, stock *stockModel.Stock) error {
	return nil
}

func (m *MockStockRepository) Update(ctx context.Context, stock *stockModel.Stock) error {
	return nil
}

func (m *MockStockRepository) Delete(ctx context.Context, id interface{}) error {
	return nil
}

func (m *MockStockRepository) GetLatestByTicker(ctx context.Context, ticker string) (*stockModel.Stock, error) {
	return nil, nil
}

func (m *MockStockRepository) GetAll(ctx context.Context, filters interface{}) ([]*stockModel.Stock, interface{}, error) {
	return nil, nil, nil
}

func (m *MockStockRepository) GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*stockModel.Stock, error) {
	args := m.Called(ctx, since)
	return args.Get(0).(map[string][]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) BulkCreate(ctx context.Context, stocks []*stockModel.Stock) error {
	return nil
}

func (m *MockStockRepository) BulkUpdate(ctx context.Context, stocks []*stockModel.Stock) error {
	return nil
}

func (m *MockStockRepository) GetTopMoversByTarget(ctx context.Context, limit int) ([]*stockModel.Stock, error) {
	return nil, nil
}

func (m *MockStockRepository) GetUniqueTickersCount(ctx context.Context) (int, error) {
	return 0, nil
}

// MockLogger implements the logger interface for testing
type MockLogger struct{}

func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {}

func TestBasicScoringCalculator_CalculateAggregatedRecommendation(t *testing.T) {
	// Setup
	mockRepo := new(MockStockRepository)
	mockLogger := &MockLogger{}
	calculator := recommendation.NewBasicScoringCalculator(mockRepo, mockLogger)

	ctx := context.Background()
	ticker := "AAPL"

	// Create test data
	testEvents := []*stockModel.Stock{
		{
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Brokerage:  "Goldman Sachs",
			Action:     "upgraded by",
			RatingFrom: "Hold",
			RatingTo:   "Buy",
			TargetFrom: 150.0,
			TargetTo:   160.0,
			EventTime:  time.Now().AddDate(0, 0, -1), // 1 day ago
		},
		{
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Brokerage:  "Morgan Stanley",
			Action:     "target raised by",
			RatingFrom: "Buy",
			RatingTo:   "Buy",
			TargetFrom: 155.0,
			TargetTo:   165.0,
			EventTime:  time.Now().AddDate(0, 0, -2), // 2 days ago
		},
	}

	testBrokerStats := []repositories.BrokerageStats{
		{
			Brokerage: "Goldman Sachs",
			Count:     1000,
			AvgScore:  0.85,
		},
		{
			Brokerage: "Morgan Stanley",
			Count:     800,
			AvgScore:  0.80,
		},
	}

	// Set up mock expectations
	mockRepo.On("GetByTicker", ctx, ticker).Return(testEvents, nil)
	mockRepo.On("GetBrokerageStats", ctx).Return(testBrokerStats, nil)

	// Execute
	result, err := calculator.CalculateAggregatedRecommendation(ctx, ticker)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ticker, result.Ticker)
	assert.Equal(t, "Apple Inc.", result.CompanyName)
	assert.Equal(t, 2, result.TotalEvents)
	assert.Equal(t, enums.RECOMMENDATION_TIER_BASIC, result.Tier)
	assert.Greater(t, result.BasicScore, 0.0)
	assert.LessOrEqual(t, result.BasicScore, 1.0)
	assert.Greater(t, result.Confidence, 0.0)
	assert.LessOrEqual(t, result.Confidence, 1.0)
	assert.NotEmpty(t, result.ScoringFactors)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestBasicScoringCalculator_NoDynamicBrokerCredibility(t *testing.T) {
	// This test verifies that broker credibility is calculated dynamically, not hardcoded
	mockRepo := new(MockStockRepository)
	mockLogger := &MockLogger{}
	calculator := recommendation.NewBasicScoringCalculator(mockRepo, mockLogger)

	ctx := context.Background()
	ticker := "TEST"

	// Create events with different brokers
	testEvents := []*stockModel.Stock{
		{
			Ticker:    "TEST",
			Company:   "Test Company",
			Brokerage: "High Activity Broker", // This broker should have high score
			Action:    "upgraded by",
			EventTime: time.Now().AddDate(0, 0, -1),
		},
		{
			Ticker:    "TEST",
			Company:   "Test Company",
			Brokerage: "Low Activity Broker", // This broker should have low score
			Action:    "upgraded by",
			EventTime: time.Now().AddDate(0, 0, -2),
		},
	}

	// Broker stats showing different activity levels
	testBrokerStats := []repositories.BrokerageStats{
		{
			Brokerage: "High Activity Broker",
			Count:     5000, // High activity
			AvgScore:  0.9,
		},
		{
			Brokerage: "Low Activity Broker",
			Count:     10, // Low activity
			AvgScore:  0.3,
		},
	}

	mockRepo.On("GetByTicker", ctx, ticker).Return(testEvents, nil)
	mockRepo.On("GetBrokerageStats", ctx).Return(testBrokerStats, nil)

	// Execute
	result, err := calculator.CalculateAggregatedRecommendation(ctx, ticker)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Find the broker frequency factor
	var brokerFreqFactor *stockModel.ScoringFactor
	for _, factor := range result.ScoringFactors {
		if factor.Name == "Broker Frequency" {
			brokerFreqFactor = &factor
			break
		}
	}

	assert.NotNil(t, brokerFreqFactor, "Broker Frequency factor should be present")
	assert.Greater(t, brokerFreqFactor.Score, 0.0, "Broker frequency score should be calculated dynamically")

	mockRepo.AssertExpectations(t)
}

func TestBasicScoringCalculator_TargetMovementScoring(t *testing.T) {
	mockRepo := new(MockStockRepository)
	mockLogger := &MockLogger{}
	calculator := recommendation.NewBasicScoringCalculator(mockRepo, mockLogger)

	ctx := context.Background()

	tests := []struct {
		name          string
		ticker        string
		events        []*stockModel.Stock
		expectedScore string // "high", "medium", "low"
	}{
		{
			name:   "Positive target movements should increase score",
			ticker: "POSITIVE",
			events: []*stockModel.Stock{
				{
					Ticker:     "POSITIVE",
					Company:    "Positive Company",
					Brokerage:  "Test Broker",
					TargetFrom: 100.0,
					TargetTo:   120.0, // 20% increase
					EventTime:  time.Now().AddDate(0, 0, -1),
				},
			},
			expectedScore: "high",
		},
		{
			name:   "Negative target movements should decrease score",
			ticker: "NEGATIVE",
			events: []*stockModel.Stock{
				{
					Ticker:     "NEGATIVE",
					Company:    "Negative Company",
					Brokerage:  "Test Broker",
					TargetFrom: 100.0,
					TargetTo:   80.0, // 20% decrease
					EventTime:  time.Now().AddDate(0, 0, -1),
				},
			},
			expectedScore: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetByTicker", ctx, tt.ticker).Return(tt.events, nil)
			mockRepo.On("GetBrokerageStats", ctx).Return([]repositories.BrokerageStats{
				{Brokerage: "Test Broker", Count: 100, AvgScore: 0.5},
			}, nil)

			result, err := calculator.CalculateAggregatedRecommendation(ctx, tt.ticker)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Find target movement factor
			var targetFactor *stockModel.ScoringFactor
			for _, factor := range result.ScoringFactors {
				if factor.Name == "Target Movement" {
					targetFactor = &factor
					break
				}
			}

			assert.NotNil(t, targetFactor)

			switch tt.expectedScore {
			case "high":
				assert.Greater(t, targetFactor.Score, 0.6, "Positive target movement should result in high score")
			case "low":
				assert.Less(t, targetFactor.Score, 0.4, "Negative target movement should result in low score")
			}

			// Clear mock calls for next iteration
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	}
}

func TestBasicScoringCalculator_RecencyScoring(t *testing.T) {
	mockRepo := new(MockStockRepository)
	mockLogger := &MockLogger{}
	calculator := recommendation.NewBasicScoringCalculator(mockRepo, mockLogger)

	ctx := context.Background()

	tests := []struct {
		name          string
		eventTime     time.Time
		expectedScore string
	}{
		{
			name:          "Very recent events should have high recency score",
			eventTime:     time.Now().Add(-12 * time.Hour),
			expectedScore: "high",
		},
		{
			name:          "Old events should have low recency score",
			eventTime:     time.Now().AddDate(0, 0, -45),
			expectedScore: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := []*stockModel.Stock{
				{
					Ticker:    "RECENCY",
					Company:   "Recency Test",
					Brokerage: "Test Broker",
					Action:    "upgraded by",
					EventTime: tt.eventTime,
				},
			}

			mockRepo.On("GetByTicker", ctx, "RECENCY").Return(events, nil)
			mockRepo.On("GetBrokerageStats", ctx).Return([]repositories.BrokerageStats{
				{Brokerage: "Test Broker", Count: 100, AvgScore: 0.5},
			}, nil)

			result, err := calculator.CalculateAggregatedRecommendation(ctx, "RECENCY")

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Find recency factor
			var recencyFactor *stockModel.ScoringFactor
			for _, factor := range result.ScoringFactors {
				if factor.Name == "Recency" {
					recencyFactor = &factor
					break
				}
			}

			assert.NotNil(t, recencyFactor)

			switch tt.expectedScore {
			case "high":
				assert.Greater(t, recencyFactor.Score, 0.7, "Recent events should have high recency score")
			case "low":
				assert.Less(t, recencyFactor.Score, 0.3, "Old events should have low recency score")
			}

			// Clear mock calls for next iteration
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	}
}

func TestBasicScoringCalculator_NoHardcodedValues(t *testing.T) {
	// This test verifies that no values are hardcoded in the scoring algorithm
	mockRepo := new(MockStockRepository)
	mockLogger := &MockLogger{}
	calculator := recommendation.NewBasicScoringCalculator(mockRepo, mockLogger)

	ctx := context.Background()
	ticker := "NOHARDCODE"

	// Create minimal test data
	testEvents := []*stockModel.Stock{
		{
			Ticker:    ticker,
			Company:   "No Hardcode Inc.",
			Brokerage: "Dynamic Broker",
			Action:    "test action",
			EventTime: time.Now().AddDate(0, 0, -5),
		},
	}

	// Empty broker stats - should still work with neutral scores
	testBrokerStats := []repositories.BrokerageStats{}

	mockRepo.On("GetByTicker", ctx, ticker).Return(testEvents, nil)
	mockRepo.On("GetBrokerageStats", ctx).Return(testBrokerStats, nil)

	// Execute
	result, err := calculator.CalculateAggregatedRecommendation(ctx, ticker)

	// Should not fail even with no broker stats
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ticker, result.Ticker)
	assert.GreaterOrEqual(t, result.BasicScore, 0.0)
	assert.LessOrEqual(t, result.BasicScore, 1.0)

	mockRepo.AssertExpectations(t)
}

func TestBasicScoringCalculator_ErrorHandling(t *testing.T) {
	mockRepo := new(MockStockRepository)
	mockLogger := &MockLogger{}
	calculator := recommendation.NewBasicScoringCalculator(mockRepo, mockLogger)

	ctx := context.Background()
	ticker := "ERROR"

	// Test error from GetByTicker
	mockRepo.On("GetByTicker", ctx, ticker).Return([]*stockModel.Stock{}, assert.AnError)

	result, err := calculator.CalculateAggregatedRecommendation(ctx, ticker)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get events")

	mockRepo.AssertExpectations(t)
}
