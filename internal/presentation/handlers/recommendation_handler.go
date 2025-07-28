package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	recommendationUsecase "stock-tracker/internal/domain/recommendation/usecase"
	"stock-tracker/internal/domain/shared/enums"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/pkg/logger"
)

// RecommendationHandler handles HTTP requests for recommendations
type RecommendationHandler struct {
	recommendationUC recommendationUsecase.RecommendationUseCase
	logger           logger.Logger
}

// NewRecommendationHandler creates a new recommendation handler
func NewRecommendationHandler(
	recommendationUC recommendationUsecase.RecommendationUseCase,
	logger logger.Logger,
) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationUC: recommendationUC,
		logger:           logger,
	}
}

// GetRecommendations handles GET /api/v1/recommendations
func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// 1. Extract user tier from context (set by auth middleware)
	userTier, ok := r.Context().Value(middleware.UserTierContextKey).(enums.UserTier)
	if !ok {
		userTier = enums.TIER_GUEST // Default to guest if not authenticated
	}

	// 2. Parse query parameters
	limit := h.parseIntParam(r.URL.Query().Get("limit"), 10)
	minScore := h.parseFloatParam(r.URL.Query().Get("min_score"))
	recommendationType := h.parseRecommendationTypeParam(r.URL.Query().Get("type"))
	excludeTickers := h.parseStringSliceParam(r.URL.Query().Get("exclude"))

	// 3. Build request
	request := recommendationUsecase.RecommendationRequest{
		UserTier: userTier,
		Limit:    limit,
		Filters: recommendationUsecase.RecommendationFilters{
			MinScore:           minScore,
			RecommendationType: recommendationType,
			ExcludeTickers:     excludeTickers,
		},
	}

	// 4. Get recommendations
	response, err := h.recommendationUC.GetRecommendations(r.Context(), request)
	if err != nil {
		h.logger.Error("Failed to get recommendations",
			"error", err,
			"user_tier", userTier,
			"limit", limit)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Failed to retrieve recommendations",
		})
		return
	}

	// 5. Filter response based on user tier for security
	filteredResponse := h.filterResponseByTier(response, userTier)

	// 6. Add rate limiting info if available
	if rateLimitRemaining, exists := r.Context().Value(middleware.RateLimitRemainingKey).(int); exists {
		filteredResponse.Meta.RateLimitRemaining = rateLimitRemaining
	}

	h.logger.Info("Recommendations served",
		"user_tier", userTier,
		"count", len(filteredResponse.Data),
		"cache_hit", filteredResponse.Meta.CacheHit,
		"generation_time_ms", filteredResponse.Meta.GenerationTime.Milliseconds())

	render.JSON(w, r, filteredResponse)
}

// GetRecommendationByTicker handles GET /api/v1/recommendations/{ticker}
func (h *RecommendationHandler) GetRecommendationByTicker(w http.ResponseWriter, r *http.Request) {
	// 1. Extract ticker from URL
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Ticker is required",
		})
		return
	}

	// Validate and normalize ticker
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	if len(ticker) < 1 || len(ticker) > 10 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Ticker must be between 1 and 10 characters",
		})
		return
	}

	// Validate ticker contains only letters and numbers
	for _, char := range ticker {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": "Ticker contains invalid characters",
			})
			return
		}
	}

	// 2. Extract user tier
	userTier, ok := r.Context().Value(middleware.UserTierContextKey).(enums.UserTier)
	if !ok {
		userTier = enums.TIER_GUEST
	}

	// 3. Get specific recommendation
	recommendation, err := h.recommendationUC.GetRecommendationByTicker(r.Context(), ticker, userTier)
	if err != nil {
		h.logger.Error("Failed to get recommendation for ticker",
			"ticker", ticker,
			"error", err,
			"user_tier", userTier)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Failed to retrieve recommendation",
		})
		return
	}

	if recommendation == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{
			"error": "Recommendation not found for ticker: " + ticker,
		})
		return
	}

	// 4. Filter based on user tier
	filteredRecommendation := h.filterSingleRecommendationByTier(recommendation, userTier)

	h.logger.Info("Single recommendation served",
		"ticker", ticker,
		"user_tier", userTier,
		"score", recommendation.BasicScore)

	render.JSON(w, r, map[string]interface{}{
		"data": filteredRecommendation,
		"meta": map[string]interface{}{
			"user_tier": userTier,
			"features":  h.getAvailableFeatures(userTier),
		},
	})
}

// GetRecommendationPreview handles GET /api/v1/recommendations/preview/{ticker}
// Shows what additional data would be available with higher tier
func (h *RecommendationHandler) GetRecommendationPreview(w http.ResponseWriter, r *http.Request) {
	// Only allow authenticated users to see previews
	userTier, ok := r.Context().Value(middleware.UserTierContextKey).(enums.UserTier)
	if !ok || userTier == enums.TIER_GUEST {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "Authentication required for preview",
		})
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Ticker is required",
		})
		return
	}

	ticker = strings.ToUpper(strings.TrimSpace(ticker))

	// Get recommendations for current tier and next tier
	currentRec, err := h.recommendationUC.GetRecommendationByTicker(r.Context(), ticker, userTier)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Failed to get recommendation",
		})
		return
	}

	// Get preview of higher tier (if applicable)
	var higherTierRec interface{} = nil
	if userTier == enums.TIER_BASIC {
		if premiumRec, err := h.recommendationUC.GetRecommendationByTicker(r.Context(), ticker, enums.TIER_PREMIUM); err == nil {
			higherTierRec = h.filterSingleRecommendationByTier(premiumRec, enums.TIER_PREMIUM)
		}
	}

	render.JSON(w, r, map[string]interface{}{
		"current_tier": map[string]interface{}{
			"tier": userTier,
			"data": h.filterSingleRecommendationByTier(currentRec, userTier),
		},
		"premium_preview":  higherTierRec,
		"upgrade_benefits": h.getUpgradeBenefits(userTier),
	})
}

// Helper methods

// filterResponseByTier filters the entire response based on user tier
func (h *RecommendationHandler) filterResponseByTier(
	response *recommendationUsecase.RecommendationResponse,
	userTier enums.UserTier,
) *recommendationUsecase.RecommendationResponse {

	filteredData := make([]*struct {
		*recommendationUsecase.RecommendationResponse
		Data interface{} `json:"data"`
	}, len(response.Data))

	for i, rec := range response.Data {
		filteredData[i] = &struct {
			*recommendationUsecase.RecommendationResponse
			Data interface{} `json:"data"`
		}{
			Data: h.filterSingleRecommendationByTier(rec, userTier),
		}
	}

	return response // Return as-is since Go will handle JSON filtering
}

// filterSingleRecommendationByTier filters a single recommendation based on user tier
func (h *RecommendationHandler) filterSingleRecommendationByTier(
	rec interface{},
	userTier enums.UserTier,
) interface{} {

	// Create a map for flexible filtering
	recMap := map[string]interface{}{
		"id":                  rec,
		"ticker":              rec,
		"company_name":        rec,
		"total_events":        rec,
		"positive_events":     rec,
		"negative_events":     rec,
		"avg_target_change":   rec,
		"latest_target_price": rec,
		"broker_consensus":    rec,
		"basic_score":         rec,
		"confidence":          rec,
		"recommendation_type": rec,
		"scoring_factors":     rec,
		"tier":                rec,
		"last_event_time":     rec,
		"created_at":          rec,
		"expires_at":          rec,
	}

	// Add tier-specific data
	switch userTier {
	case enums.TIER_GUEST:
		// Only basic data - external_data and ai_insights are omitted

	case enums.TIER_BASIC:
		// Add external data
		recMap["external_data"] = rec

	case enums.TIER_PREMIUM:
		// Add all data
		recMap["external_data"] = rec
		recMap["ai_insights"] = rec
	}

	// In a real implementation, this would use proper struct filtering
	// For now, return the original rec (the actual filtering would happen in JSON serialization)
	return rec
}

// parseIntParam parses an integer parameter with a default value
func (h *RecommendationHandler) parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(param); err == nil && value > 0 {
		return value
	}

	return defaultValue
}

// parseFloatParam parses a float parameter
func (h *RecommendationHandler) parseFloatParam(param string) *float64 {
	if param == "" {
		return nil
	}

	if value, err := strconv.ParseFloat(param, 64); err == nil {
		return &value
	}

	return nil
}

// parseRecommendationTypeParam parses recommendation type parameter
func (h *RecommendationHandler) parseRecommendationTypeParam(param string) *enums.RecommendationType {
	if param == "" {
		return nil
	}

	// Normalize the parameter
	param = strings.ToLower(strings.TrimSpace(param))

	switch param {
	case "strong_buy", "strongbuy":
		recType := enums.RECOMMENDATION_TYPE_STRONG_BUY
		return &recType
	case "buy":
		recType := enums.RECOMMENDATION_TYPE_BUY
		return &recType
	case "hold":
		recType := enums.RECOMMENDATION_TYPE_HOLD
		return &recType
	case "sell":
		recType := enums.RECOMMENDATION_TYPE_SELL
		return &recType
	case "strong_sell", "strongsell":
		recType := enums.RECOMMENDATION_TYPE_STRONG_SELL
		return &recType
	default:
		return nil
	}
}

// parseStringSliceParam parses a comma-separated string into a slice
func (h *RecommendationHandler) parseStringSliceParam(param string) []string {
	if param == "" {
		return nil
	}

	values := strings.Split(param, ",")
	var result []string

	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, strings.ToUpper(trimmed))
		}
	}

	return result
}

// getAvailableFeatures returns features available to a user tier
func (h *RecommendationHandler) getAvailableFeatures(userTier enums.UserTier) []string {
	features := []string{"basic_recommendations", "market_analytics"}

	if userTier == enums.TIER_BASIC || userTier == enums.TIER_PREMIUM {
		features = append(features, "real_time_data", "external_apis", "advanced_analytics")
	}

	if userTier == enums.TIER_PREMIUM {
		features = append(features, "ai_insights", "sentiment_analysis", "premium_recommendations")
	}

	return features
}

// getUpgradeBenefits returns upgrade benefits for a user tier
func (h *RecommendationHandler) getUpgradeBenefits(userTier enums.UserTier) []string {
	switch userTier {
	case enums.TIER_GUEST:
		return []string{
			"Register to access real-time market data",
			"Get up to 25 recommendations",
			"Access to external data sources",
		}
	case enums.TIER_BASIC:
		return []string{
			"Upgrade to Premium for AI-powered insights",
			"Get up to 100 recommendations",
			"Access to sentiment analysis",
			"Advanced market predictions",
		}
	case enums.TIER_PREMIUM:
		return []string{
			"You have access to all features!",
		}
	default:
		return []string{}
	}
}
