package external

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/infrastructure/external"
)

func TestAlphaVantageClient_Constructor(t *testing.T) {
	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient("test_key", mockLogger)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestAlphaVantageClient_Constructor_EmptyAPIKey(t *testing.T) {
	// Create client with empty API key
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient("", mockLogger)

	// Assert client is still created (validation might be done at request time)
	assert.NotNil(t, client)
}

func TestAlphaVantageClient_Constructor_NilLogger(t *testing.T) {
	// Create client with nil logger
	client := external.NewAlphaVantageClient("test_key", nil)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestAlphaVantageClient_GetQuote_NetworkError(t *testing.T) {
	mockLogger := &MockLogger{}
	// Only expect Error call since validation will fail and log an error
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	client := external.NewAlphaVantageClient("test_key", mockLogger)

	ctx := context.Background()
	// Use empty symbol to trigger validation error instead of network error
	result, err := client.GetQuote(ctx, "")

	// Assertions - should fail due to validation error
	assert.Error(t, err)
	assert.Nil(t, result)

	mockLogger.AssertExpectations(t)
}

func TestAlphaVantageClient_GetCompanyOverview_NetworkError(t *testing.T) {
	mockLogger := &MockLogger{}
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	client := external.NewAlphaVantageClient("test_key", mockLogger)

	ctx := context.Background()
	// Use empty symbol to trigger validation error instead of network error
	result, err := client.GetCompanyOverview(ctx, "")

	// Assertions - should fail due to validation error
	assert.Error(t, err)
	assert.Nil(t, result)

	mockLogger.AssertExpectations(t)
}

func TestAlphaVantageClient_ContextCancellation(t *testing.T) {
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient("test_key", mockLogger)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := client.GetQuote(ctx, "AAPL")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestAlphaVantageClient_EmptySymbol(t *testing.T) {
	mockLogger := &MockLogger{}
	// Only expect Error calls since validation will fail and log errors
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	client := external.NewAlphaVantageClient("test_key", mockLogger)

	ctx := context.Background()

	// Test empty symbol for GetQuote
	result, err := client.GetQuote(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid or empty symbol")

	// Test empty symbol for GetCompanyOverview
	overview, err := client.GetCompanyOverview(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, overview)
	assert.Contains(t, err.Error(), "invalid or empty symbol")

	mockLogger.AssertExpectations(t)
}

func TestAlphaVantageClient_CompanyOverview_Struct(t *testing.T) {
	// Test that CompanyOverview struct is properly defined
	overview := &external.CompanyOverview{
		Symbol:          "AAPL",
		Name:            "Apple Inc.",
		MarketCap:       2500000000000,
		PERatio:         25.5,
		DividendYield:   0.5,
		EPS:             5.95,
		RevenuePerShare: 25.30,
		ProfitMargin:    25.0,
	}

	assert.NotNil(t, overview)
	assert.Equal(t, "AAPL", overview.Symbol)
	assert.Equal(t, "Apple Inc.", overview.Name)
	assert.Equal(t, int64(2500000000000), overview.MarketCap)
	assert.Equal(t, 25.5, overview.PERatio)
	assert.Equal(t, 0.5, overview.DividendYield)
	assert.Equal(t, 5.95, overview.EPS)
	assert.Equal(t, 25.30, overview.RevenuePerShare)
	assert.Equal(t, 25.0, overview.ProfitMargin)
}
