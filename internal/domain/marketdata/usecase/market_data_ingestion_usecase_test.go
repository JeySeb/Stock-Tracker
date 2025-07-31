package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	marketdatamodel "stock-tracker/internal/domain/marketdata/model"
	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/infrastructure/external"
	"stock-tracker/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMarketDataRepository is a mock implementation of MarketDataRepository
type MockMarketDataRepository struct {
	mock.Mock
}

func (m *MockMarketDataRepository) SaveMarketData(ctx context.Context, marketData *marketdatamodel.MarketData) error {
	args := m.Called(ctx, marketData)
	return args.Error(0)
}

func (m *MockMarketDataRepository) SaveIngestionLog(ctx context.Context, log *marketdatamodel.MarketDataIngestionLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockMarketDataRepository) UpdateIngestionLog(ctx context.Context, log *marketdatamodel.MarketDataIngestionLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockMarketDataRepository) GetUniqueTickers(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMarketDataRepository) ExistsMarketData(ctx context.Context, ticker string, dataSource marketdatamodel.DataSource, dataTimestamp time.Time) (bool, error) {
	args := m.Called(ctx, ticker, dataSource, dataTimestamp)
	return args.Bool(0), args.Error(1)
}

func (m *MockMarketDataRepository) GetLatestMarketData(ctx context.Context, ticker string) (*marketdatamodel.MarketData, error) {
	args := m.Called(ctx, ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*marketdatamodel.MarketData), args.Error(1)
}

func (m *MockMarketDataRepository) GetMarketDataStats(ctx context.Context, days int) (map[string]interface{}, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockYahooFinanceClient is a mock implementation of YahooFinanceClient
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

func TestNewMarketDataIngestionUseCase(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockRepo, useCase.marketDataRepo)
	assert.Equal(t, mockClient, useCase.yahooClient)
	assert.Equal(t, mockLogger, useCase.logger)
}

func TestMarketDataIngestionUseCase_IngestMarketData_Success(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock repository responses
	tickers := []string{"AAPL", "MSFT", "GOOGL"}
	mockRepo.On("GetUniqueTickers", mock.Anything).Return(tickers, nil)
	mockRepo.On("SaveIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(nil)
	mockRepo.On("UpdateIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(nil)

	// Mock client responses for each ticker
	for _, ticker := range tickers {
		mockRepo.On("ExistsMarketData", mock.Anything, ticker, marketdatamodel.DataSourceYahooFinance, mock.AnythingOfType("time.Time")).Return(false, nil)
		mockRepo.On("SaveMarketData", mock.Anything, mock.Anything).Return(nil)

		externalData := &model.ExternalStockData{
			CurrentPrice:     100.0,
			DayChange:        5.0,
			DayChangePercent: 5.0,
			Volume:           1000000,
			MarketCap:        5000000000,
			LastUpdated:      time.Now(),
		}
		mockClient.On("GetQuote", mock.Anything, ticker).Return(externalData, nil)
	}

	// Execute
	err := useCase.IngestMarketData(context.Background())

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestMarketDataIngestionUseCase_IngestMarketData_NoTickers(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock empty tickers response
	mockRepo.On("GetUniqueTickers", mock.Anything).Return([]string{}, nil)

	// Execute
	err := useCase.IngestMarketData(context.Background())

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockClient.AssertNotCalled(t, "GetQuote")
}

func TestMarketDataIngestionUseCase_IngestMarketData_GetTickersError(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock error response
	mockRepo.On("GetUniqueTickers", mock.Anything).Return([]string{}, errors.New("database error"))

	// Execute
	err := useCase.IngestMarketData(context.Background())

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get unique tickers")
	mockRepo.AssertExpectations(t)
	mockClient.AssertNotCalled(t, "GetQuote")
}

func TestMarketDataIngestionUseCase_IngestMarketData_SaveLogError(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock repository responses
	tickers := []string{"AAPL"}
	mockRepo.On("GetUniqueTickers", mock.Anything).Return(tickers, nil)
	mockRepo.On("SaveIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(errors.New("save error"))

	// Execute
	err := useCase.IngestMarketData(context.Background())

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save initial ingestion log")
	mockRepo.AssertExpectations(t)
	mockClient.AssertNotCalled(t, "GetQuote")
}

func TestMarketDataIngestionUseCase_IngestMarketData_ClientError(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock repository responses
	tickers := []string{"AAPL"}
	mockRepo.On("GetUniqueTickers", mock.Anything).Return(tickers, nil)
	mockRepo.On("SaveIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(nil)
	mockRepo.On("UpdateIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(nil)

	// Mock client error
	mockClient.On("GetQuote", mock.Anything, "AAPL").Return(nil, errors.New("API error"))

	// Execute
	err := useCase.IngestMarketData(context.Background())

	// Assertions
	assert.NoError(t, err) // Should not fail the entire process, just log the error
	mockRepo.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestMarketDataIngestionUseCase_IngestMarketData_DataExists(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock repository responses
	tickers := []string{"AAPL"}
	mockRepo.On("GetUniqueTickers", mock.Anything).Return(tickers, nil)
	mockRepo.On("SaveIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(nil)
	mockRepo.On("UpdateIngestionLog", mock.Anything, mock.AnythingOfType("*model.MarketDataIngestionLog")).Return(nil)

	// Mock data already exists
	mockRepo.On("ExistsMarketData", mock.Anything, "AAPL", marketdatamodel.DataSourceYahooFinance, mock.AnythingOfType("time.Time")).Return(true, nil)

	externalData := &model.ExternalStockData{
		CurrentPrice:     100.0,
		DayChange:        5.0,
		DayChangePercent: 5.0,
		Volume:           1000000,
		MarketCap:        5000000000,
		LastUpdated:      time.Now(),
	}
	mockClient.On("GetQuote", mock.Anything, "AAPL").Return(externalData, nil)

	// Execute
	err := useCase.IngestMarketData(context.Background())

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockClient.AssertExpectations(t)
	// Should not call SaveMarketData since data already exists
	mockRepo.AssertNotCalled(t, "SaveMarketData")
}

func TestMarketDataIngestionUseCase_GetMarketDataStats(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock stats response
	expectedStats := map[string]interface{}{
		"total_records":   int64(100),
		"unique_tickers":  int64(10),
		"yahoo_records":   int64(80),
		"alpha_records":   int64(20),
		"avg_price":       float64(150.5),
		"avg_volatility":  float64(2.5),
	}
	mockRepo.On("GetMarketDataStats", mock.Anything, 7).Return(expectedStats, nil)

	// Execute
	stats, err := useCase.GetMarketDataStats(context.Background(), 7)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedStats, stats)
	mockRepo.AssertExpectations(t)
}

func TestMarketDataIngestionUseCase_GetMarketDataStats_Error(t *testing.T) {
	mockRepo := &MockMarketDataRepository{}
	mockClient := &MockYahooFinanceClient{}
	mockLogger := logger.NewSimpleLogger()

	useCase := NewMarketDataIngestionUseCase(mockRepo, mockClient, mockLogger)

	// Mock error response
	mockRepo.On("GetMarketDataStats", mock.Anything, 7).Return(nil, errors.New("database error"))

	// Execute
	stats, err := useCase.GetMarketDataStats(context.Background(), 7)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
} 