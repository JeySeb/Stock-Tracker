package validation

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
