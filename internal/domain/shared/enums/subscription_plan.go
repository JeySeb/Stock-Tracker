package enums

// SubscriptionPlan represents the type of subscription plan
type SubscriptionPlan string

const (
	// PLAN_MONTHLY represents a monthly subscription plan
	PLAN_MONTHLY SubscriptionPlan = "monthly"

	// PLAN_YEARLY represents a yearly subscription plan 
	PLAN_YEARLY SubscriptionPlan = "yearly"
)

// IsValid checks if the SubscriptionPlan value is valid
func (p SubscriptionPlan) IsValid() bool {
	switch p {
	case PLAN_MONTHLY, PLAN_YEARLY:
		return true
	default:
		return false
	}
}

// String returns the string representation of the SubscriptionPlan
func (p SubscriptionPlan) String() string {
	return string(p)
}
