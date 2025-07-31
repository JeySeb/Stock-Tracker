package model

import "time"

type ExternalStockData struct {
	CurrentPrice     float64   `json:"current_price"`
	DayChange        float64   `json:"day_change"`
	DayChangePercent float64   `json:"day_change_percent"`
	Volume           int64     `json:"volume"`
	MarketCap        int64     `json:"market_cap"`
	PERatio          *float64  `json:"pe_ratio,omitempty"`
	DividendYield    *float64  `json:"dividend_yield,omitempty"`
	Week52High       *float64  `json:"week_52_high,omitempty"`
	Week52Low        *float64  `json:"week_52_low,omitempty"`
	AvgVolume        *int64    `json:"avg_volume,omitempty"`
	LastUpdated      time.Time `json:"last_updated"`
}
