package model

import (
	"time"

	"github.com/google/uuid"
)

type MarketTrend struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Sector         string    `json:"sector" db:"sector"`
	Industry       string    `json:"industry,omitempty" db:"industry"`
	UpgradeCount   int       `json:"upgrade_count" db:"upgrade_count"`
	DowngradeCount int       `json:"downgrade_count" db:"downgrade_count"`
	TotalActions   int       `json:"total_actions" db:"total_actions"`
	TrendScore     float64   `json:"trend_score" db:"trend_score"`         // -1 to 1
	TrendDirection string    `json:"trend_direction" db:"trend_direction"` // "Bullish", "Bearish", "Neutral"
	Period         string    `json:"period" db:"period"`                   // "7d", "30d", "90d"
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

func NewMarketTrend(sector string, period string) *MarketTrend {
	return &MarketTrend{
		ID:        uuid.New(),
		Sector:    sector,
		Period:    period,
		CreatedAt: time.Now(),
	}
}

func (mt *MarketTrend) CalculateTrendScore() {
	if mt.TotalActions == 0 {
		mt.TrendScore = 0
		mt.TrendDirection = "Neutral"
		return
	}

	// Calculate trend score: (upgrades - downgrades) / total_actions
	mt.TrendScore = float64(mt.UpgradeCount-mt.DowngradeCount) / float64(mt.TotalActions)

	// Determine trend direction
	switch {
	case mt.TrendScore > 0.3:
		mt.TrendDirection = "Bullish"
	case mt.TrendScore < -0.3:
		mt.TrendDirection = "Bearish"
	default:
		mt.TrendDirection = "Neutral"
	}
}
