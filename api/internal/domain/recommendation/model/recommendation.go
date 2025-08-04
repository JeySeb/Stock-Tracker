package model

import (
	"time"

	"stock-tracker/internal/domain/shared/enums"

	"github.com/google/uuid"
)

type Recommendation struct {
	ID                 uuid.UUID                `json:"id" db:"id"`
	Ticker             string                   `json:"ticker" db:"ticker"`
	CompanyName        string                   `json:"company_name" db:"company_name"`
	Score              float64                  `json:"score" db:"score"`
	Confidence         float64                  `json:"confidence" db:"confidence"`
	Tier               enums.RecommendationTier `json:"tier" db:"tier"`
	RecommendationType enums.RecommendationType `json:"recommendation_type" db:"recommendation_type"`

	// Basic factors (available to all users)
	BasicFactors []ScoringFactor `json:"basic_factors" db:"basic_factors"`

	// Enriched data (registered users only)
	ExternalData *ExternalStockData `json:"external_data,omitempty" db:"external_data"`

	// Premium insights (premium users only)
	AIInsights *AIGeneratedInsights `json:"ai_insights,omitempty" db:"ai_insights"`

	Explanation string    `json:"explanation" db:"explanation"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}

func NewRecommendation(ticker, companyName string, score, confidence float64, tier enums.RecommendationTier) *Recommendation {
	return &Recommendation{
		ID:          uuid.New(),
		Ticker:      ticker,
		CompanyName: companyName,
		Score:       score,
		Confidence:  confidence,
		Tier:        tier,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(getExpirationByTier(tier)),
	}
}

func getExpirationByTier(tier enums.RecommendationTier) time.Duration {
	switch tier {
	case enums.RECOMMENDATION_TIER_BASIC:
		return 24 * time.Hour // 24 hours for basic
	case enums.RECOMMENDATION_TIER_ENRICHED:
		return 12 * time.Hour // 12 hours for enriched
	case enums.RECOMMENDATION_TIER_PREMIUM:
		return 4 * time.Hour // 4 hours for premium (more frequent updates)
	default:
		return 24 * time.Hour
	}
}

func (r *Recommendation) DetermineType() {
	switch {
	case r.Score >= 0.75:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_STRONG_BUY
	case r.Score >= 0.55:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_BUY
	case r.Score >= 0.35:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_HOLD
	case r.Score >= 0.15:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_SELL
	default:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_STRONG_SELL
	}
}

func (r *Recommendation) CanAccessExternalData(userTier enums.UserTier) bool {
	return userTier == enums.TIER_BASIC || userTier == enums.TIER_PREMIUM
}

func (r *Recommendation) CanAccessAIInsights(userTier enums.UserTier) bool {
	return userTier == enums.TIER_PREMIUM
}
