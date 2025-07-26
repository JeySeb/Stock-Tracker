package recommendation

import (
	"context"
	"fmt"
	"math"
	"time"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	"stock-tracker/internal/infrastructure/cache"
	"stock-tracker/internal/infrastructure/external"
	"stock-tracker/pkg/logger"
)

// ExternalDataEnricher enriches basic recommendations with external data
type ExternalDataEnricher struct {
	yahooClient        external.YahooFinanceClient
	alphaVantageClient external.AlphaVantageClient
	cache              cache.Cache
	logger             logger.Logger
}

// EnrichmentResult contains the result of external data enrichment
type EnrichmentResult struct {
	EnrichedScore     float64
	ExternalData      *model.ExternalStockData
	EnrichmentFactors []model.ScoringFactor
}

// NewExternalDataEnricher creates a new instance of the external data enricher
func NewExternalDataEnricher(
	yahooClient external.YahooFinanceClient,
	alphaVantageClient external.AlphaVantageClient,
	cache cache.Cache,
	logger logger.Logger,
) *ExternalDataEnricher {
	return &ExternalDataEnricher{
		yahooClient:        yahooClient,
		alphaVantageClient: alphaVantageClient,
		cache:              cache,
		logger:             logger,
	}
}

// EnrichRecommendation enhances a basic recommendation with external data
func (enricher *ExternalDataEnricher) EnrichRecommendation(
	ctx context.Context,
	baseRecommendation *model.AggregatedRecommendation,
) (*model.AggregatedRecommendation, error) {

	// 1. Get real-time external data
	externalData, err := enricher.getRealtimeData(ctx, baseRecommendation.Ticker)
	if err != nil {
		enricher.logger.Warn("Failed to get external data",
			"ticker", baseRecommendation.Ticker,
			"error", err)
		// Continue without external data - graceful degradation
		return baseRecommendation, nil
	}

	// 2. Calculate enriched factors
	enrichmentResult := enricher.calculateEnrichedFactors(baseRecommendation, externalData)

	// 3. Create enriched recommendation
	enrichedRecommendation := *baseRecommendation // Copy the base
	enrichedRecommendation.Tier = enums.RECOMMENDATION_TIER_ENRICHED
	enrichedRecommendation.ExternalData = externalData
	enrichedRecommendation.BasicScore = enrichmentResult.EnrichedScore
	enrichedRecommendation.ExpiresAt = time.Now().Add(4 * time.Hour) // More frequent updates for real-time data

	// 4. Merge scoring factors
	enrichedRecommendation.ScoringFactors = append(
		enrichedRecommendation.ScoringFactors,
		enrichmentResult.EnrichmentFactors...,
	)

	enrichedRecommendation.DetermineType()

	enricher.logger.Info("Recommendation enriched with external data",
		"ticker", baseRecommendation.Ticker,
		"original_score", baseRecommendation.BasicScore,
		"enriched_score", enrichedRecommendation.BasicScore)

	return &enrichedRecommendation, nil
}

// getRealtimeData retrieves real-time data with caching
func (enricher *ExternalDataEnricher) getRealtimeData(ctx context.Context, ticker string) (*model.ExternalStockData, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("external_data:%s", ticker)
	var cachedData *model.ExternalStockData

	if err := enricher.cache.Get(ctx, cacheKey, &cachedData); err == nil {
		enricher.logger.Debug("Serving external data from cache", "ticker", ticker)
		return cachedData, nil
	}

	// Fetch from external API (prefer Yahoo Finance for reliability)
	data, err := enricher.yahooClient.GetQuote(ctx, ticker)
	if err != nil {
		enricher.logger.Warn("Yahoo Finance failed, trying Alpha Vantage",
			"ticker", ticker,
			"error", err)

		// Fallback to Alpha Vantage
		data, err = enricher.alphaVantageClient.GetQuote(ctx, ticker)
		if err != nil {
			return nil, fmt.Errorf("all external data sources failed for %s: %w", ticker, err)
		}
	}

	// Cache for 5 minutes (balance between freshness and API limits)
	if err := enricher.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil {
		enricher.logger.Warn("Failed to cache external data", "ticker", ticker, "error", err)
	}

	enricher.logger.Debug("External data fetched and cached", "ticker", ticker)
	return data, nil
}

// calculateEnrichedFactors computes enrichment factors based on external data
func (enricher *ExternalDataEnricher) calculateEnrichedFactors(
	base *model.AggregatedRecommendation,
	external *model.ExternalStockData,
) EnrichmentResult {

	// Define weights for combining basic score with external factors
	baseWeight := 0.70     // 70% weight to basic score
	externalWeight := 0.30 // 30% weight to external factors

	// Calculate individual external factors
	upsideFactor := enricher.calculateUpsidePotential(base, external)
	volumeFactor := enricher.calculateVolumeActivity(external)
	positionFactor := enricher.calculateYearPosition(external)

	// Combine external factors
	combinedExternalScore := (upsideFactor + volumeFactor + positionFactor) / 3.0

	// Calculate final enriched score
	enrichedScore := baseWeight*base.BasicScore + externalWeight*combinedExternalScore

	// Create enrichment factors for transparency
	enrichmentFactors := []model.ScoringFactor{
		{
			Name:        "Real-time Upside Potential",
			Score:       upsideFactor,
			Weight:      0.15,
			Explanation: enricher.explainUpsidePotential(base, external, upsideFactor),
			Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
		},
		{
			Name:        "Volume Activity",
			Score:       volumeFactor,
			Weight:      0.10,
			Explanation: enricher.explainVolumeActivity(external, volumeFactor),
			Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
		},
		{
			Name:        "52-Week Position",
			Score:       positionFactor,
			Weight:      0.05,
			Explanation: enricher.explainYearPosition(external, positionFactor),
			Tier:        enums.RECOMMENDATION_TIER_ENRICHED,
		},
	}

	return EnrichmentResult{
		EnrichedScore:     math.Min(1.0, math.Max(0.0, enrichedScore)),
		ExternalData:      external,
		EnrichmentFactors: enrichmentFactors,
	}
}

// calculateUpsidePotential calculates upside potential vs current market price
func (enricher *ExternalDataEnricher) calculateUpsidePotential(
	base *model.AggregatedRecommendation,
	external *model.ExternalStockData,
) float64 {

	if base.LatestTargetPrice <= 0 || external.CurrentPrice <= 0 {
		return 0.5 // neutral if no data available
	}

	// Calculate upside percentage
	upside := (base.LatestTargetPrice - external.CurrentPrice) / external.CurrentPrice

	// Normalize upside to 0-1 scale
	// 0% upside = 0.5, 20% upside = 0.8, 50% upside = 1.0
	// Negative upside gets lower scores
	normalizedUpside := 0.5 + math.Min(0.5, math.Max(-0.5, upside))

	return math.Max(0.0, math.Min(1.0, normalizedUpside))
}

// calculateVolumeActivity calculates volume activity score
func (enricher *ExternalDataEnricher) calculateVolumeActivity(external *model.ExternalStockData) float64 {
	if external.Volume <= 0 || external.AvgVolume == nil || *external.AvgVolume <= 0 {
		return 0.5 // neutral if no volume data
	}

	// Calculate volume ratio vs average
	volumeRatio := float64(external.Volume) / float64(*external.AvgVolume)

	// Normalize volume ratio
	// 1x average = 0.5, 2x average = 0.8, 3x+ average = 1.0
	// Below average gets lower scores
	switch {
	case volumeRatio >= 3.0:
		return 1.0
	case volumeRatio >= 2.0:
		return 0.8
	case volumeRatio >= 1.0:
		return 0.5 + (volumeRatio-1.0)*0.3 // Linear scale from 0.5 to 0.8
	case volumeRatio >= 0.5:
		return 0.25 + (volumeRatio-0.5)*0.5 // Linear scale from 0.25 to 0.5
	default:
		return 0.1 // Very low volume
	}
}

// calculateYearPosition calculates position within 52-week range
func (enricher *ExternalDataEnricher) calculateYearPosition(external *model.ExternalStockData) float64 {
	if external.Week52High == nil || external.Week52Low == nil ||
		*external.Week52High <= *external.Week52Low {
		return 0.5 // neutral if no 52-week data
	}

	// Calculate position in 52-week range
	position := (external.CurrentPrice - *external.Week52Low) / (*external.Week52High - *external.Week52Low)

	// Position directly maps to score (0-1)
	return math.Max(0.0, math.Min(1.0, position))
}

// Explanation methods for enrichment factors
func (enricher *ExternalDataEnricher) explainUpsidePotential(
	base *model.AggregatedRecommendation,
	external *model.ExternalStockData,
	score float64,
) string {

	if base.LatestTargetPrice <= 0 || external.CurrentPrice <= 0 {
		return "No price target or current price data available"
	}

	upside := (base.LatestTargetPrice - external.CurrentPrice) / external.CurrentPrice * 100

	if upside > 20 {
		return fmt.Sprintf("Strong upside potential: %.1f%% to target ($%.2f → $%.2f)",
			upside, external.CurrentPrice, base.LatestTargetPrice)
	} else if upside > 0 {
		return fmt.Sprintf("Modest upside: %.1f%% to target ($%.2f → $%.2f)",
			upside, external.CurrentPrice, base.LatestTargetPrice)
	} else {
		return fmt.Sprintf("Trading above target: %.1f%% above target price", -upside)
	}
}

func (enricher *ExternalDataEnricher) explainVolumeActivity(external *model.ExternalStockData, score float64) string {
	if external.Volume <= 0 || external.AvgVolume == nil || *external.AvgVolume <= 0 {
		return "Volume data not available"
	}

	volumeRatio := float64(external.Volume) / float64(*external.AvgVolume)

	switch {
	case volumeRatio >= 3.0:
		return fmt.Sprintf("Very high volume activity: %.1fx average volume", volumeRatio)
	case volumeRatio >= 2.0:
		return fmt.Sprintf("High volume activity: %.1fx average volume", volumeRatio)
	case volumeRatio >= 1.0:
		return fmt.Sprintf("Above average volume: %.1fx average", volumeRatio)
	case volumeRatio >= 0.5:
		return fmt.Sprintf("Below average volume: %.1fx average", volumeRatio)
	default:
		return "Very low volume activity"
	}
}

func (enricher *ExternalDataEnricher) explainYearPosition(external *model.ExternalStockData, score float64) string {
	if external.Week52High == nil || external.Week52Low == nil {
		return "52-week range data not available"
	}

	position := (external.CurrentPrice - *external.Week52Low) / (*external.Week52High - *external.Week52Low) * 100

	switch {
	case position >= 80:
		return fmt.Sprintf("Near 52-week high: %.1f%% of range", position)
	case position >= 60:
		return fmt.Sprintf("Above midpoint of 52-week range: %.1f%%", position)
	case position >= 40:
		return fmt.Sprintf("Mid-range of 52-week prices: %.1f%%", position)
	case position >= 20:
		return fmt.Sprintf("Below midpoint of 52-week range: %.1f%%", position)
	default:
		return fmt.Sprintf("Near 52-week low: %.1f%% of range", position)
	}
}

// GetEnrichmentPreview provides a preview of what enrichment would add (for testing)
func (enricher *ExternalDataEnricher) GetEnrichmentPreview(
	ctx context.Context,
	ticker string,
) (*model.ExternalStockData, error) {
	return enricher.getRealtimeData(ctx, ticker)
}
