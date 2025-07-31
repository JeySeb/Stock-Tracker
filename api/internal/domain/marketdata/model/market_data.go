package model

import (
	"time"

	recommendationmodel "stock-tracker/internal/domain/recommendation/model"

	"github.com/google/uuid"
)

// DataSource represents the source of market data
type DataSource string

const (
	DataSourceYahooFinance DataSource = "yahoo_finance"
	DataSourceAlphaVantage DataSource = "alpha_vantage"
	DataSourceManual       DataSource = "manual"
)

// DataQuality represents the quality level of the data
type DataQuality string

const (
	DataQualityExcellent DataQuality = "excellent"
	DataQualityGood      DataQuality = "good"
	DataQualityFair      DataQuality = "fair"
	DataQualityPoor      DataQuality = "poor"
)

// MarketData represents market data from external sources
type MarketData struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	Ticker      string      `json:"ticker" db:"ticker"`
	DataSource  DataSource  `json:"data_source" db:"data_source"`
	DataQuality DataQuality `json:"data_quality" db:"data_quality"`

	// Price data
	CurrentPrice     float64 `json:"current_price" db:"current_price"`
	DayChange        float64 `json:"day_change" db:"day_change"`
	DayChangePercent float64 `json:"day_change_percent" db:"day_change_percent"`

	// Volume and market metrics
	Volume    int64  `json:"volume" db:"volume"`
	MarketCap *int64 `json:"market_cap,omitempty" db:"market_cap"`

	// Fundamental ratios (nullable)
	PERatio       *float64 `json:"pe_ratio,omitempty" db:"pe_ratio"`
	DividendYield *float64 `json:"dividend_yield,omitempty" db:"dividend_yield"`

	// Technical levels
	Week52High *float64 `json:"week_52_high,omitempty" db:"week_52_high"`
	Week52Low  *float64 `json:"week_52_low,omitempty" db:"week_52_low"`
	AvgVolume  *int64   `json:"avg_volume,omitempty" db:"avg_volume"`

	// Metadata
	CollectedAt   time.Time `json:"collected_at" db:"collected_at"`
	DataTimestamp time.Time `json:"data_timestamp" db:"data_timestamp"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// MarketDataIngestionLog represents the ingestion process log
type MarketDataIngestionLog struct {
	ID                uuid.UUID   `json:"id" db:"id"`
	BatchID           string      `json:"batch_id" db:"batch_id"`
	DataSource        DataSource  `json:"data_source" db:"data_source"`
	TotalTickers      int         `json:"total_tickers" db:"total_tickers"`
	SuccessfulTickers int         `json:"successful_tickers" db:"successful_tickers"`
	FailedTickers     int         `json:"failed_tickers" db:"failed_tickers"`
	SkippedTickers    int         `json:"skipped_tickers" db:"skipped_tickers"`
	Status            string      `json:"status" db:"status"`
	ErrorDetails      interface{} `json:"error_details,omitempty" db:"error_details"`
	StartedAt         time.Time   `json:"started_at" db:"started_at"`
	CompletedAt       *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

// NewMarketData creates a new market data instance
// DEPRECATED: Use ConvertFromRecommendationModel instead
func NewMarketData(ticker string, dataSource DataSource, externalData *ExternalStockData) *MarketData {
	now := time.Now()

	return &MarketData{
		ID:               uuid.New(),
		Ticker:           ticker,
		DataSource:       dataSource,
		DataQuality:      DataQualityGood, // Default quality
		CurrentPrice:     externalData.CurrentPrice,
		DayChange:        externalData.DayChange,
		DayChangePercent: externalData.DayChangePercent,
		Volume:           externalData.Volume,
		MarketCap:        externalData.MarketCap,
		PERatio:          externalData.PERatio,
		DividendYield:    externalData.DividendYield,
		Week52High:       externalData.Week52High,
		Week52Low:        externalData.Week52Low,
		AvgVolume:        externalData.AvgVolume,
		CollectedAt:      now,
		DataTimestamp:    externalData.LastUpdated,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// NewMarketDataIngestionLog creates a new ingestion log
func NewMarketDataIngestionLog(batchID string, dataSource DataSource, totalTickers int) *MarketDataIngestionLog {
	now := time.Now()

	return &MarketDataIngestionLog{
		ID:           uuid.New(),
		BatchID:      batchID,
		DataSource:   dataSource,
		TotalTickers: totalTickers,
		Status:       "running",
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// ConvertFromRecommendationModel converts from recommendation model to market data model
func ConvertFromRecommendationModel(ticker string, dataSource DataSource, recModel *recommendationmodel.ExternalStockData) *MarketData {
	now := time.Now()

	// Convert MarketCap from int64 to *int64
	var marketCap *int64
	if recModel.MarketCap > 0 {
		marketCap = &recModel.MarketCap
	}

	return &MarketData{
		ID:               uuid.New(),
		Ticker:           ticker,
		DataSource:       dataSource,
		DataQuality:      DataQualityGood, // Default quality
		CurrentPrice:     recModel.CurrentPrice,
		DayChange:        recModel.DayChange,
		DayChangePercent: recModel.DayChangePercent,
		Volume:           recModel.Volume,
		MarketCap:        marketCap,
		PERatio:          recModel.PERatio,
		DividendYield:    recModel.DividendYield,
		Week52High:       recModel.Week52High,
		Week52Low:        recModel.Week52Low,
		AvgVolume:        recModel.AvgVolume,
		CollectedAt:      now,
		DataTimestamp:    recModel.LastUpdated,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Note: This model is kept for backward compatibility but should not be used.
// Use the recommendation model's ExternalStockData instead.
type ExternalStockData struct {
	CurrentPrice     float64   `json:"current_price"`
	DayChange        float64   `json:"day_change"`
	DayChangePercent float64   `json:"day_change_percent"`
	Volume           int64     `json:"volume"`
	MarketCap        *int64    `json:"market_cap,omitempty"`
	PERatio          *float64  `json:"pe_ratio,omitempty"`
	DividendYield    *float64  `json:"dividend_yield,omitempty"`
	Week52High       *float64  `json:"week_52_high,omitempty"`
	Week52Low        *float64  `json:"week_52_low,omitempty"`
	AvgVolume        *int64    `json:"avg_volume,omitempty"`
	LastUpdated      time.Time `json:"last_updated"`
}
