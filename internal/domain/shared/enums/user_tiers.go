package enums

// UserTier represents the access tier level of a user
type UserTier string

const (
	// TIER_GUEST represents non-registered users with limited access
	TIER_GUEST UserTier = "guest"

	// TIER_BASIC represents registered users with standard access
	TIER_BASIC UserTier = "basic"

	// TIER_PREMIUM represents premium subscribers with full access
	TIER_PREMIUM UserTier = "premium"
)

// API Rate Limits per hour for each tier
const (
	GuestRateLimit   = 100
	BasicRateLimit   = 500
	PremiumRateLimit = 2000
	DefaultRateLimit = 50
)

// GetRateLimit returns the API rate limit for the given user tier
func (t UserTier) GetRateLimit() int {
	switch t {
	case TIER_GUEST:
		return GuestRateLimit
	case TIER_BASIC:
		return BasicRateLimit
	case TIER_PREMIUM:
		return PremiumRateLimit
	default:
		return DefaultRateLimit
	}
}

// IsValid checks if the UserTier value is valid
func (t UserTier) IsValid() bool {
	switch t {
	case TIER_GUEST, TIER_BASIC, TIER_PREMIUM:
		return true
	default:
		return false
	}
}

// String returns the string representation of the UserTier
func (t UserTier) String() string {
	return string(t)
}
