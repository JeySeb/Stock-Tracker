package validation

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/recommendation/validation"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/internal/domain/stocks/repositories"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
)

// TestLogger for testing - simple implementation
type TestLogger struct{}

func (l *TestLogger) Info(msg string, keyvals ...interface{})  {}
func (l *TestLogger) Error(msg string, keyvals ...interface{}) {}
func (l *TestLogger) Warn(msg string, keyvals ...interface{})  {}
func (l *TestLogger) Debug(msg string, keyvals ...interface{}) {}

// MockStockRepository for testing - implements full interface
type MockStockRepository struct {
	mock.Mock
}

// CRUD operations
func (m *MockStockRepository) Create(ctx context.Context, stock *stockModel.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) GetByID(ctx context.Context, id uuid.UUID) (*stockModel.Stock, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) Update(ctx context.Context, stock *stockModel.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Query operations
func (m *MockStockRepository) GetByTicker(ctx context.Context, ticker string) ([]*stockModel.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).([]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetLatestByTicker(ctx context.Context, ticker string) (*stockModel.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetAll(ctx context.Context, filters stockValidation.StockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*stockModel.Stock), args.Get(1).(*stockValidation.Pagination), args.Error(2)
}

func (m *MockStockRepository) GetAllWithEnhancedFilters(ctx context.Context, filters stockValidation.EnhancedStockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*stockModel.Stock), args.Get(1).(*stockValidation.Pagination), args.Error(2)
}

func (m *MockStockRepository) GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*stockModel.Stock, error) {
	args := m.Called(ctx, since)
	return args.Get(0).(map[string][]*stockModel.Stock), args.Error(1)
}

// Batch operations
func (m *MockStockRepository) BulkCreate(ctx context.Context, stocks []*stockModel.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

func (m *MockStockRepository) BulkUpdate(ctx context.Context, stocks []*stockModel.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

// Analytics queries
func (m *MockStockRepository) GetTopMoversByTarget(ctx context.Context, limit int) ([]*stockModel.Stock, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetUniqueTickersCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStockRepository) GetBrokerageStats(ctx context.Context) ([]repositories.BrokerageStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repositories.BrokerageStats), args.Error(1)
}

func TestDirectionalScoringApproach(t *testing.T) {
	mockRepo := new(MockStockRepository)
	testLogger := &TestLogger{}
	calc := validation.NewBasicScoringCalculator(mockRepo, testLogger)

	testCases := []struct {
		name                   string
		events                 []*stockModel.Stock
		expectedRecommendation enums.RecommendationType
		expecteScoreRange      [2]float64 // [min, max]
		businessRationale      string
	}{
		{
			name: "Strong Sell Consensus - Multiple Downgrades to Sell",
			events: []*stockModel.Stock{
				createStockEvent("AAPL", "Goldman Sachs", "Buy", "Sell", 180.0, 150.0),
				createStockEvent("AAPL", "JP Morgan", "Hold", "Sell", 175.0, 145.0),
				createStockEvent("AAPL", "Morgan Stanley", "Outperform", "Strong Sell", 185.0, 140.0),
				createStockEvent("AAPL", "Citigroup", "Buy", "Sell", 190.0, 155.0),
			},
			expectedRecommendation: enums.RECOMMENDATION_TYPE_STRONG_SELL,
			expecteScoreRange:      [2]float64{0.0, 0.15},
			businessRationale:      "Multiple respected brokers downgrading to sell with high certainty should generate STRONG_SELL",
		},
		{
			name: "Sell Consensus - Mixed but Predominantly Negative",
			events: []*stockModel.Stock{
				createStockEvent("TSLA", "Deutsche Bank", "Hold", "Sell", 250.0, 200.0),
				createStockEvent("TSLA", "Wells Fargo", "Buy", "Underperform", 260.0, 210.0),
				createStockEvent("TSLA", "UBS", "Neutral", "Hold", 255.0, 245.0),
			},
			expectedRecommendation: enums.RECOMMENDATION_TYPE_SELL,
			expecteScoreRange:      [2]float64{0.15, 0.35},
			businessRationale:      "Predominantly negative sentiment should generate SELL recommendation",
		},
		{
			name: "Strong Buy Consensus - Multiple Upgrades",
			events: []*stockModel.Stock{
				createStockEvent("NVDA", "Goldman Sachs", "Hold", "Strong Buy", 400.0, 500.0),
				createStockEvent("NVDA", "JP Morgan", "Buy", "Strong Buy", 420.0, 520.0),
				createStockEvent("NVDA", "Morgan Stanley", "Neutral", "Buy", 410.0, 480.0),
				createStockEvent("NVDA", "Barclays", "Underperform", "Outperform", 390.0, 460.0),
			},
			expectedRecommendation: enums.RECOMMENDATION_TYPE_STRONG_BUY,
			expecteScoreRange:      [2]float64{0.75, 1.0},
			businessRationale:      "Multiple upgrades to strong buy should generate STRONG_BUY",
		},
		{
			name: "Hold Consensus - Neutral Mixed Signals",
			events: []*stockModel.Stock{
				createStockEvent("MSFT", "Goldman Sachs", "Buy", "Hold", 350.0, 345.0),
				createStockEvent("MSFT", "JP Morgan", "Sell", "Hold", 340.0, 350.0),
				createStockEvent("MSFT", "Citigroup", "Hold", "Neutral", 345.0, 348.0),
			},
			expectedRecommendation: enums.RECOMMENDATION_TYPE_HOLD,
			expecteScoreRange:      [2]float64{0.35, 0.55},
			businessRationale:      "Mixed signals converging on neutral should generate HOLD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock to return test events
			mockRepo.On("GetByTicker", mock.Anything, mock.AnythingOfType("string")).
				Return(tc.events, nil).Once()
			mockRepo.On("GetBrokerageStats", mock.Anything).
				Return([]repositories.BrokerageStats{
					{Brokerage: "Goldman Sachs", Count: 100, AvgScore: 0.9},
					{Brokerage: "JP Morgan", Count: 95, AvgScore: 0.9},
					{Brokerage: "Morgan Stanley", Count: 90, AvgScore: 0.85},
				}, nil).Once()

			// Calculate recommendation
			recommendation, err := calc.CalculateAggregatedRecommendation(context.Background(), "TEST")

			assert.NoError(t, err)
			assert.NotNil(t, recommendation)

			// Validate recommendation type
			assert.Equal(t, tc.expectedRecommendation, recommendation.RecommendationType,
				"Expected %s but got %s. Business rationale: %s",
				tc.expectedRecommendation, recommendation.RecommendationType, tc.businessRationale)

			// Validate score range
			assert.GreaterOrEqual(t, recommendation.BasicScore, tc.expecteScoreRange[0],
				"Score %f is below expected minimum %f", recommendation.BasicScore, tc.expecteScoreRange[0])
			assert.LessOrEqual(t, recommendation.BasicScore, tc.expecteScoreRange[1],
				"Score %f is above expected maximum %f", recommendation.BasicScore, tc.expecteScoreRange[1])

			// Log scoring factors for analysis
			t.Logf("Recommendation: %s, Score: %.3f", recommendation.RecommendationType, recommendation.BasicScore)
			for _, factor := range recommendation.ScoringFactors {
				t.Logf("  %s: %.3f (weight: %.2f) - %s",
					factor.Name, factor.Score, factor.Weight, factor.Explanation)
			}
		})
	}
}

func TestDirectionalCertaintyCalculation(t *testing.T) {
	testCases := []struct {
		name                     string
		ratingFrom               string
		ratingTo                 string
		expectedDirectionalScore float64
		expectedStrength         float64
	}{
		{
			name:                     "Strong Buy - High Positive Certainty",
			ratingFrom:               "Hold",
			ratingTo:                 "Strong Buy",
			expectedDirectionalScore: 1.0, // Maximum positive certainty
			expectedStrength:         1.0, // Maximum strength
		},
		{
			name:                     "Strong Sell - High Negative Certainty",
			ratingFrom:               "Buy",
			ratingTo:                 "Strong Sell",
			expectedDirectionalScore: -1.0, // Maximum negative certainty
			expectedStrength:         1.0,  // Maximum strength
		},
		{
			name:                     "Hold - Neutral Certainty",
			ratingFrom:               "Buy",
			ratingTo:                 "Hold",
			expectedDirectionalScore: 0.0, // Neutral
			expectedStrength:         0.0, // No strength at neutral
		},
		{
			name:                     "Buy - Moderate Positive Certainty",
			ratingFrom:               "Hold",
			ratingTo:                 "Buy",
			expectedDirectionalScore: 0.6, // 0.6 because Buy = 0.8, distance from 0.5 = 0.3, normalized = 0.6
			expectedStrength:         0.6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &stockModel.Stock{
				RatingFrom: tc.ratingFrom,
				RatingTo:   tc.ratingTo,
			}

			directionalCertainty := stock.GetDirectionalCertainty()
			strength := stock.GetRecommendationStrength()

			assert.InDelta(t, tc.expectedDirectionalScore, directionalCertainty, 0.01,
				"Expected directional certainty %.2f but got %.2f", tc.expectedDirectionalScore, directionalCertainty)
			assert.InDelta(t, tc.expectedStrength, strength, 0.01,
				"Expected strength %.2f but got %.2f", tc.expectedStrength, strength)
		})
	}
}

// Helper function to create test stock events
func createStockEvent(ticker, brokerage, ratingFrom, ratingTo string, targetFrom, targetTo float64) *stockModel.Stock {
	return &stockModel.Stock{
		ID:         uuid.New(),
		Ticker:     ticker,
		Company:    "Test Company",
		Brokerage:  brokerage,
		Action:     "analyst action",
		RatingFrom: ratingFrom,
		RatingTo:   ratingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		EventTime:  time.Now().AddDate(0, 0, -1), // 1 day ago
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
