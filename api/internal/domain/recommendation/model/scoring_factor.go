package model

import "stock-tracker/internal/domain/shared/enums"

type ScoringFactor struct {
	Name        string                   `json:"name"`
	Score       float64                  `json:"score"`
	Weight      float64                  `json:"weight"`
	Explanation string                   `json:"explanation"`
	Tier        enums.RecommendationTier `json:"tier"` // "basic", "enriched", "premium"
}
