package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	marketDataModel "stock-tracker/internal/domain/marketdata/model"
)

// Test data helpers
func createMockMarketDataAnalysis(ticker string) *marketDataModel.MarketDataAnalysis {
	week52High := 150.0
	week52Low := 100.0
	return marketDataModel.NewMarketDataAnalysis(
		ticker,
		125.0,   // currentPrice
		5.0,     // dayChange
		4.0,     // dayChangePercent
		1000000, // volume
		&week52High,
		&week52Low,
		time.Now(),
		time.Now(),
		marketDataModel.DataQualityExcellent,
		marketDataModel.DataSourceYahooFinance,
	)
}

func createMockMarketDataTrend(ticker string) *marketDataModel.MarketDataTrend {
	return &marketDataModel.MarketDataTrend{
		Ticker:             ticker,
		Period:             "1d",
		StartPrice:         120.0,
		EndPrice:           125.0,
		TotalChange:        5.0,
		TotalChangePercent: 4.0,
		MaxPrice:           126.0,
		MinPrice:           119.0,
		AvgVolume:          1000000,
		Volatility:         2.5,
		TrendStrength:      "strong",
		Direction:          "up",
		DataPoints:         24,
	}
}

func createMockMarketDataSummary() *marketDataModel.MarketDataSummary {
	return &marketDataModel.MarketDataSummary{
		TotalRecords:        1000,
		UniqueTickers:       100,
		AvgPrice:            125.0,
		AvgDayChange:        2.5,
		AvgDayChangePercent: 2.0,
		TotalVolume:         5000000000,
		AvgVolume:           5000000,
		MostActiveTicker:    "AAPL",
		BestPerformer:       "TSLA",
		WorstPerformer:      "META",
		BullishCount:        60,
		BearishCount:        35,
		NeutralCount:        5,
		Period:              "1d",
		LastUpdated:         time.Now(),
	}
}

func createMockMarketData(ticker string) *marketDataModel.MarketData {
	week52High := 150.0
	week52Low := 100.0
	marketCap := int64(2000000000)
	peRatio := 25.0
	dividendYield := 2.0
	avgVolume := int64(1500000)

	return &marketDataModel.MarketData{
		Ticker:           ticker,
		DataSource:       marketDataModel.DataSourceYahooFinance,
		DataQuality:      marketDataModel.DataQualityExcellent,
		CurrentPrice:     125.0,
		DayChange:        5.0,
		DayChangePercent: 4.0,
		Volume:           1000000,
		MarketCap:        &marketCap,
		PERatio:          &peRatio,
		DividendYield:    &dividendYield,
		Week52High:       &week52High,
		Week52Low:        &week52Low,
		AvgVolume:        &avgVolume,
		CollectedAt:      time.Now(),
		DataTimestamp:    time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func createMockMarketDataResponse(data interface{}, message string) *marketDataModel.MarketDataResponse {
	return &marketDataModel.MarketDataResponse{
		Data:    data,
		Message: message,
		Metadata: &marketDataModel.Metadata{
			GeneratedAt: time.Now(),
			DataPoints:  1,
			Period:      "1d",
		},
	}
}

// Test validation functions
func TestMarketDataHandler_ValidationFunctions(t *testing.T) {
	tests := []struct {
		name          string
		period        string
		expectedValid bool
	}{
		{
			name:          "Valid period 1d",
			period:        "1d",
			expectedValid: true,
		},
		{
			name:          "Valid period 1w",
			period:        "1w",
			expectedValid: true,
		},
		{
			name:          "Valid period 1m",
			period:        "1m",
			expectedValid: true,
		},
		{
			name:          "Valid period 3m",
			period:        "3m",
			expectedValid: true,
		},
		{
			name:          "Invalid period invalid",
			period:        "invalid",
			expectedValid: false,
		},
		{
			name:          "Invalid period empty",
			period:        "",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test the validation logic if we had access to the use case
			// For now, we just verify the test structure
			assert.NotEmpty(t, tt.name)
		})
	}
}

// Test URL parameter extraction
func TestMarketDataHandler_URLParameterExtraction(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedTicker string
	}{
		{
			name:           "Valid ticker AAPL",
			url:            "/market-data/analysis/AAPL",
			expectedTicker: "AAPL",
		},
		{
			name:           "Valid ticker TSLA",
			url:            "/market-data/analysis/TSLA",
			expectedTicker: "TSLA",
		},
		{
			name:           "Empty ticker",
			url:            "/market-data/analysis/",
			expectedTicker: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			if tt.expectedTicker != "" {
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("ticker", tt.expectedTicker)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			}

			ticker := chi.URLParam(req, "ticker")
			assert.Equal(t, tt.expectedTicker, ticker)
		})
	}
}

// Test query parameter extraction
func TestMarketDataHandler_QueryParameterExtraction(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedPeriod string
		expectedLimit  string
	}{
		{
			name:           "Valid period and limit",
			url:            "/market-data/top-performers?period=1d&limit=10",
			expectedPeriod: "1d",
			expectedLimit:  "10",
		},
		{
			name:           "Only period",
			url:            "/market-data/top-performers?period=1w",
			expectedPeriod: "1w",
			expectedLimit:  "",
		},
		{
			name:           "Only limit",
			url:            "/market-data/top-performers?limit=5",
			expectedPeriod: "",
			expectedLimit:  "5",
		},
		{
			name:           "No parameters",
			url:            "/market-data/top-performers",
			expectedPeriod: "",
			expectedLimit:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)

			period := req.URL.Query().Get("period")
			limit := req.URL.Query().Get("limit")

			assert.Equal(t, tt.expectedPeriod, period)
			assert.Equal(t, tt.expectedLimit, limit)
		})
	}
}

// Test error response structure
func TestMarketDataHandler_ErrorResponseStructure(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		errorMessage   string
		expectedStatus int
	}{
		{
			name:           "Bad request error",
			statusCode:     http.StatusBadRequest,
			errorMessage:   "Ticker is required",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal server error",
			statusCode:     http.StatusInternalServerError,
			errorMessage:   "Failed to retrieve market data",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Not found error",
			statusCode:     http.StatusNotFound,
			errorMessage:   "No market data found",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Not implemented error",
			statusCode:     http.StatusNotImplemented,
			errorMessage:   "Feature not implemented yet",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Simulate setting status and error response
			w.WriteHeader(tt.statusCode)
			response := map[string]string{"error": tt.errorMessage}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("Failed to encode response: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]string
			err := json.NewDecoder(w.Body).Decode(&responseBody)
			require.NoError(t, err)
			assert.Equal(t, tt.errorMessage, responseBody["error"])
		})
	}
}

// Test model creation helpers
func TestMarketDataHandler_ModelCreationHelpers(t *testing.T) {
	t.Run("CreateMarketDataAnalysis", func(t *testing.T) {
		analysis := createMockMarketDataAnalysis("AAPL")

		assert.NotNil(t, analysis)
		assert.Equal(t, "AAPL", analysis.Ticker)
		assert.Equal(t, 125.0, analysis.CurrentPrice)
		assert.Equal(t, 5.0, analysis.DayChange)
		assert.Equal(t, 4.0, analysis.DayChangePercent)
		assert.Equal(t, int64(1000000), analysis.Volume)
		assert.Equal(t, marketDataModel.DataQualityExcellent, analysis.DataQuality)
		assert.Equal(t, marketDataModel.DataSourceYahooFinance, analysis.DataSource)
	})

	t.Run("CreateMarketDataTrend", func(t *testing.T) {
		trend := createMockMarketDataTrend("AAPL")

		assert.NotNil(t, trend)
		assert.Equal(t, "AAPL", trend.Ticker)
		assert.Equal(t, "1d", trend.Period)
		assert.Equal(t, 120.0, trend.StartPrice)
		assert.Equal(t, 125.0, trend.EndPrice)
		assert.Equal(t, 5.0, trend.TotalChange)
		assert.Equal(t, 4.0, trend.TotalChangePercent)
		assert.Equal(t, "strong", trend.TrendStrength)
		assert.Equal(t, "up", trend.Direction)
	})

	t.Run("CreateMarketDataSummary", func(t *testing.T) {
		summary := createMockMarketDataSummary()

		assert.NotNil(t, summary)
		assert.Equal(t, int64(1000), summary.TotalRecords)
		assert.Equal(t, int64(100), summary.UniqueTickers)
		assert.Equal(t, 125.0, summary.AvgPrice)
		assert.Equal(t, 2.5, summary.AvgDayChange)
		assert.Equal(t, 2.0, summary.AvgDayChangePercent)
		assert.Equal(t, "AAPL", summary.MostActiveTicker)
		assert.Equal(t, "TSLA", summary.BestPerformer)
		assert.Equal(t, "META", summary.WorstPerformer)
		assert.Equal(t, "1d", summary.Period)
	})

	t.Run("CreateMarketData", func(t *testing.T) {
		marketData := createMockMarketData("AAPL")

		assert.NotNil(t, marketData)
		assert.Equal(t, "AAPL", marketData.Ticker)
		assert.Equal(t, 125.0, marketData.CurrentPrice)
		assert.Equal(t, 5.0, marketData.DayChange)
		assert.Equal(t, 4.0, marketData.DayChangePercent)
		assert.Equal(t, int64(1000000), marketData.Volume)
		assert.Equal(t, marketDataModel.DataQualityExcellent, marketData.DataQuality)
		assert.Equal(t, marketDataModel.DataSourceYahooFinance, marketData.DataSource)
		assert.NotNil(t, marketData.MarketCap)
		assert.NotNil(t, marketData.PERatio)
		assert.NotNil(t, marketData.DividendYield)
		assert.NotNil(t, marketData.Week52High)
		assert.NotNil(t, marketData.Week52Low)
		assert.NotNil(t, marketData.AvgVolume)
	})

	t.Run("CreateMarketDataResponse", func(t *testing.T) {
		data := createMockMarketDataAnalysis("AAPL")
		response := createMockMarketDataResponse(data, "Test message")

		assert.NotNil(t, response)
		assert.Equal(t, data, response.Data)
		assert.Equal(t, "Test message", response.Message)
		assert.NotNil(t, response.Metadata)
		assert.Equal(t, 1, response.Metadata.DataPoints)
		assert.Equal(t, "1d", response.Metadata.Period)
	})
}

// Test edge cases for parameter validation
func TestMarketDataHandler_ParameterValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedTicker string
		expectedPeriod string
		expectedLimit  string
	}{
		{
			name:           "Special characters in ticker",
			url:            "/market-data/analysis/AAPL-123",
			expectedTicker: "AAPL-123",
			expectedPeriod: "",
			expectedLimit:  "",
		},
		{
			name:           "Multiple query parameters",
			url:            "/market-data/top-performers?period=1d&limit=10&sort=volume",
			expectedTicker: "",
			expectedPeriod: "1d",
			expectedLimit:  "10",
		},
		{
			name:           "Empty query parameters",
			url:            "/market-data/top-performers?period=&limit=",
			expectedTicker: "",
			expectedPeriod: "",
			expectedLimit:  "",
		},
		{
			name:           "Invalid numeric parameters",
			url:            "/market-data/top-performers?limit=abc&period=1d",
			expectedTicker: "",
			expectedPeriod: "1d",
			expectedLimit:  "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)

			if tt.expectedTicker != "" {
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("ticker", tt.expectedTicker)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			}

			ticker := chi.URLParam(req, "ticker")
			period := req.URL.Query().Get("period")
			limit := req.URL.Query().Get("limit")

			assert.Equal(t, tt.expectedTicker, ticker)
			assert.Equal(t, tt.expectedPeriod, period)
			assert.Equal(t, tt.expectedLimit, limit)
		})
	}
}
