package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/domain/valueObjects"
	"stock-tracker/tests/mocks"
)

type StockIngestionUseCaseSuite struct {
	suite.Suite
	stockRepo  *mocks.MockStockRepository
	brokerRepo *mocks.MockBrokerRepository
	apiClient  *mocks.MockStockAPIClient
	logger     *mocks.MockLogger
	useCase    *usecases.StockIngestionUseCase
}

func (suite *StockIngestionUseCaseSuite) SetupTest() {
	suite.stockRepo = &mocks.MockStockRepository{}
	suite.brokerRepo = &mocks.MockBrokerRepository{}
	suite.apiClient = &mocks.MockStockAPIClient{}
	suite.logger = &mocks.MockLogger{}

	suite.useCase = usecases.NewStockIngestionUseCase(
		suite.stockRepo,
		suite.brokerRepo,
		suite.apiClient,
		suite.logger,
	)
}

func (suite *StockIngestionUseCaseSuite) TearDownTest() {
	suite.stockRepo.AssertExpectations(suite.T())
	suite.brokerRepo.AssertExpectations(suite.T())
	suite.apiClient.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func TestStockIngestionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(StockIngestionUseCaseSuite))
}

func (suite *StockIngestionUseCaseSuite) TestIngestStocks_Success() {
	// Arrange
	ctx := context.Background()
	testStocks := []*entities.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple Inc.",
			Brokerage: "Goldman Sachs",
			Action:    "upgraded by",
			EventTime: time.Now(),
		},
		{
			Ticker:    "GOOGL",
			Company:   "Alphabet Inc.",
			Brokerage: "Morgan Stanley",
			Action:    "initiated by",
			EventTime: time.Now(),
		},
	}

	existingBrokers := []*entities.Broker{
		entities.NewBroker("Goldman Sachs", 0.95),
	}

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(testStocks, nil)
	suite.brokerRepo.On("GetAll", ctx).Return(existingBrokers, nil)
	suite.brokerRepo.On("Create", ctx, mock.AnythingOfType("*entities.Broker")).Return(nil)
	suite.stockRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]*entities.Stock")).Return(nil)

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *StockIngestionUseCaseSuite) TestIngestStocks_APIClientError() {
	// Arrange
	ctx := context.Background()
	expectedError := errors.New("API client error")

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return([]*entities.Stock{}, expectedError)

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to fetch stocks")
}

func (suite *StockIngestionUseCaseSuite) TestIngestStocks_EmptyStocks() {
	// Arrange
	ctx := context.Background()
	emptyStocks := []*entities.Stock{}

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(emptyStocks, nil)

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *StockIngestionUseCaseSuite) TestIngestStocks_BrokerRepositoryError() {
	// Arrange
	ctx := context.Background()
	testStocks := []*entities.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple Inc.",
			Brokerage: "Goldman Sachs",
			Action:    "upgraded by",
			EventTime: time.Now(),
		},
	}
	expectedError := errors.New("broker repository error")

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(testStocks, nil)
	suite.brokerRepo.On("GetAll", ctx).Return([]*entities.Broker{}, expectedError)

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to enrich stocks with brokers")
}

func (suite *StockIngestionUseCaseSuite) TestIngestStocks_StockRepositoryError() {
	// Arrange
	ctx := context.Background()
	testStocks := []*entities.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple Inc.",
			Brokerage: "Goldman Sachs",
			Action:    "upgraded by",
			EventTime: time.Now(),
		},
	}
	existingBrokers := []*entities.Broker{
		entities.NewBroker("Goldman Sachs", 0.95),
	}
	expectedError := errors.New("stock repository error")

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(testStocks, nil)
	suite.brokerRepo.On("GetAll", ctx).Return(existingBrokers, nil)
	suite.stockRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]*entities.Stock")).Return(expectedError)

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "error during stock ingestion")
}

func (suite *StockIngestionUseCaseSuite) TestEnrichWithBrokerInfo_NewBroker() {
	// Arrange
	ctx := context.Background()
	testStocks := []*entities.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple Inc.",
			Brokerage: "New Brokerage",
			Action:    "upgraded by",
			EventTime: time.Now(),
		},
	}
	existingBrokers := []*entities.Broker{}

	// Setup expectations - test through IngestStocks which calls the private method
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(testStocks, nil)
	suite.brokerRepo.On("GetAll", ctx).Return(existingBrokers, nil)
	suite.brokerRepo.On("Create", ctx, mock.MatchedBy(func(broker *entities.Broker) bool {
		return broker.Name == "New Brokerage" && broker.CredibilityScore == 0.60
	})).Return(nil)
	suite.stockRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]*entities.Stock")).Return(nil)

	// Act - test the private method through the public interface
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *StockIngestionUseCaseSuite) TestEnrichWithBrokerInfo_ExistingBroker() {
	// Arrange
	ctx := context.Background()
	existingBroker := entities.NewBroker("Goldman Sachs", 0.95)
	testStocks := []*entities.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple Inc.",
			Brokerage: "Goldman Sachs",
			Action:    "upgraded by",
			EventTime: time.Now(),
		},
	}
	existingBrokers := []*entities.Broker{existingBroker}

	// Setup expectations - test through IngestStocks which calls the private method
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(testStocks, nil)
	suite.brokerRepo.On("GetAll", ctx).Return(existingBrokers, nil)
	suite.stockRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]*entities.Stock")).Return(nil)

	// Act - test the private method through the public interface
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *StockIngestionUseCaseSuite) TestGetStats() {
	// Arrange
	ctx := context.Background()

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.stockRepo.On("GetAll", ctx, mock.AnythingOfType("valueObjects.StockFilters")).Return([]*entities.Stock{}, &valueObjects.Pagination{TotalItems: 5}, nil)

	// Act
	stats, err := suite.useCase.GetStats(ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stats)

	// Verify stats structure
	statsMap, ok := stats.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Contains(suite.T(), statsMap, "total_stocks")
	assert.Contains(suite.T(), statsMap, "last_updated")
}

func (suite *StockIngestionUseCaseSuite) TestGetStocks() {
	// Arrange
	ctx := context.Background()
	filters := valueObjects.StockFilters{
		Ticker: "AAPL",
		Limit:  10,
	}

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.stockRepo.On("GetAll", ctx, filters).Return([]*entities.Stock{}, &valueObjects.Pagination{TotalItems: 0}, nil)

	// Act
	stocks, pagination, err := suite.useCase.GetStocks(ctx, filters)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stocks)
	assert.NotNil(suite.T(), pagination)
}

func (suite *StockIngestionUseCaseSuite) TestGetStocksByTicker() {
	// Arrange
	ctx := context.Background()
	ticker := "AAPL"

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.stockRepo.On("GetByTicker", ctx, ticker).Return([]*entities.Stock{}, nil)

	// Act
	stocks, err := suite.useCase.GetStocksByTicker(ctx, ticker)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stocks)
}

// Test concurrent processing
func (suite *StockIngestionUseCaseSuite) TestIngestStocks_ConcurrentProcessing() {
	// Arrange
	ctx := context.Background()

	// Create a large number of stocks to test batching
	var testStocks []*entities.Stock
	for i := 0; i < 250; i++ { // More than one batch
		stock := &entities.Stock{
			Ticker:    fmt.Sprintf("STOCK%d", i),
			Company:   fmt.Sprintf("Company %d", i),
			Brokerage: "Test Brokerage",
			Action:    "upgraded by",
			EventTime: time.Now(),
		}
		testStocks = append(testStocks, stock)
	}

	existingBrokers := []*entities.Broker{
		entities.NewBroker("Test Brokerage", 0.80),
	}

	// Setup expectations
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.apiClient.On("FetchAllStocks", ctx).Return(testStocks, nil)
	suite.brokerRepo.On("GetAll", ctx).Return(existingBrokers, nil)

	// Expect multiple BulkCreate calls for different batches
	suite.stockRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]*entities.Stock")).Return(nil).Times(3) // 250 stocks / 100 batch size = 3 batches

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test context cancellation
func (suite *StockIngestionUseCaseSuite) TestIngestStocks_ContextCancellation() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Setup expectations with context that gets cancelled
	suite.apiClient.On("FetchAllStocks", mock.MatchedBy(func(ctx context.Context) bool {
		// Check if context is cancelled
		return ctx.Err() != nil
	})).Return([]*entities.Stock{}, context.Canceled)

	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	suite.logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	// Act
	err := suite.useCase.IngestStocks(ctx)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to fetch stocks")
}
