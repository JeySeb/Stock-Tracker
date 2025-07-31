package enums

// SubscriptionPrice represents the price of a subscription plan
type SubscriptionPrice float64

const (
	// PRICE_MONTHLY represents the price of a monthly subscription plan
	PRICE_MONTHLY SubscriptionPrice = 29.99

	// PRICE_YEARLY represents the price of a yearly subscription plan
	PRICE_YEARLY SubscriptionPrice = 299.99
)