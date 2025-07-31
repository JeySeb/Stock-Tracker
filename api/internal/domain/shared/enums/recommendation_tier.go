package enums

// RecommendationTier represents the tier level of a recommendation
type RecommendationTier string

const (
	// RecommendationTierBasic represents basic tier recommendations with limited data
	RECOMMENDATION_TIER_BASIC RecommendationTier = "basic"

	// RecommendationTierEnriched represents enriched tier recommendations with additional data
	RECOMMENDATION_TIER_ENRICHED RecommendationTier = "enriched"

	// RecommendationTierPremium represents premium tier recommendations with full data access
	RECOMMENDATION_TIER_PREMIUM RecommendationTier = "premium"
)

// IsValid checks if the RecommendationTier value is valid
func (t RecommendationTier) IsValid() bool {
	switch t {
	case RECOMMENDATION_TIER_BASIC, RECOMMENDATION_TIER_ENRICHED, RECOMMENDATION_TIER_PREMIUM:
		return true
	default:
		return false
	}
}

// String returns the string representation of the RecommendationTier
func (t RecommendationTier) String() string {
	return string(t)
}
