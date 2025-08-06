package model

import (
	"time"

	"stock-tracker/internal/domain/shared/enums"

	"github.com/google/uuid"
)

// AggregatedRecommendation represents a consolidated recommendation based on multiple stock events
type AggregatedRecommendation struct {
	ID          uuid.UUID `json:"id"`
	Ticker      string    `json:"ticker"`
	CompanyName string    `json:"company_name"`

	// Aggregated Basic Data (calculated from events in DB)
	TotalEvents       int     `json:"total_events"`
	PositiveEvents    int     `json:"positive_events"`
	NegativeEvents    int     `json:"negative_events"`
	AvgTargetChange   float64 `json:"avg_target_change"`
	LatestTargetPrice float64 `json:"latest_target_price"`
	BrokerConsensus   float64 `json:"broker_consensus"` // -1 to 1

	// Calculated Scores
	BasicScore         float64                  `json:"basic_score"` // 0 to 1
	Confidence         float64                  `json:"confidence"`  // 0 to 1
	RecommendationType enums.RecommendationType `json:"recommendation_type"`

	// Scoring Details
	ScoringFactors []ScoringFactor `json:"scoring_factors"`

	// Tiered Data
	Tier         enums.RecommendationTier `json:"tier"`
	ExternalData *ExternalStockData       `json:"external_data,omitempty"`
	AIInsights   *AIGeneratedInsights     `json:"ai_insights,omitempty"`

	// Metadata
	LastEventTime time.Time `json:"last_event_time"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// EventStatistics holds statistical data calculated from multiple events
type EventStatistics struct {
	TotalEvents       int
	PositiveCount     int
	NegativeCount     int
	AvgTargetChange   float64
	LatestTargetPrice float64
	LatestEventTime   time.Time
}

// ScoringWeights defines the weights for different scoring factors
type ScoringWeights struct {
	BrokerFrequency      float64 // Weight for broker credibility/frequency
	TargetMovement       float64 // Weight for price target changes
	RatingChange         float64 // Weight for rating upgrades/downgrades
	Recency              float64 // Weight for how recent the recommendations are
	Consensus            float64 // Weight for directional consensus [-1, +1]
	Certainty            float64 // Weight for non-directional certainty [0, 1]
	DirectionalCertainty float64 // Weight for directional certainty [-1, +1]
	ConsensusStrength    float64 // Weight for consensus strength regardless of direction
}

// NewAggregatedRecommendation creates a new aggregated recommendation instance
func NewAggregatedRecommendation(ticker, companyName string, tier enums.RecommendationTier) *AggregatedRecommendation {
	return &AggregatedRecommendation{
		ID:          uuid.New(),
		Ticker:      ticker,
		CompanyName: companyName,
		Tier:        tier,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(getAggregatedExpirationByTier(tier)),
	}
}

// getAggregatedExpirationByTier returns the appropriate TTL based on recommendation tier
func getAggregatedExpirationByTier(tier enums.RecommendationTier) time.Duration {
	switch tier {
	case enums.RECOMMENDATION_TIER_BASIC:
		return 12 * time.Hour // Basic data changes less frequently
	case enums.RECOMMENDATION_TIER_ENRICHED:
		return 4 * time.Hour // External data updates more frequently
	case enums.RECOMMENDATION_TIER_PREMIUM:
		return 2 * time.Hour // Premium users get fresher data
	default:
		return 12 * time.Hour
	}
}

// DetermineType sets the recommendation type based on the basic score
// Updated thresholds to better distribute recommendations across all types
func (r *AggregatedRecommendation) DetermineType() {
	switch {
	case r.BasicScore >= 0.70:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_STRONG_BUY
	case r.BasicScore >= 0.55:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_BUY
	case r.BasicScore >= 0.40:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_HOLD
	case r.BasicScore >= 0.25:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_SELL
	default:
		r.RecommendationType = enums.RECOMMENDATION_TYPE_STRONG_SELL
	}
}

// IsExpired checks if the recommendation has expired
func (r *AggregatedRecommendation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// CanAccessExternalData checks if the user tier can access external data
func (r *AggregatedRecommendation) CanAccessExternalData(userTier enums.UserTier) bool {
	return userTier == enums.TIER_BASIC || userTier == enums.TIER_PREMIUM
}

// CanAccessAIInsights checks if the user tier can access AI insights
func (r *AggregatedRecommendation) CanAccessAIInsights(userTier enums.UserTier) bool {
	return userTier == enums.TIER_PREMIUM
}

// GetPrimaryFactors returns the most important scoring factors for display
func (r *AggregatedRecommendation) GetPrimaryFactors() []ScoringFactor {
	var primary []ScoringFactor
	for _, factor := range r.ScoringFactors {
		if factor.Weight >= 0.2 { // Only factors with significant weight
			primary = append(primary, factor)
		}
	}
	return primary
}

// GetUpside calculates the upside potential if external data is available
func (r *AggregatedRecommendation) GetUpside() *float64 {
	if r.ExternalData != nil && r.LatestTargetPrice > 0 && r.ExternalData.CurrentPrice > 0 {
		upside := (r.LatestTargetPrice - r.ExternalData.CurrentPrice) / r.ExternalData.CurrentPrice
		return &upside
	}
	return nil
}
