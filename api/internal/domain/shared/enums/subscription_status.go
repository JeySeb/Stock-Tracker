package enums

// SubscriptionStatus represents the current status of a subscription
type SubscriptionStatus string

const (
	// STATUS_ACTIVE represents an active subscription
	STATUS_ACTIVE SubscriptionStatus = "active"

	// STATUS_CANCELLED represents a cancelled subscription
	STATUS_CANCELLED SubscriptionStatus = "cancelled"

	// STATUS_EXPIRED represents an expired subscription
	STATUS_EXPIRED SubscriptionStatus = "expired"

	// STATUS_PENDING represents a pending subscription
	STATUS_PENDING SubscriptionStatus = "pending"
)

// IsValid checks if the SubscriptionStatus value is valid
func (s SubscriptionStatus) IsValid() bool {
	switch s {
	case STATUS_ACTIVE, STATUS_CANCELLED, STATUS_EXPIRED, STATUS_PENDING:
		return true
	default:
		return false
	}
}

// String returns the string representation of the SubscriptionStatus
func (s SubscriptionStatus) String() string {
	return string(s)
}
