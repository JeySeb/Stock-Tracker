package validation

import (
	"errors"
	"time"
)

type StockFilters struct {
	Ticker     string     `json:"ticker,omitempty" form:"ticker"`
	Company    string     `json:"company,omitempty" form:"company"`
	Brokerage  string     `json:"brokerage,omitempty" form:"brokerage"`
	Action     string     `json:"action,omitempty" form:"action"`
	RatingFrom string     `json:"rating_from,omitempty" form:"rating_from"`
	RatingTo   string     `json:"rating_to,omitempty" form:"rating_to"`
	DateFrom   *time.Time `json:"date_from,omitempty" form:"date_from"`
	DateTo     *time.Time `json:"date_to,omitempty" form:"date_to"`
	SortBy     string     `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder  string     `json:"sort_order,omitempty" form:"sort_order"`
	Limit      int        `json:"limit,omitempty" form:"limit"`
	Offset     int        `json:"offset,omitempty" form:"offset"`
}

// EnhancedStockFilters provides more robust filtering capabilities
type EnhancedStockFilters struct {
	// Text-based filters (support multiple values)
	Tickers    []string `json:"tickers,omitempty" form:"tickers"`
	Companies  []string `json:"companies,omitempty" form:"companies"`
	Brokerages []string `json:"brokerages,omitempty" form:"brokerages"`
	Actions    []string `json:"actions,omitempty" form:"actions"`

	// Rating filters
	RatingFrom string `json:"rating_from,omitempty" form:"rating_from"`
	RatingTo   string `json:"rating_to,omitempty" form:"rating_to"`

	// Enhanced date filters
	DateFrom   *time.Time  `json:"date_from,omitempty" form:"date_from"`
	DateTo     *time.Time  `json:"date_to,omitempty" form:"date_to"`
	DateRanges []DateRange `json:"date_ranges,omitempty" form:"date_ranges"`

	// Time-based filters
	LastHours  *int `json:"last_hours,omitempty" form:"last_hours"`
	LastDays   *int `json:"last_days,omitempty" form:"last_days"`
	LastWeeks  *int `json:"last_weeks,omitempty" form:"last_weeks"`
	LastMonths *int `json:"last_months,omitempty" form:"last_months"`

	// Target price filters
	TargetFrom *float64 `json:"target_from,omitempty" form:"target_from"`
	TargetTo   *float64 `json:"target_to,omitempty" form:"target_to"`

	// Advanced filters
	MinTargetChange *float64 `json:"min_target_change,omitempty" form:"min_target_change"` // Minimum target price change percentage
	MaxTargetChange *float64 `json:"max_target_change,omitempty" form:"max_target_change"` // Maximum target price change percentage
	HasTargetPrice  *bool    `json:"has_target_price,omitempty" form:"has_target_price"`   // Filter stocks with/without target prices
	HasRating       *bool    `json:"has_rating,omitempty" form:"has_rating"`               // Filter stocks with/without ratings

	// Brokerage credibility filters
	MinBrokerScore *float64 `json:"min_broker_score,omitempty" form:"min_broker_score"` // Minimum brokerage credibility score
	MaxBrokerScore *float64 `json:"max_broker_score,omitempty" form:"max_broker_score"` // Maximum brokerage credibility score

	// Pagination and sorting
	SortBy    string `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder string `json:"sort_order,omitempty" form:"sort_order"`
	Limit     int    `json:"limit,omitempty" form:"limit"`
	Offset    int    `json:"offset,omitempty" form:"offset"`
}

// DateRange represents a specific date range for filtering
type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (f *StockFilters) SetDefaults() {
	if f.Limit <= 0 {
		f.Limit = 50 // TODO: check this value
	}
	if f.Limit > 1000 {
		f.Limit = 1000 // TODO: check this value
	}
	if f.SortBy == "" {
		f.SortBy = "event_time"
	}
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
	// Don't set default date filters - let the query return all data
	// when no date range is specified by the user
}

func (f *EnhancedStockFilters) SetDefaults() {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 1000 {
		f.Limit = 1000
	}
	if f.SortBy == "" {
		f.SortBy = "event_time"
	}
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
}

func (f *StockFilters) Validate() error {
	if f.DateFrom != nil && f.DateTo != nil && f.DateFrom.After(*f.DateTo) {
		return errors.New("date_from must be before date_to")
	}

	if f.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	if f.Limit > 1000 {
		return errors.New("limit must be less than 1000")
	}

	if f.Offset < 0 {
		return errors.New("offset must be greater than 0")
	}

	return nil
}

func (f *EnhancedStockFilters) Validate() error {
	// Validate date ranges
	if f.DateFrom != nil && f.DateTo != nil && f.DateFrom.After(*f.DateTo) {
		return errors.New("date_from must be before date_to")
	}

	// Validate date ranges array
	for _, dateRange := range f.DateRanges {
		if dateRange.From.After(dateRange.To) {
			return errors.New("date range from must be before to")
		}
	}

	// Validate time-based filters
	if f.LastHours != nil && *f.LastHours <= 0 {
		return errors.New("last_hours must be greater than 0")
	}
	if f.LastDays != nil && *f.LastDays <= 0 {
		return errors.New("last_days must be greater than 0")
	}
	if f.LastWeeks != nil && *f.LastWeeks <= 0 {
		return errors.New("last_weeks must be greater than 0")
	}
	if f.LastMonths != nil && *f.LastMonths <= 0 {
		return errors.New("last_months must be greater than 0")
	}

	// Validate target price filters
	if f.TargetFrom != nil && f.TargetTo != nil && *f.TargetFrom > *f.TargetTo {
		return errors.New("target_from must be less than or equal to target_to")
	}

	// Validate target change filters
	if f.MinTargetChange != nil && f.MaxTargetChange != nil && *f.MinTargetChange > *f.MaxTargetChange {
		return errors.New("min_target_change must be less than or equal to max_target_change")
	}

	// Validate brokerage score filters
	if f.MinBrokerScore != nil && f.MaxBrokerScore != nil && *f.MinBrokerScore > *f.MaxBrokerScore {
		return errors.New("min_broker_score must be less than or equal to max_broker_score")
	}

	// Validate pagination
	if f.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	if f.Limit > 1000 {
		return errors.New("limit must be less than 1000")
	}
	if f.Offset < 0 {
		return errors.New("offset must be greater than or equal to 0")
	}

	return nil
}

type Pagination struct {
	Page       int  `json:"page" form:"page"`
	Limit      int  `json:"limit" form:"limit"`
	TotalPages int  `json:"total_pages"`
	TotalItems int  `json:"total_items"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
