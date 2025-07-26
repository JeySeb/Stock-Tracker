package model

type ExternalStockData struct {
	CurrentPrice     float64  `json:"current_price"`
	DayChange        float64  `json:"day_change"`
	DayChangePercent float64  `json:"day_change_percent"`
	Volume           int64    `json:"volume"`
	MarketCap        int64    `json:"market_cap"`
	PERatio          *float64 `json:"pe_ratio,omitempty"`
}
