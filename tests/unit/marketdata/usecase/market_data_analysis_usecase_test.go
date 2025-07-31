package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
	"stock-tracker/internal/domain/marketdata/usecase"
	"stock-tracker/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMarketDataAnalysisRepository is a mock implementation of the repository
type MockMarketDataAnalysisRepository struct {
	mock.Mock
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataByTicker(ctx context.Context, ticker string, filters *model.MarketDataFilters) ([]*model.MarketData, error) {
	args := m.Called(ctx, ticker, filters)
	return args.Get(0).([]*model.MarketData), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetLatestMarketDataByTicker(ctx context.Context, ticker string) (*model.MarketData, error) {
	args := m.Called(ctx, ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MarketData), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataByDateRange(ctx context.Context, startDate, endDate time.Time, filters *model.MarketDataFilters) ([]*model.MarketData, error) {
	args := m.Called(ctx, startDate, endDate, filters)
	return args.Get(0).([]*model.MarketData), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataAnalysis(ctx context.Context, ticker string) (*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataTrend(ctx context.Context, ticker string, period string) (*model.MarketDataTrend, error) {
	args := m.Called(ctx, ticker, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MarketDataTrend), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataComparison(ctx context.Context, ticker1, ticker2 string, date time.Time) (*model.MarketDataComparison, error) {
	args := m.Called(ctx, ticker1, ticker2, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MarketDataComparison), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataSummary(ctx context.Context, period string) (*model.MarketDataSummary, error) {
	args := m.Called(ctx, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MarketDataSummary), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetTopPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetWorstPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMostVolatile(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMostActive(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetCorrelatedTickers(ctx context.Context, ticker string, threshold float64, period string) ([]string, error) {
	args := m.Called(ctx, ticker, threshold, period)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetPriceBreakouts(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetVolumeSurges(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetHighRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetLowRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetRiskDistribution(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetTrendingTickers(ctx context.Context, direction string, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, direction, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetTrendStrength(ctx context.Context, ticker string, period string) (string, error) {
	args := m.Called(ctx, ticker, period)
	return args.String(0), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetTrendReversal(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	args := m.Called(ctx, limit, period)
	return args.Get(0).([]*model.MarketDataAnalysis), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetPriceStatistics(ctx context.Context, ticker string, period string) (map[string]float64, error) {
	args := m.Called(ctx, ticker, period)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetVolumeStatistics(ctx context.Context, ticker string, period string) (map[string]int64, error) {
	args := m.Called(ctx, ticker, period)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetVolatilityStatistics(ctx context.Context, ticker string, period string) (map[string]float64, error) {
	args := m.Called(ctx, ticker, period)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketAlerts(ctx context.Context, alertType string, severity string, limit int) ([]*model.MarketDataAlert, error) {
	args := m.Called(ctx, alertType, severity, limit)
	return args.Get(0).([]*model.MarketDataAlert), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) CreateMarketAlert(ctx context.Context, alert *model.MarketDataAlert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

func (m *MockMarketDataAnalysisRepository) UpdateMarketAlert(ctx context.Context, alert *model.MarketDataAlert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

func (m *MockMarketDataAnalysisRepository) DeleteMarketAlert(ctx context.Context, alertID string) error {
	args := m.Called(ctx, alertID)
	return args.Error(0)
}

func (m *MockMarketDataAnalysisRepository) GetDataQualityStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetDataSourceStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetDataFreshness(ctx context.Context) (map[string]time.Time, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]time.Time), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataWithStockAnalysis(ctx context.Context, ticker string) (map[string]interface{}, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetCorrelationWithBrokerActions(ctx context.Context, ticker string, period string) (map[string]interface{}, error) {
	args := m.Called(ctx, ticker, period)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockMarketDataAnalysisRepository) GetMarketDataImpactOnRecommendations(ctx context.Context, ticker string) (map[string]interface{}, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestMarketDataAnalysisUseCase_GetMarketDataAnalysis_Success(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	expectedAnalysis := &model.MarketDataAnalysis{
		Ticker: "AAPL",
	}

	mockRepo.On("GetMarketDataAnalysis", mock.Anything, "AAPL").Return(expectedAnalysis, nil)

	analysis, err := useCase.GetMarketDataAnalysis(context.Background(), "AAPL")

	assert.NoError(t, err)
	assert.Equal(t, expectedAnalysis, analysis)
	mockRepo.AssertExpectations(t)
}

func TestMarketDataAnalysisUseCase_GetMarketDataAnalysis_Error(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	mockRepo.On("GetMarketDataAnalysis", mock.Anything, "AAPL").Return(nil, errors.New("database error"))

	analysis, err := useCase.GetMarketDataAnalysis(context.Background(), "AAPL")

	assert.Error(t, err)
	assert.Nil(t, analysis)
	assert.Contains(t, err.Error(), "failed to get market data analysis")
	mockRepo.AssertExpectations(t)
}

func TestMarketDataAnalysisUseCase_GetMarketDataTrend_Success(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	expectedTrend := &model.MarketDataTrend{
		Ticker: "AAPL",
		Period: "1w",
	}

	mockRepo.On("GetMarketDataTrend", mock.Anything, "AAPL", "1w").Return(expectedTrend, nil)

	trend, err := useCase.GetMarketDataTrend(context.Background(), "AAPL", "1w")

	assert.NoError(t, err)
	assert.Equal(t, expectedTrend, trend)
	mockRepo.AssertExpectations(t)
}

func TestMarketDataAnalysisUseCase_GetMarketDataSummary_Success(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	expectedSummary := &model.MarketDataSummary{
		TotalRecords: 100,
		Period:       "1d",
	}

	mockRepo.On("GetMarketDataSummary", mock.Anything, "1d").Return(expectedSummary, nil)

	summary, err := useCase.GetMarketDataSummary(context.Background(), "1d")

	assert.NoError(t, err)
	assert.Equal(t, expectedSummary, summary)
	mockRepo.AssertExpectations(t)
}

func TestMarketDataAnalysisUseCase_GetTopPerformers_Success(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	expectedPerformers := []*model.MarketDataAnalysis{
		{Ticker: "AAPL"},
		{Ticker: "MSFT"},
	}

	mockRepo.On("GetTopPerformers", mock.Anything, 10, "1d").Return(expectedPerformers, nil)

	performers, err := useCase.GetTopPerformers(context.Background(), 10, "1d")

	assert.NoError(t, err)
	assert.Equal(t, expectedPerformers, performers)
	mockRepo.AssertExpectations(t)
}

func TestMarketDataAnalysisUseCase_GetWorstPerformers_Success(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	expectedPerformers := []*model.MarketDataAnalysis{
		{Ticker: "TSLA"},
		{Ticker: "NFLX"},
	}

	mockRepo.On("GetWorstPerformers", mock.Anything, 10, "1d").Return(expectedPerformers, nil)

	performers, err := useCase.GetWorstPerformers(context.Background(), 10, "1d")

	assert.NoError(t, err)
	assert.Equal(t, expectedPerformers, performers)
	mockRepo.AssertExpectations(t)
}

func TestMarketDataAnalysisUseCase_ParseFiltersFromRequest(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	params := map[string]string{
		"ticker":          "AAPL",
		"data_source":     "yahoo_finance",
		"data_quality":    "good",
		"limit":           "20",
		"offset":          "10",
		"sort_by":         "price",
		"sort_order":      "asc",
		"start_date":      "2024-01-01",
		"end_date":        "2024-01-31",
		"min_price":       "100.0",
		"max_price":       "200.0",
		"min_change":      "-5.0",
		"max_change":      "5.0",
		"min_volume":      "1000000",
		"max_volume":      "10000000",
		"trend_direction": "bullish",
		"risk_level":      "medium",
	}

	filters := useCase.ParseFiltersFromRequest(params)

	assert.Equal(t, "AAPL", filters.Ticker)
	assert.Equal(t, "yahoo_finance", filters.DataSource)
	assert.Equal(t, "good", filters.DataQuality)
	assert.Equal(t, 20, filters.Limit)
	assert.Equal(t, 10, filters.Offset)
	assert.Equal(t, "price", filters.SortBy)
	assert.Equal(t, "asc", filters.SortOrder)
	assert.Equal(t, "bullish", filters.TrendDirection)
	assert.Equal(t, "medium", filters.RiskLevel)
	assert.NotNil(t, filters.MinPrice)
	assert.NotNil(t, filters.MaxPrice)
	assert.NotNil(t, filters.MinChange)
	assert.NotNil(t, filters.MaxChange)
	assert.NotNil(t, filters.MinVolume)
	assert.NotNil(t, filters.MaxVolume)
}

func TestMarketDataAnalysisUseCase_ValidatePeriod(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	assert.True(t, useCase.ValidatePeriod("1d"))
	assert.True(t, useCase.ValidatePeriod("1w"))
	assert.True(t, useCase.ValidatePeriod("1m"))
	assert.True(t, useCase.ValidatePeriod("3m"))
	assert.False(t, useCase.ValidatePeriod("invalid"))
	assert.False(t, useCase.ValidatePeriod(""))
}

func TestMarketDataAnalysisUseCase_GetDefaultPeriod(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	assert.Equal(t, "1d", useCase.GetDefaultPeriod())
}

func TestMarketDataAnalysisUseCase_ValidateLimit(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	assert.Equal(t, 10, useCase.ValidateLimit(0))
	assert.Equal(t, 10, useCase.ValidateLimit(-5))
	assert.Equal(t, 50, useCase.ValidateLimit(50))
	assert.Equal(t, 100, useCase.ValidateLimit(150))
	assert.Equal(t, 100, useCase.ValidateLimit(200))
}

func TestMarketDataAnalysisUseCase_CreateMarketDataResponse(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	data := []string{"test1", "test2"}
	pagination := &model.Pagination{
		Total:   2,
		Page:    1,
		PerPage: 10,
	}
	message := "Test message"

	response := useCase.CreateMarketDataResponse(data, pagination, message)

	assert.Equal(t, data, response.Data)
	assert.Equal(t, pagination, response.Pagination)
	assert.Equal(t, message, response.Message)
	assert.NotNil(t, response.Metadata)
}

func TestMarketDataAnalysisUseCase_CalculatePagination(t *testing.T) {
	mockRepo := new(MockMarketDataAnalysisRepository)
	mockLogger := logger.NewSimpleLogger()
	useCase := usecase.NewMarketDataAnalysisUseCase(mockRepo, mockLogger)

	pagination := useCase.CalculatePagination(100, 2, 10)

	assert.Equal(t, 100, pagination.Total)
	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 10, pagination.PerPage)
	assert.Equal(t, 10, pagination.TotalPages)
	assert.True(t, pagination.HasNext)
	assert.True(t, pagination.HasPrevious)
}
