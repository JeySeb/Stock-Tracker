package usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	recommendationModel "stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/pkg/logger"
)

// TieredRecommendationUseCase handles business logic for tiered recommendations
type TieredRecommendationUseCase struct {
	stockRepo            repositories.StockRepository
	basicCalculator      BasicScoringCalculatorInterface
	externalDataEnricher ExternalDataEnricherInterface
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
	Data []*recommendationModel.AggregatedRecommendation `json:"data"`
	Meta RecommendationMeta                              `json:"meta"`
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
	basicCalculator BasicScoringCalculatorInterface,
	externalDataEnricher ExternalDataEnricherInterface,
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
) (*recommendationModel.AggregatedRecommendation, error) {

	// Check cache first
	cacheKey := fmt.Sprintf("recommendation:%s:%s", ticker, userTier)
	var cachedRecommendation *recommendationModel.AggregatedRecommendation

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

// generateFreshRecommendations creates new recommendations from scratch with optimizations
func (uc *TieredRecommendationUseCase) generateFreshRecommendations(
	ctx context.Context,
	userTier enums.UserTier,
	limit int,
	filters RecommendationFilters,
) ([]*recommendationModel.AggregatedRecommendation, error) {

	// 1. Get tickers with recent activity
	since := time.Now().AddDate(0, 0, -30) // Last 30 days
	tickerEvents, err := uc.stockRepo.GetRecentByTickers(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tickers: %w", err)
	}

	// 2. Apply early filtering and limit to avoid processing too many tickers
	eligibleTickers := uc.filterAndLimitTickers(tickerEvents, filters, limit*3) // Process 3x limit to ensure good results after filtering

	if len(eligibleTickers) == 0 {
		return []*recommendationModel.AggregatedRecommendation{}, nil
	}

	// 3. Use optimized batch processing with the data we already have
	recommendations, err := uc.generateRecommendationsFromEvents(ctx, eligibleTickers, userTier, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations from events: %w", err)
	}

	// 4. Sort by score (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].BasicScore > recommendations[j].BasicScore
	})

	// 5. Apply final limit
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	uc.logger.Info("Generated optimized recommendations",
		"processed_tickers", len(eligibleTickers),
		"final_recommendations", len(recommendations))

	return recommendations, nil
}

// filterAndLimitTickers applies early filtering and limits the number of tickers to process
func (uc *TieredRecommendationUseCase) filterAndLimitTickers(
	tickerEvents map[string][]*stockModel.Stock,
	filters RecommendationFilters,
	maxTickers int,
) map[string][]*stockModel.Stock {

	result := make(map[string][]*stockModel.Stock)
	count := 0

	// Sort tickers by activity level (number of recent events)
	type tickerActivity struct {
		ticker string
		events []*stockModel.Stock
		score  float64 // Simple activity score
	}

	var activities []tickerActivity
	for ticker, events := range tickerEvents {
		if len(events) == 0 || uc.shouldSkipTicker(ticker, filters) {
			continue
		}

		// Calculate simple activity score based on recency and count
		score := float64(len(events))
		if len(events) > 0 {
			daysSince := time.Since(events[0].EventTime).Hours() / 24
			recencyBonus := math.Max(0, 30-daysSince) / 30 // Bonus for recent activity
			score += recencyBonus * 10
		}

		activities = append(activities, tickerActivity{
			ticker: ticker,
			events: events,
			score:  score,
		})
	}

	// Sort by activity score (descending)
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].score > activities[j].score
	})

	// Take top N most active tickers
	for _, activity := range activities {
		if count >= maxTickers {
			break
		}
		result[activity.ticker] = activity.events
		count++
	}

	return result
}

// generateRecommendationsFromEvents optimized to use existing event data instead of re-querying
func (uc *TieredRecommendationUseCase) generateRecommendationsFromEvents(
	ctx context.Context,
	tickerEvents map[string][]*stockModel.Stock,
	userTier enums.UserTier,
	filters RecommendationFilters,
) ([]*recommendationModel.AggregatedRecommendation, error) {

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Process tickers with concurrency for better performance
	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)
	defer close(semaphore) // Ensure semaphore channel is closed

	var recommendations []*recommendationModel.AggregatedRecommendation
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Channel for collecting errors from goroutines
	errorChan := make(chan error, len(tickerEvents))
	defer close(errorChan)

	// Create a context with timeout for the goroutines
	goroutineCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel() // Ensure context is cancelled to free resources

	for ticker, events := range tickerEvents {
		// Check context cancellation before starting new goroutine
		select {
		case <-goroutineCtx.Done():
			uc.logger.Warn("Context cancelled, stopping goroutine creation")
			wg.Wait() // Wait for any running goroutines to finish
			return recommendations, goroutineCtx.Err()
		default:
		}

		wg.Add(1)

		go func(ticker string, events []*stockModel.Stock) {
			defer func() {
				wg.Done()
				// Recover from any panics in goroutines
				if r := recover(); r != nil {
					uc.logger.Error("Goroutine panic recovered",
						"ticker", ticker,
						"panic", r)
					errorChan <- fmt.Errorf("goroutine panic for ticker %s: %v", ticker, r)
				}
			}()

			// Acquire semaphore with context cancellation check
			select {
			case semaphore <- struct{}{}:
				defer func() {
					select {
					case <-semaphore:
					default:
					}
				}()
			case <-goroutineCtx.Done():
				uc.logger.Debug("Goroutine cancelled before acquiring semaphore", "ticker", ticker)
				return
			}

			// Check context again after acquiring semaphore
			select {
			case <-goroutineCtx.Done():
				uc.logger.Debug("Goroutine cancelled after acquiring semaphore", "ticker", ticker)
				return
			default:
			}

			// Calculate recommendation using the existing events data
			basicRec, err := uc.basicCalculator.CalculateAggregatedRecommendationFromEvents(goroutineCtx, ticker, events)
			if err != nil {
				uc.logger.Warn("Failed to calculate basic recommendation",
					"ticker", ticker,
					"error", err)
				select {
				case errorChan <- fmt.Errorf("calculation failed for ticker %s: %w", ticker, err):
				default:
				}
				return
			}

			// Apply score filter
			if filters.MinScore != nil && basicRec.BasicScore < *filters.MinScore {
				return
			}

			// Enrich based on user tier with context propagation
			finalRec, err := uc.enrichForTier(goroutineCtx, basicRec, userTier)
			if err != nil {
				uc.logger.Warn("Failed to enrich recommendation",
					"ticker", ticker,
					"error", err)
				// Use basic recommendation if enrichment fails
				finalRec = basicRec
			}

			// Apply recommendation type filter
			if filters.RecommendationType != nil && finalRec.RecommendationType != *filters.RecommendationType {
				return
			}

			// Thread-safe append with context check
			select {
			case <-goroutineCtx.Done():
				uc.logger.Debug("Goroutine cancelled before adding recommendation", "ticker", ticker)
				return
			default:
			}

			mutex.Lock()
			recommendations = append(recommendations, finalRec)
			mutex.Unlock()
		}(ticker, events)
	}

	// Wait for all goroutines to complete or timeout
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		// All goroutines completed successfully
		uc.logger.Debug("All goroutines completed successfully",
			"processed_tickers", len(tickerEvents),
			"recommendations_generated", len(recommendations))
	case <-goroutineCtx.Done():
		// Context cancelled or timeout
		uc.logger.Warn("Goroutines cancelled due to timeout or context cancellation",
			"error", goroutineCtx.Err(),
			"partial_recommendations", len(recommendations))
		// Wait a bit more for graceful shutdown
		time.Sleep(100 * time.Millisecond)
	}

	// Check for any errors from goroutines
	select {
	case err := <-errorChan:
		uc.logger.Warn("Goroutine reported error", "error", err)
		// Don't fail the entire operation for individual ticker failures
	default:
	}

	return recommendations, nil
}

// enrichForTier enriches a recommendation based on user tier
func (uc *TieredRecommendationUseCase) enrichForTier(
	ctx context.Context,
	baseRecommendation *recommendationModel.AggregatedRecommendation,
	userTier enums.UserTier,
) (*recommendationModel.AggregatedRecommendation, error) {

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
		features = append(features, "ai_insights", "sentiment_analysis", "advanced_recommendations", "priority_support")
	}

	return features
}

// buildCacheKey creates a unique cache key for the request
func (uc *TieredRecommendationUseCase) buildCacheKey(userTier enums.UserTier, limit int, filters RecommendationFilters) string {
	// Simple cache key - in production this would be more sophisticated
	return fmt.Sprintf("recommendations:%s:%d", userTier, limit)
}

// tryGetFromCache attempts to retrieve recommendations from cache
func (uc *TieredRecommendationUseCase) tryGetFromCache(ctx context.Context, cacheKey string) ([]*recommendationModel.AggregatedRecommendation, bool) {
	var cachedRecommendations []*recommendationModel.AggregatedRecommendation

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
