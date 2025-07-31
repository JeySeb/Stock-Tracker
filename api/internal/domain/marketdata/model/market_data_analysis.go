package model

import (
	"time"

	"github.com/google/uuid"
)

// MarketDataAnalysis represents analytical data for market analysis
type MarketDataAnalysis struct {
	ID               uuid.UUID   `json:"id" db:"id"`
	Ticker           string      `json:"ticker" db:"ticker"`
	CurrentPrice     float64     `json:"current_price" db:"current_price"`
	DayChange        float64     `json:"day_change" db:"day_change"`
	DayChangePercent float64     `json:"day_change_percent" db:"day_change_percent"`
	Volume           int64       `json:"volume" db:"volume"`
	Week52High       *float64    `json:"week_52_high,omitempty" db:"week_52_high"`
	Week52Low        *float64    `json:"week_52_low,omitempty" db:"week_52_low"`
	DataTimestamp    time.Time   `json:"data_timestamp" db:"data_timestamp"`
	CollectedAt      time.Time   `json:"collected_at" db:"collected_at"`
	DataQuality      DataQuality `json:"data_quality" db:"data_quality"`
	DataSource       DataSource  `json:"data_source" db:"data_source"`

	// Calculated fields
	PricePosition   float64 `json:"price_position" db:"price_position"`     // Position within 52-week range
	VolumeActivity  float64 `json:"volume_activity" db:"volume_activity"`   // Volume relative to average
	VolatilityScore float64 `json:"volatility_score" db:"volatility_score"` // Day change volatility
	TrendDirection  string  `json:"trend_direction" db:"trend_direction"`   // "bullish", "bearish", "neutral"
	RiskLevel       string  `json:"risk_level" db:"risk_level"`             // "low", "medium", "high"
}

// MarketDataTrend represents trend analysis for a ticker
type MarketDataTrend struct {
	Ticker             string  `json:"ticker" db:"ticker"`
	Period             string  `json:"period" db:"period"` // "1d", "1w", "1m", "3m"
	StartPrice         float64 `json:"start_price" db:"start_price"`
	EndPrice           float64 `json:"end_price" db:"end_price"`
	TotalChange        float64 `json:"total_change" db:"total_change"`
	TotalChangePercent float64 `json:"total_change_percent" db:"total_change_percent"`
	MaxPrice           float64 `json:"max_price" db:"max_price"`
	MinPrice           float64 `json:"min_price" db:"min_price"`
	AvgVolume          int64   `json:"avg_volume" db:"avg_volume"`
	Volatility         float64 `json:"volatility" db:"volatility"`
	TrendStrength      string  `json:"trend_strength" db:"trend_strength"` // "strong", "moderate", "weak"
	Direction          string  `json:"direction" db:"direction"`           // "up", "down", "sideways"
	DataPoints         int     `json:"data_points" db:"data_points"`
}

// MarketDataComparison represents comparison between tickers
type MarketDataComparison struct {
	Ticker1          string    `json:"ticker1" db:"ticker1"`
	Ticker2          string    `json:"ticker2" db:"ticker2"`
	ComparisonDate   time.Time `json:"comparison_date" db:"comparison_date"`
	Price1           float64   `json:"price1" db:"price1"`
	Price2           float64   `json:"price2" db:"price2"`
	Change1          float64   `json:"change1" db:"change1"`
	Change2          float64   `json:"change2" db:"change2"`
	ChangePercent1   float64   `json:"change_percent1" db:"change_percent1"`
	ChangePercent2   float64   `json:"change_percent2" db:"change_percent2"`
	Volume1          int64     `json:"volume1" db:"volume1"`
	Volume2          int64     `json:"volume2" db:"volume2"`
	RelativeStrength float64   `json:"relative_strength" db:"relative_strength"` // Ticker1 vs Ticker2
	Correlation      float64   `json:"correlation" db:"correlation"`             // Price correlation
}

// MarketDataSummary represents summary statistics for market data
type MarketDataSummary struct {
	TotalRecords        int64     `json:"total_records" db:"total_records"`
	UniqueTickers       int64     `json:"unique_tickers" db:"unique_tickers"`
	AvgPrice            float64   `json:"avg_price" db:"avg_price"`
	AvgDayChange        float64   `json:"avg_day_change" db:"avg_day_change"`
	AvgDayChangePercent float64   `json:"avg_day_change_percent" db:"avg_day_change_percent"`
	TotalVolume         int64     `json:"total_volume" db:"total_volume"`
	AvgVolume           int64     `json:"avg_volume" db:"avg_volume"`
	MostActiveTicker    string    `json:"most_active_ticker" db:"most_active_ticker"`
	BestPerformer       string    `json:"best_performer" db:"best_performer"`
	WorstPerformer      string    `json:"worst_performer" db:"worst_performer"`
	BullishCount        int64     `json:"bullish_count" db:"bullish_count"`
	BearishCount        int64     `json:"bearish_count" db:"bearish_count"`
	NeutralCount        int64     `json:"neutral_count" db:"neutral_count"`
	Period              string    `json:"period" db:"period"`
	LastUpdated         time.Time `json:"last_updated" db:"last_updated"`
}

// MarketDataAlert represents alerts for significant market movements
type MarketDataAlert struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	Ticker             string    `json:"ticker" db:"ticker"`
	AlertType          string    `json:"alert_type" db:"alert_type"` // "price_spike", "volume_surge", "breakout", "breakdown"
	Severity           string    `json:"severity" db:"severity"`     // "low", "medium", "high", "critical"
	Message            string    `json:"message" db:"message"`
	CurrentPrice       float64   `json:"current_price" db:"current_price"`
	PreviousPrice      float64   `json:"previous_price" db:"previous_price"`
	PriceChange        float64   `json:"price_change" db:"price_change"`
	PriceChangePercent float64   `json:"price_change_percent" db:"price_change_percent"`
	Volume             int64     `json:"volume" db:"volume"`
	Threshold          float64   `json:"threshold" db:"threshold"`
	TriggeredAt        time.Time `json:"triggered_at" db:"triggered_at"`
	IsActive           bool      `json:"is_active" db:"is_active"`
}

// MarketDataFilters represents filters for market data queries
type MarketDataFilters struct {
	Ticker         string    `json:"ticker,omitempty"`
	DataSource     string    `json:"data_source,omitempty"`
	DataQuality    string    `json:"data_quality,omitempty"`
	StartDate      time.Time `json:"start_date,omitempty"`
	EndDate        time.Time `json:"end_date,omitempty"`
	MinPrice       *float64  `json:"min_price,omitempty"`
	MaxPrice       *float64  `json:"max_price,omitempty"`
	MinChange      *float64  `json:"min_change,omitempty"`
	MaxChange      *float64  `json:"max_change,omitempty"`
	MinVolume      *int64    `json:"min_volume,omitempty"`
	MaxVolume      *int64    `json:"max_volume,omitempty"`
	TrendDirection string    `json:"trend_direction,omitempty"`
	RiskLevel      string    `json:"risk_level,omitempty"`
	SortBy         string    `json:"sort_by,omitempty"`
	SortOrder      string    `json:"sort_order,omitempty"`
	Limit          int       `json:"limit,omitempty"`
	Offset         int       `json:"offset,omitempty"`
}

// SetDefaults sets default values for filters
func (f *MarketDataFilters) SetDefaults() {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.SortBy == "" {
		f.SortBy = "data_timestamp"
	}
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
}

// MarketDataResponse represents API response structure
type MarketDataResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Message    string      `json:"message,omitempty"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

// Metadata represents additional response metadata
type Metadata struct {
	GeneratedAt time.Time              `json:"generated_at"`
	DataPoints  int                    `json:"data_points"`
	Period      string                 `json:"period,omitempty"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
}

// NewMarketDataAnalysis creates a new market data analysis instance
func NewMarketDataAnalysis(ticker string, currentPrice, dayChange, dayChangePercent float64, volume int64, week52High, week52Low *float64, dataTimestamp, collectedAt time.Time, dataQuality DataQuality, dataSource DataSource) *MarketDataAnalysis {
	analysis := &MarketDataAnalysis{
		ID:               uuid.New(),
		Ticker:           ticker,
		CurrentPrice:     currentPrice,
		DayChange:        dayChange,
		DayChangePercent: dayChangePercent,
		Volume:           volume,
		Week52High:       week52High,
		Week52Low:        week52Low,
		DataTimestamp:    dataTimestamp,
		CollectedAt:      collectedAt,
		DataQuality:      dataQuality,
		DataSource:       dataSource,
	}

	// Calculate derived fields
	analysis.CalculateDerivedFields()
	return analysis
}

// CalculateDerivedFields calculates all derived analytical fields
func (a *MarketDataAnalysis) CalculateDerivedFields() {
	// Calculate price position within 52-week range
	if a.Week52High != nil && a.Week52Low != nil && *a.Week52High > *a.Week52Low {
		priceRange := *a.Week52High - *a.Week52Low
		if priceRange > 0 {
			a.PricePosition = (a.CurrentPrice - *a.Week52Low) / priceRange
		}
	}

	// Calculate volatility score (absolute day change percent)
	a.VolatilityScore = abs(a.DayChangePercent)

	// Determine trend direction
	a.TrendDirection = a.determineTrendDirection()

	// Determine risk level
	a.RiskLevel = a.determineRiskLevel()
}

// determineTrendDirection determines if the trend is bullish, bearish, or neutral
func (a *MarketDataAnalysis) determineTrendDirection() string {
	if a.DayChangePercent > 2.0 {
		return "bullish"
	} else if a.DayChangePercent < -2.0 {
		return "bearish"
	}
	return "neutral"
}

// determineRiskLevel determines the risk level based on volatility and price position
func (a *MarketDataAnalysis) determineRiskLevel() string {
	volatilityScore := a.VolatilityScore
	pricePosition := a.PricePosition

	// High risk: high volatility or extreme price positions
	if volatilityScore > 5.0 || pricePosition > 0.9 || pricePosition < 0.1 {
		return "high"
	}
	// Medium risk: moderate volatility
	if volatilityScore > 2.0 {
		return "medium"
	}
	// Low risk: low volatility and moderate price position
	return "low"
}

// abs returns absolute value of float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// IsBullish returns true if the stock is showing bullish signals
func (a *MarketDataAnalysis) IsBullish() bool {
	return a.TrendDirection == "bullish"
}

// IsBearish returns true if the stock is showing bearish signals
func (a *MarketDataAnalysis) IsBearish() bool {
	return a.TrendDirection == "bearish"
}

// IsHighRisk returns true if the stock is considered high risk
func (a *MarketDataAnalysis) IsHighRisk() bool {
	return a.RiskLevel == "high"
}

// GetPriceTarget calculates potential price targets based on current position
func (a *MarketDataAnalysis) GetPriceTarget() (support, resistance float64) {
	if a.Week52Low != nil {
		support = *a.Week52Low
	}
	if a.Week52High != nil {
		resistance = *a.Week52High
	}
	return
}
