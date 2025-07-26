package enums

// RecommendationType represents the type of stock recommendation
type RecommendationType string

const (
	// RECOMMENDATION_TYPE_STRONG_BUY indicates a strong recommendation to buy the stock
	RECOMMENDATION_TYPE_STRONG_BUY RecommendationType = "Strong Buy"

	// RECOMMENDATION_TYPE_BUY indicates a recommendation to buy the stock
	RECOMMENDATION_TYPE_BUY RecommendationType = "Buy"

	// RECOMMENDATION_TYPE_HOLD indicates a recommendation to hold the stock
	RECOMMENDATION_TYPE_HOLD RecommendationType = "Hold"

	// RECOMMENDATION_TYPE_SELL indicates a recommendation to sell the stock
	RECOMMENDATION_TYPE_SELL RecommendationType = "Sell"

	// RECOMMENDATION_TYPE_STRONG_SELL indicates a strong recommendation to sell the stock
	RECOMMENDATION_TYPE_STRONG_SELL RecommendationType = "Strong Sell"
)
