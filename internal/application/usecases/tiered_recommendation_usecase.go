package usecases

import (
	"context"
	"fmt"
	"sort"
	"time"

	"stock-tracker/internal/application/recommendation"
	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	"stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/pkg/logger"
)

// TieredRecommendationUseCase handles business logic for tiered recommendations
type TieredRecommendationUseCase struct {
	stockRepo            repositories.StockRepository
	basicCalculator      *recommendation.BasicScoringCalculator
	externalDataEnricher *recommendation.ExternalDataEnricher
	cache                cache.Cache
	logger               logger.Logger
}

// RecommendationRequest represents a request for recommendations
type RecommendationRequest struct {
	UserTier enums.UserTier
	Limit    int
	Filters  RecommendationFilters
}

// RecommendationFilters allows filtering recommendations
type RecommendationFilters struct {
	MinScore           *float64
	RecommendationType *enums.RecommendationType
	Sectors            []string
	ExcludeTickers     []string
}

// RecommendationResponse represents the response with recommendations and metadata
type RecommendationResponse struct {
	Data []*model.AggregatedRecommendation `json:"data"`
	Meta RecommendationMeta                `json:"meta"`
}

// RecommendationMeta provides metadata about the recommendation response
type RecommendationMeta struct {
	Count              int            `json:"count"`
	UserTier           enums.UserTier `json:"user_tier"`
	Features           []string       `json:"features"`
	CacheHit           bool           `json:"cache_hit"`
	GenerationTime     time.Duration  `json:"generation_time"`
	RateLimitRemaining int            `json:"rate_limit_remaining,omitempty"`
}

// NewTieredRecommendationUseCase creates a new instance of the use case
func NewTieredRecommendationUseCase(
	stockRepo repositories.StockRepository,
	basicCalculator *recommendation.BasicScoringCalculator,
	externalDataEnricher *recommendation.ExternalDataEnricher,
	cache cache.Cache,
	logger logger.Logger,
) *TieredRecommendationUseCase {
	return &TieredRecommendationUseCase{
		stockRepo:            stockRepo,
		basicCalculator:      basicCalculator,
		externalDataEnricher: externalDataEnricher,
		cache:                cache,
		logger:               logger,
	}
}

// GetRecommendations retrieves recommendations based on user tier and request parameters
func (uc *TieredRecommendationUseCase) GetRecommendations(
	ctx context.Context,
	request RecommendationRequest,
) (*RecommendationResponse, error) {

	startTime := time.Now()

	// 1. Apply tier-based limits
	limit := uc.applyTierLimits(request.UserTier, request.Limit)

	// 2. Try to get from cache
	cacheKey := uc.buildCacheKey(request.UserTier, limit, request.Filters)
	cachedRecommendations, cacheHit := uc.tryGetFromCache(ctx, cacheKey)

	if cacheHit && len(cachedRecommendations) > 0 {
		uc.logger.Info("Serving recommendations from cache",
			"tier", request.UserTier,
			"count", len(cachedRecommendations),
			"cache_key", cacheKey)

		return &RecommendationResponse{
			Data: cachedRecommendations,
			Meta: RecommendationMeta{
				Count:          len(cachedRecommendations),
				UserTier:       request.UserTier,
				Features:       uc.getAvailableFeatures(request.UserTier),
				CacheHit:       true,
				GenerationTime: time.Since(startTime),
			},
		}, nil
	}

	// 3. Generate fresh recommendations
	recommendations, err := uc.generateFreshRecommendations(ctx, request.UserTier, limit, request.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	// 4. Cache the results
	cacheTTL := uc.getCacheTTLForTier(request.UserTier)
	if err := uc.cache.Set(ctx, cacheKey, recommendations, cacheTTL); err != nil {
		uc.logger.Warn("Failed to cache recommendations", "error", err)
	}

	generationTime := time.Since(startTime)
	uc.logger.Info("Generated fresh recommendations",
		"tier", request.UserTier,
		"count", len(recommendations),
		"generation_time", generationTime)

	return &RecommendationResponse{
		Data: recommendations,
		Meta: RecommendationMeta{
			Count:          len(recommendations),
			UserTier:       request.UserTier,
			Features:       uc.getAvailableFeatures(request.UserTier),
			CacheHit:       false,
			GenerationTime: generationTime,
		},
	}, nil
}

// GetRecommendationByTicker retrieves a specific recommendation for a ticker
func (uc *TieredRecommendationUseCase) GetRecommendationByTicker(
	ctx context.Context,
	ticker string,
	userTier enums.UserTier,
) (*model.AggregatedRecommendation, error) {

	// Check cache first
	cacheKey := fmt.Sprintf("recommendation:%s:%s", ticker, userTier)
	var cachedRecommendation *model.AggregatedRecommendation

	if err := uc.cache.Get(ctx, cacheKey, &cachedRecommendation); err == nil {
		uc.logger.Debug("Serving single recommendation from cache", "ticker", ticker)
		return cachedRecommendation, nil
	}

	// Generate fresh recommendation
	basicRecommendation, err := uc.basicCalculator.CalculateAggregatedRecommendation(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate recommendation for %s: %w", ticker, err)
	}

	// Enrich based on user tier
	finalRecommendation, err := uc.enrichForTier(ctx, basicRecommendation, userTier)
	if err != nil {
		uc.logger.Warn("Failed to enrich recommendation", "ticker", ticker, "error", err)
		// Return basic recommendation if enrichment fails
		finalRecommendation = basicRecommendation
	}

	// Cache the result
	cacheTTL := uc.getCacheTTLForTier(userTier)
	if err := uc.cache.Set(ctx, cacheKey, finalRecommendation, cacheTTL); err != nil {
		uc.logger.Warn("Failed to cache single recommendation", "ticker", ticker, "error", err)
	}

	return finalRecommendation, nil
}

// generateFreshRecommendations creates new recommendations from scratch
func (uc *TieredRecommendationUseCase) generateFreshRecommendations(
	ctx context.Context,
	userTier enums.UserTier,
	limit int,
	filters RecommendationFilters,
) ([]*model.AggregatedRecommendation, error) {

	// 1. Get tickers with recent activity
	since := time.Now().AddDate(0, 0, -30) // Last 30 days
	tickerEvents, err := uc.stockRepo.GetRecentByTickers(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tickers: %w", err)
	}

	var recommendations []*model.AggregatedRecommendation

	// 2. Process each ticker
	for ticker, events := range tickerEvents {
		if len(events) == 0 {
			continue
		}

		// Apply filters
		if uc.shouldSkipTicker(ticker, filters) {
			continue
		}

		// Calculate basic recommendation
		basicRec, err := uc.basicCalculator.CalculateAggregatedRecommendation(ctx, ticker)
		if err != nil {
			uc.logger.Warn("Failed to calculate basic recommendation",
				"ticker", ticker,
				"error", err)
			continue
		}

		// Apply score filter
		if filters.MinScore != nil && basicRec.BasicScore < *filters.MinScore {
			continue
		}

		// Enrich based on user tier
		finalRec, err := uc.enrichForTier(ctx, basicRec, userTier)
		if err != nil {
			uc.logger.Warn("Failed to enrich recommendation",
				"ticker", ticker,
				"error", err)
			// Use basic recommendation if enrichment fails
			recommendations = append(recommendations, basicRec)
			continue
		}

		// Apply recommendation type filter
		if filters.RecommendationType != nil && finalRec.RecommendationType != *filters.RecommendationType {
			continue
		}

		recommendations = append(recommendations, finalRec)
	}

	// 3. Sort by score (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].BasicScore > recommendations[j].BasicScore
	})

	// 4. Apply limit
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// enrichForTier enriches a recommendation based on user tier
func (uc *TieredRecommendationUseCase) enrichForTier(
	ctx context.Context,
	baseRecommendation *model.AggregatedRecommendation,
	userTier enums.UserTier,
) (*model.AggregatedRecommendation, error) {

	switch userTier {
	case enums.TIER_GUEST:
		// Guests get only basic recommendations
		return baseRecommendation, nil

	case enums.TIER_BASIC:
		// Basic users get external data enrichment
		return uc.externalDataEnricher.EnrichRecommendation(ctx, baseRecommendation)

	case enums.TIER_PREMIUM:
		// Premium users get external data + AI insights (prepared for Phase 6)
		enriched, err := uc.externalDataEnricher.EnrichRecommendation(ctx, baseRecommendation)
		if err != nil {
			return baseRecommendation, err
		}

		// TODO: Add AI insights enrichment in Phase 6
		enriched.Tier = enums.RECOMMENDATION_TIER_PREMIUM
		return enriched, nil

	default:
		uc.logger.Warn("Unknown user tier", "tier", userTier)
		return baseRecommendation, nil
	}
}

// Helper methods

// applyTierLimits enforces tier-based limits on the number of recommendations
func (uc *TieredRecommendationUseCase) applyTierLimits(userTier enums.UserTier, requestedLimit int) int {
	switch userTier {
	case enums.TIER_GUEST:
		return min(requestedLimit, 10) // Max 10 for guests
	case enums.TIER_BASIC:
		return min(requestedLimit, 25) // Max 25 for basic users
	case enums.TIER_PREMIUM:
		return min(requestedLimit, 100) // Max 100 for premium users
	default:
		return min(requestedLimit, 5)
	}
}

// getCacheTTLForTier returns the appropriate cache TTL based on user tier
func (uc *TieredRecommendationUseCase) getCacheTTLForTier(userTier enums.UserTier) time.Duration {
	switch userTier {
	case enums.TIER_GUEST:
		return 12 * time.Hour // Basic data changes less frequently
	case enums.TIER_BASIC:
		return 4 * time.Hour // External data updates more frequently
	case enums.TIER_PREMIUM:
		return 2 * time.Hour // Premium users get fresher data
	default:
		return 12 * time.Hour
	}
}

// getAvailableFeatures returns the features available to a user tier
func (uc *TieredRecommendationUseCase) getAvailableFeatures(userTier enums.UserTier) []string {
	features := []string{"basic_recommendations", "market_analytics"}

	if userTier == enums.TIER_BASIC || userTier == enums.TIER_PREMIUM {
		features = append(features, "real_time_data", "external_apis", "advanced_analytics")
	}

	if userTier == enums.TIER_PREMIUM {
		features = append(features, "ai_insights", "sentiment_analysis", "advanced_recommendations")
	}

	return features
}

// buildCacheKey creates a unique cache key for the request
func (uc *TieredRecommendationUseCase) buildCacheKey(userTier enums.UserTier, limit int, filters RecommendationFilters) string {
	// Simple cache key - in production this would be more sophisticated
	return fmt.Sprintf("recommendations:%s:%d", userTier, limit)
}

// tryGetFromCache attempts to retrieve recommendations from cache
func (uc *TieredRecommendationUseCase) tryGetFromCache(ctx context.Context, cacheKey string) ([]*model.AggregatedRecommendation, bool) {
	var cachedRecommendations []*model.AggregatedRecommendation

	if err := uc.cache.Get(ctx, cacheKey, &cachedRecommendations); err == nil {
		return cachedRecommendations, true
	}

	return nil, false
}

// shouldSkipTicker checks if a ticker should be skipped based on filters
func (uc *TieredRecommendationUseCase) shouldSkipTicker(ticker string, filters RecommendationFilters) bool {
	// Check exclude list
	for _, excludeTicker := range filters.ExcludeTickers {
		if ticker == excludeTicker {
			return true
		}
	}

	// TODO: Add sector filtering logic here when sector mapping is available

	return false
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
