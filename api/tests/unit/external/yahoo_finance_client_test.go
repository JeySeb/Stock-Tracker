package external

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/infrastructure/external"
)

// MockLogger implementation
type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func TestYahooFinanceClient_Constructor(t *testing.T) {
	// Create client
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(mockLogger)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestYahooFinanceClient_Constructor_NilLogger(t *testing.T) {
	// Create client with nil logger
	client := external.NewYahooFinanceClient(nil)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestYahooFinanceClient_Constructor_WithConfig(t *testing.T) {
	// Create client with custom config
	mockLogger := &MockLogger{}
	config := &external.ClientConfig{
		Timeout:        10 * time.Second,
		MaxRetries:     2,
		RetryDelay:     500 * time.Millisecond,
		UserAgent:      "TestAgent/1.0",
		RateLimitDelay: 50 * time.Millisecond,
	}

	client := external.NewYahooFinanceClient(mockLogger, config)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestYahooFinanceClient_GetQuote_InvalidSymbol(t *testing.T) {
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(mockLogger)

	ctx := context.Background()

	// Test empty symbol
	result, err := client.GetQuote(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "symbol validation failed")

	// Test invalid symbol with special characters
	result, err = client.GetQuote(ctx, "AAPL@")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "symbol validation failed")

	// Test symbol too long
	result, err = client.GetQuote(ctx, "VERYLONGSYMBOL")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "symbol validation failed")
}

func TestYahooFinanceClient_GetHistoricalData(t *testing.T) {
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	client := external.NewYahooFinanceClient(mockLogger)

	ctx := context.Background()

	// Test valid request
	result, err := client.GetHistoricalData(ctx, "AAPL", "1mo")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result)) // Currently returns empty slice

	// Test empty period
	result, err = client.GetHistoricalData(ctx, "AAPL", "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "period cannot be empty")

	// Test invalid symbol
	result, err = client.GetHistoricalData(ctx, "", "1mo")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "symbol validation failed")

	mockLogger.AssertExpectations(t)
}

func TestYahooFinanceClient_ContextCancellation(t *testing.T) {
	mockLogger := &MockLogger{}
	// Add Debug expectation for the logger call made during the request
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	// Add Error expectation for when the request fails
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	client := external.NewYahooFinanceClient(mockLogger)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := client.GetQuote(ctx, "AAPL")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestYahooFinanceClient_DefaultConfig(t *testing.T) {
	// Test that default config is properly set
	config := external.DefaultClientConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 15*time.Second, config.Timeout)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1*time.Second, config.RetryDelay)
	assert.Equal(t, "Mozilla/5.0 (compatible; StockTracker/1.0)", config.UserAgent)
	assert.Equal(t, 100*time.Millisecond, config.RateLimitDelay)
}

func TestYahooFinanceClient_ErrorTypes(t *testing.T) {
	// Test that error types are properly defined
	assert.NotNil(t, external.ErrInvalidSymbol)
	assert.NotNil(t, external.ErrAPIQuotaExceeded)
	assert.NotNil(t, external.ErrSymbolNotFound)
	assert.NotNil(t, external.ErrAPIUnavailable)
	assert.NotNil(t, external.ErrInvalidResponse)
}
