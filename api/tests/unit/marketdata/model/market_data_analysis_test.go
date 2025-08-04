package model

import (
	"testing"
	"time"

	marketdatamodel "stock-tracker/internal/domain/marketdata/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewMarketDataAnalysis(t *testing.T) {
	now := time.Now()
	week52High := 150.0
	week52Low := 100.0

	analysis := marketdatamodel.NewMarketDataAnalysis(
		"AAPL",
		125.0,
		5.0,
		4.0,
		1000000,
		&week52High,
		&week52Low,
		now,
		now,
		marketdatamodel.DataQualityGood,
		marketdatamodel.DataSourceYahooFinance,
	)

	assert.NotNil(t, analysis)
	assert.Equal(t, "AAPL", analysis.Ticker)
	assert.Equal(t, 125.0, analysis.CurrentPrice)
	assert.Equal(t, 5.0, analysis.DayChange)
	assert.Equal(t, 4.0, analysis.DayChangePercent)
	assert.Equal(t, int64(1000000), analysis.Volume)
	assert.Equal(t, &week52High, analysis.Week52High)
	assert.Equal(t, &week52Low, analysis.Week52Low)
	assert.Equal(t, now, analysis.DataTimestamp)
	assert.Equal(t, now, analysis.CollectedAt)
	assert.Equal(t, marketdatamodel.DataQualityGood, analysis.DataQuality)
	assert.Equal(t, marketdatamodel.DataSourceYahooFinance, analysis.DataSource)
	assert.NotEqual(t, uuid.Nil, analysis.ID)
}

func TestMarketDataAnalysis_CalculateDerivedFields(t *testing.T) {
	now := time.Now()
	week52High := 150.0
	week52Low := 100.0

	analysis := &marketdatamodel.MarketDataAnalysis{
		ID:               uuid.New(),
		Ticker:           "AAPL",
		CurrentPrice:     125.0,
		DayChange:        5.0,
		DayChangePercent: 4.0,
		Volume:           1000000,
		Week52High:       &week52High,
		Week52Low:        &week52Low,
		DataTimestamp:    now,
		CollectedAt:      now,
		DataQuality:      marketdatamodel.DataQualityGood,
		DataSource:       marketdatamodel.DataSourceYahooFinance,
	}
	analysis.CalculateDerivedFields()

	// Test price position calculation
	expectedPosition := (125.0 - 100.0) / (150.0 - 100.0) // 0.5
	assert.Equal(t, expectedPosition, analysis.PricePosition)

	// Test volatility score
	assert.Equal(t, 4.0, analysis.VolatilityScore)

	// Test trend direction (4% change should be bullish)
	assert.Equal(t, "bullish", analysis.TrendDirection)

	// Test risk level (4% volatility should be medium risk)
	assert.Equal(t, "medium", analysis.RiskLevel)
}

func TestMarketDataAnalysis_DetermineTrendDirection(t *testing.T) {
	tests := []struct {
		name          string
		changePercent float64
		expectedTrend string
	}{
		{"bullish", 3.0, "bullish"},
		{"bearish", -3.0, "bearish"},
		{"neutral", 1.0, "neutral"},
		{"neutral_negative", -1.0, "neutral"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &marketdatamodel.MarketDataAnalysis{
				DayChangePercent: tt.changePercent,
			}
			analysis.CalculateDerivedFields()
			// Note: determineTrendDirection is unexported, so we'll test the public IsBullish/IsBearish methods instead
			switch tt.expectedTrend {
			case "bullish":
				assert.True(t, analysis.IsBullish())
			case "bearish":
				assert.True(t, analysis.IsBearish())
			default:
				assert.False(t, analysis.IsBullish())
				assert.False(t, analysis.IsBearish())
			}
		})
	}
}

func TestMarketDataAnalysis_DetermineRiskLevel(t *testing.T) {
	tests := []struct {
		name             string
		dayChangePercent float64
		pricePosition    float64
		expectedRisk     string
	}{
		{"high_volatility", 6.0, 0.5, "high"},
		{"high_position_upper", 2.0, 0.95, "high"},
		{"high_position_lower", 2.0, 0.05, "high"},
		{"medium_volatility", 3.0, 0.5, "medium"},
		{"low_risk", 1.0, 0.5, "low"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &marketdatamodel.MarketDataAnalysis{
				DayChangePercent: tt.dayChangePercent,
				PricePosition:    tt.pricePosition,
			}
			analysis.CalculateDerivedFields()
			switch tt.expectedRisk {
			case "high":
				assert.True(t, analysis.IsHighRisk())
			case "medium":
				assert.Equal(t, "medium", analysis.RiskLevel)
			default:
				assert.False(t, analysis.IsHighRisk())
			}
		})
	}
}

func TestMarketDataAnalysis_IsBullish(t *testing.T) {
	analysis := &marketdatamodel.MarketDataAnalysis{TrendDirection: "bullish"}
	assert.True(t, analysis.IsBullish())

	analysis.TrendDirection = "bearish"
	assert.False(t, analysis.IsBullish())
}

func TestMarketDataAnalysis_IsBearish(t *testing.T) {
	analysis := &marketdatamodel.MarketDataAnalysis{TrendDirection: "bearish"}
	assert.True(t, analysis.IsBearish())

	analysis.TrendDirection = "bullish"
	assert.False(t, analysis.IsBearish())
}

func TestMarketDataAnalysis_IsHighRisk(t *testing.T) {
	analysis := &marketdatamodel.MarketDataAnalysis{RiskLevel: "high"}
	assert.True(t, analysis.IsHighRisk())

	analysis.RiskLevel = "medium"
	assert.False(t, analysis.IsHighRisk())
}

func TestMarketDataAnalysis_GetPriceTarget(t *testing.T) {
	week52High := 150.0
	week52Low := 100.0

	analysis := &marketdatamodel.MarketDataAnalysis{
		Week52High: &week52High,
		Week52Low:  &week52Low,
	}

	support, resistance := analysis.GetPriceTarget()
	assert.Equal(t, 100.0, support)
	assert.Equal(t, 150.0, resistance)
}

func TestMarketDataFilters_SetDefaults(t *testing.T) {
	filters := &marketdatamodel.MarketDataFilters{}
	filters.SetDefaults()

	assert.Equal(t, 50, filters.Limit)
	assert.Equal(t, 0, filters.Offset)
	assert.Equal(t, "data_timestamp", filters.SortBy)
	assert.Equal(t, "desc", filters.SortOrder)
}

func TestMarketDataFilters_SetDefaults_WithValues(t *testing.T) {
	filters := &marketdatamodel.MarketDataFilters{
		Limit:     100,
		Offset:    10,
		SortBy:    "price",
		SortOrder: "asc",
	}
	filters.SetDefaults()

	// Should not override existing values
	assert.Equal(t, 100, filters.Limit)
	assert.Equal(t, 10, filters.Offset)
	assert.Equal(t, "price", filters.SortBy)
	assert.Equal(t, "asc", filters.SortOrder)
}

func TestAbs(t *testing.T) {
	// Note: abs is unexported, so we'll test it indirectly through the public methods
	// that use it internally
}

func TestMarketDataResponse_Creation(t *testing.T) {
	data := []string{"test1", "test2"}
	pagination := &marketdatamodel.Pagination{
		Total:   2,
		Page:    1,
		PerPage: 10,
	}
	message := "Test message"

	response := &marketdatamodel.MarketDataResponse{
		Data:       data,
		Pagination: pagination,
		Message:    message,
	}

	assert.Equal(t, data, response.Data)
	assert.Equal(t, pagination, response.Pagination)
	assert.Equal(t, message, response.Message)
}

func TestPagination_Calculation(t *testing.T) {
	pagination := &marketdatamodel.Pagination{
		Total:       100,
		Page:        2,
		PerPage:     10,
		TotalPages:  10,
		HasNext:     true,
		HasPrevious: true,
	}

	assert.Equal(t, 100, pagination.Total)
	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 10, pagination.PerPage)
	assert.Equal(t, 10, pagination.TotalPages)
	assert.True(t, pagination.HasNext)
	assert.True(t, pagination.HasPrevious)
}

func TestMetadata_Creation(t *testing.T) {
	now := time.Now()
	filters := map[string]interface{}{
		"ticker": "AAPL",
		"limit":  10,
	}

	metadata := &marketdatamodel.Metadata{
		GeneratedAt: now,
		DataPoints:  100,
		Period:      "1d",
		Filters:     filters,
	}

	assert.Equal(t, now, metadata.GeneratedAt)
	assert.Equal(t, 100, metadata.DataPoints)
	assert.Equal(t, "1d", metadata.Period)
	assert.Equal(t, filters, metadata.Filters)
}
