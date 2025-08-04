package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"stock-tracker/internal/domain/shared/enums"
	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	stockUsecase "stock-tracker/internal/domain/stocks/usecase"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/render"
)

// BrokerHandler handles HTTP requests for broker-related operations
type BrokerHandler struct {
	brokerUC stockUsecase.BrokerUseCase
	logger   logger.Logger
}

// NewBrokerHandler creates a new broker handler
func NewBrokerHandler(
	brokerUC stockUsecase.BrokerUseCase,
	logger logger.Logger,
) *BrokerHandler {
	return &BrokerHandler{
		brokerUC: brokerUC,
		logger:   logger,
	}
}

// GetBrokersWithScores handles GET /api/v1/brokers/scores
func (h *BrokerHandler) GetBrokersWithScores(w http.ResponseWriter, r *http.Request) {
	// 1. Extract user tier from context (set by auth middleware)
	userTier, ok := r.Context().Value(middleware.UserTierContextKey).(enums.UserTier)
	if !ok {
		userTier = enums.TIER_GUEST // Default to guest if not authenticated
	}

	// 2. Parse query parameters
	limit := h.parseIntParam(r.URL.Query().Get("limit"))
	orderBy := h.parseOrderByParam(r.URL.Query().Get("order"))

	// 3. Get brokers with scores
	brokers, err := h.brokerUC.GetBrokersWithScores(r.Context(), limit, orderBy)
	if err != nil {
		h.logger.Error("Failed to get brokers with scores",
			"error", err,
			"user_tier", userTier,
			"limit", limit,
			"order", orderBy)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Failed to retrieve brokers with scores",
		})
		return
	}

	// 4. Filter response based on user tier for security
	filteredBrokers := h.filterBrokersByTier(brokers, userTier)

	// 5. Add rate limiting info if available
	meta := map[string]interface{}{
		"count":     len(filteredBrokers),
		"user_tier": userTier,
		"features":  h.getAvailableFeatures(userTier),
	}

	if rateLimitRemaining, exists := r.Context().Value(middleware.RateLimitRemainingKey).(int); exists {
		meta["rate_limit_remaining"] = rateLimitRemaining
	}

	h.logger.Info("Brokers with scores served",
		"user_tier", userTier,
		"count", len(filteredBrokers),
		"limit", limit,
		"order", orderBy)

	render.JSON(w, r, map[string]interface{}{
		"data": filteredBrokers,
		"meta": meta,
	})
}

// Helper methods

// filterBrokersByTier filters the brokers based on user tier for security
func (h *BrokerHandler) filterBrokersByTier(
	brokers []*stockRepos.BrokerWithScore,
	userTier enums.UserTier,
) []*stockRepos.BrokerWithScore {
	// For now, all tiers can see the same broker data
	// In the future, this could be extended to show different levels of detail
	return brokers
}

// parseIntParam parses an integer parameter
func (h *BrokerHandler) parseIntParam(param string) *int {
	if param == "" {
		return nil
	}

	if value, err := strconv.Atoi(param); err == nil && value > 0 {
		return &value
	}

	return nil
}

// parseOrderByParam parses the orderBy parameter
func (h *BrokerHandler) parseOrderByParam(param string) string {
	if param == "" {
		return "desc" // Default to descending order
	}

	param = strings.ToLower(strings.TrimSpace(param))

	switch param {
	case "asc", "ascending":
		return "asc"
	case "desc", "descending":
		return "desc"
	default:
		return "desc" // Default to descending order
	}
}

// getAvailableFeatures returns features available to a user tier
func (h *BrokerHandler) getAvailableFeatures(userTier enums.UserTier) []string {
	features := []string{"broker_scores", "basic_analytics"}

	if userTier == enums.TIER_BASIC || userTier == enums.TIER_PREMIUM {
		features = append(features, "detailed_scores", "credibility_metrics")
	}

	if userTier == enums.TIER_PREMIUM {
		features = append(features, "advanced_analytics", "trend_analysis")
	}

	return features
}
