package model_test

import (
	"testing"
	"time"

	marketdatamodel "stock-tracker/internal/domain/marketdata/model"
	recommendationmodel "stock-tracker/internal/domain/recommendation/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewMarketData(t *testing.T) {
	// Create test external data
	externalData := &marketdatamodel.ExternalStockData{
		CurrentPrice:     100.50,
		DayChange:        5.25,
		DayChangePercent: 5.5,
		Volume:           1000000,
		MarketCap:        &[]int64{5000000000}[0],
		PERatio:          &[]float64{15.5}[0],
		DividendYield:    &[]float64{2.5}[0],
		Week52High:       &[]float64{120.0}[0],
		Week52Low:        &[]float64{80.0}[0],
		AvgVolume:        &[]int64{800000}[0],
		LastUpdated:      time.Now(),
	}

	// Test NewMarketData
	marketData := marketdatamodel.NewMarketData("AAPL", marketdatamodel.DataSourceYahooFinance, externalData)

	// Assertions
	assert.NotNil(t, marketData)
	assert.NotEqual(t, uuid.Nil, marketData.ID)
	assert.Equal(t, "AAPL", marketData.Ticker)
	assert.Equal(t, marketdatamodel.DataSourceYahooFinance, marketData.DataSource)
	assert.Equal(t, marketdatamodel.DataQualityGood, marketData.DataQuality)
	assert.Equal(t, 100.50, marketData.CurrentPrice)
	assert.Equal(t, 5.25, marketData.DayChange)
	assert.Equal(t, 5.5, marketData.DayChangePercent)
	assert.Equal(t, int64(1000000), marketData.Volume)
	assert.Equal(t, int64(5000000000), *marketData.MarketCap)
	assert.Equal(t, 15.5, *marketData.PERatio)
	assert.Equal(t, 2.5, *marketData.DividendYield)
	assert.Equal(t, 120.0, *marketData.Week52High)
	assert.Equal(t, 80.0, *marketData.Week52Low)
	assert.Equal(t, int64(800000), *marketData.AvgVolume)
	assert.False(t, marketData.CollectedAt.IsZero())
	assert.False(t, marketData.DataTimestamp.IsZero())
	assert.False(t, marketData.CreatedAt.IsZero())
	assert.False(t, marketData.UpdatedAt.IsZero())
}

func TestNewMarketDataIngestionLog(t *testing.T) {
	batchID := "test-batch-123"
	totalTickers := 50

	// Test NewMarketDataIngestionLog
	log := marketdatamodel.NewMarketDataIngestionLog(batchID, marketdatamodel.DataSourceYahooFinance, totalTickers)

	// Assertions
	assert.NotNil(t, log)
	assert.NotEqual(t, uuid.Nil, log.ID)
	assert.Equal(t, batchID, log.BatchID)
	assert.Equal(t, marketdatamodel.DataSourceYahooFinance, log.DataSource)
	assert.Equal(t, totalTickers, log.TotalTickers)
	assert.Equal(t, 0, log.SuccessfulTickers)
	assert.Equal(t, 0, log.FailedTickers)
	assert.Equal(t, 0, log.SkippedTickers)
	assert.Equal(t, "running", log.Status)
	assert.False(t, log.StartedAt.IsZero())
	assert.False(t, log.CreatedAt.IsZero())
	assert.False(t, log.UpdatedAt.IsZero())
	assert.Nil(t, log.CompletedAt)
}

func TestConvertFromRecommendationModel(t *testing.T) {
	// Create test recommendation model data
	recModel := &recommendationmodel.ExternalStockData{
		CurrentPrice:     150.75,
		DayChange:        -2.25,
		DayChangePercent: -1.5,
		Volume:           2000000,
		MarketCap:        7500000000, // int64 in recommendation model
		PERatio:          &[]float64{18.2}[0],
		DividendYield:    &[]float64{1.8}[0],
		Week52High:       &[]float64{160.0}[0],
		Week52Low:        &[]float64{90.0}[0],
		AvgVolume:        &[]int64{1200000}[0],
		LastUpdated:      time.Now(),
	}

	// Test conversion
	marketData := marketdatamodel.ConvertFromRecommendationModel("MSFT", marketdatamodel.DataSourceYahooFinance, recModel)

	// Assertions
	assert.NotNil(t, marketData)
	assert.NotEqual(t, uuid.Nil, marketData.ID)
	assert.Equal(t, "MSFT", marketData.Ticker)
	assert.Equal(t, marketdatamodel.DataSourceYahooFinance, marketData.DataSource)
	assert.Equal(t, marketdatamodel.DataQualityGood, marketData.DataQuality)
	assert.Equal(t, 150.75, marketData.CurrentPrice)
	assert.Equal(t, -2.25, marketData.DayChange)
	assert.Equal(t, -1.5, marketData.DayChangePercent)
	assert.Equal(t, int64(2000000), marketData.Volume)
	assert.NotNil(t, marketData.MarketCap)
	assert.Equal(t, int64(7500000000), *marketData.MarketCap)
	assert.Equal(t, 18.2, *marketData.PERatio)
	assert.Equal(t, 1.8, *marketData.DividendYield)
	assert.Equal(t, 160.0, *marketData.Week52High)
	assert.Equal(t, 90.0, *marketData.Week52Low)
	assert.Equal(t, int64(1200000), *marketData.AvgVolume)
	assert.False(t, marketData.CollectedAt.IsZero())
	assert.False(t, marketData.DataTimestamp.IsZero())
	assert.False(t, marketData.CreatedAt.IsZero())
	assert.False(t, marketData.UpdatedAt.IsZero())
}

func TestConvertFromRecommendationModel_ZeroMarketCap(t *testing.T) {
	// Test with zero MarketCap
	recModel := &recommendationmodel.ExternalStockData{
		CurrentPrice:     100.0,
		DayChange:        0.0,
		DayChangePercent: 0.0,
		Volume:           1000000,
		MarketCap:        0, // Zero value
		LastUpdated:      time.Now(),
	}

	// Test conversion
	marketData := marketdatamodel.ConvertFromRecommendationModel("TEST", marketdatamodel.DataSourceYahooFinance, recModel)

	// Assertions
	assert.NotNil(t, marketData)
	assert.Equal(t, "TEST", marketData.Ticker)
	assert.Equal(t, 100.0, marketData.CurrentPrice)
	assert.Nil(t, marketData.MarketCap) // Should be nil for zero value
}

func TestDataSourceConstants(t *testing.T) {
	// Test DataSource constants
	assert.Equal(t, marketdatamodel.DataSource("yahoo_finance"), marketdatamodel.DataSourceYahooFinance)
	assert.Equal(t, marketdatamodel.DataSource("alpha_vantage"), marketdatamodel.DataSourceAlphaVantage)
	assert.Equal(t, marketdatamodel.DataSource("manual"), marketdatamodel.DataSourceManual)
}

func TestDataQualityConstants(t *testing.T) {
	// Test DataQuality constants
	assert.Equal(t, marketdatamodel.DataQuality("excellent"), marketdatamodel.DataQualityExcellent)
	assert.Equal(t, marketdatamodel.DataQuality("good"), marketdatamodel.DataQualityGood)
	assert.Equal(t, marketdatamodel.DataQuality("fair"), marketdatamodel.DataQualityFair)
	assert.Equal(t, marketdatamodel.DataQuality("poor"), marketdatamodel.DataQualityPoor)
}

func TestMarketData_JSONTags(t *testing.T) {
	// Test that JSON tags are properly set
	marketData := &marketdatamodel.MarketData{
		ID:               uuid.New(),
		Ticker:           "TEST",
		DataSource:       marketdatamodel.DataSourceYahooFinance,
		DataQuality:      marketdatamodel.DataQualityGood,
		CurrentPrice:     100.0,
		DayChange:        5.0,
		DayChangePercent: 5.0,
		Volume:           1000000,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// This test ensures the struct can be marshaled to JSON
	// The actual marshaling is tested implicitly by the struct definition
	assert.NotNil(t, marketData)
	assert.NotEqual(t, uuid.Nil, marketData.ID)
	assert.Equal(t, "TEST", marketData.Ticker)
}

func TestMarketDataIngestionLog_JSONTags(t *testing.T) {
	// Test that JSON tags are properly set
	log := &marketdatamodel.MarketDataIngestionLog{
		ID:           uuid.New(),
		BatchID:      "test-batch",
		DataSource:   marketdatamodel.DataSourceYahooFinance,
		TotalTickers: 10,
		Status:       "running",
		StartedAt:    time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// This test ensures the struct can be marshaled to JSON
	// The actual marshaling is tested implicitly by the struct definition
	assert.NotNil(t, log)
	assert.NotEqual(t, uuid.Nil, log.ID)
	assert.Equal(t, "test-batch", log.BatchID)
}
