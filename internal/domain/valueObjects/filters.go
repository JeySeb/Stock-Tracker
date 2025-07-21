package valueObjects

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

type Pagination struct {
	Page       int  `json:"page" form:"page"`
	Limit      int  `json:"limit" form:"limit"`
	TotalPages int  `json:"total_pages"`
	TotalItems int  `json:"total_items"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
