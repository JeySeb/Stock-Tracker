package validation

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/pkg/logger"
)

// ScoringConfig holds configuration for the scoring calculator
type ScoringConfig struct {
	// Cache configuration
	CacheTTL time.Duration `json:"cache_ttl"`

	// Data filtering
	DataRetentionDays int `json:"data_retention_days"`

	// Scoring weights
	Weights model.ScoringWeights `json:"weights"`

	// Recency scoring thresholds
	RecencyThresholds RecencyThresholds `json:"recency_thresholds"`

	// Target movement scoring limits
	TargetMovementCap float64 `json:"target_movement_cap"`

	// Neutral scores for missing data
	NeutralScore float64 `json:"neutral_score"`

	// Confidence calculation
	ConfidenceThresholds ConfidenceThresholds `json:"confidence_thresholds"`
}

// RecencyThresholds defines scoring thresholds for event recency
type RecencyThresholds struct {
	VeryRecent      float64 `json:"very_recent_days"`  // <= 1 day
	Recent          float64 `json:"recent_days"`       // <= 7 days
	Medium          float64 `json:"medium_days"`       // <= 30 days
	VeryRecentScore float64 `json:"very_recent_score"` // 1.0
	RecentScore     float64 `json:"recent_score"`      // 0.8
	MediumScore     float64 `json:"medium_score"`      // 0.6
	OldScore        float64 `json:"old_score"`         // 0.2
	DecayRate       float64 `json:"decay_rate"`        // 0.01
}

// ConfidenceThresholds defines confidence calculation parameters
type ConfidenceThresholds struct {
	MinEvents     int     `json:"min_events"`
	LowThreshold  int     `json:"low_threshold"`
	HighThreshold int     `json:"high_threshold"`
	MinConfidence float64 `json:"min_confidence"`
	MaxConfidence float64 `json:"max_confidence"`
}

// DefaultScoringConfig returns the default scoring configuration
func DefaultScoringConfig() *ScoringConfig {
	return &ScoringConfig{
		CacheTTL:          30 * time.Minute,
		DataRetentionDays: 90,
		Weights: model.ScoringWeights{
			BrokerFrequency: 0.25,
			TargetMovement:  0.25,
			RatingChange:    0.25,
			Recency:         0.15,
			Consensus:       0.10,
		},
		RecencyThresholds: RecencyThresholds{
			VeryRecent:      1.0,
			Recent:          7.0,
			Medium:          30.0,
			VeryRecentScore: 1.0,
			RecentScore:     0.8,
			MediumScore:     0.6,
			OldScore:        0.2,
			DecayRate:       0.01,
		},
		TargetMovementCap: 0.5, // ±50%
		NeutralScore:      0.5,
		ConfidenceThresholds: ConfidenceThresholds{
			MinEvents:     1,
			LowThreshold:  3,
			HighThreshold: 10,
			MinConfidence: 0.3,
			MaxConfidence: 0.95,
		},
	}
}

// BasicScoringCalculator calculates recommendation scores without hardcoded values
type BasicScoringCalculator struct {
	stockRepo        repositories.StockRepository
	logger           logger.Logger
	config           *ScoringConfig
	brokerStatsCache map[string]float64
	cacheMutex       sync.RWMutex
	cacheExpiry      time.Time
}

// BasicScoringFactors represents the different factors used in scoring
type BasicScoringFactors struct {
	BrokerFrequencyScore float64 `json:"broker_frequency_score"` // Calculated dynamically
	TargetMovementScore  float64 `json:"target_movement_score"`  // Based on TargetFrom/TargetTo
	RatingChangeScore    float64 `json:"rating_change_score"`    // Based on RatingFrom/RatingTo
	RecencyScore         float64 `json:"recency_score"`          // Based on EventTime
	ConsensusScore       float64 `json:"consensus_score"`        // Aggregation of events by ticker
}

// NewBasicScoringCalculator creates a new instance of the scoring calculator with optional configuration
func NewBasicScoringCalculator(
	stockRepo repositories.StockRepository,
	logger logger.Logger,
	config ...*ScoringConfig,
) *BasicScoringCalculator {
	scoringConfig := DefaultScoringConfig()
	if len(config) > 0 && config[0] != nil {
		scoringConfig = config[0]
	}

	return &BasicScoringCalculator{
		stockRepo:        stockRepo,
		logger:           logger,
		config:           scoringConfig,
		brokerStatsCache: make(map[string]float64),
	}
}

// CalculateAggregatedRecommendation calculates a recommendation based on all events for a ticker
func (calc *BasicScoringCalculator) CalculateAggregatedRecommendation(
	ctx context.Context,
	ticker string,
) (*model.AggregatedRecommendation, error) {

	// 1. Get ALL events for this ticker using configurable retention period
	since := time.Now().AddDate(0, 0, -calc.config.DataRetentionDays)
	events, err := calc.stockRepo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for ticker %s: %w", ticker, err)
	}

	// Filter by date
	recentEvents := calc.filterRecentEvents(events, since)
	if len(recentEvents) == 0 {
		return nil, fmt.Errorf("no recent events for ticker %s", ticker)
	}

	return calc.CalculateAggregatedRecommendationFromEvents(ctx, ticker, recentEvents)
}

// CalculateAggregatedRecommendationFromEvents calculates a recommendation from provided events (optimized version)
func (calc *BasicScoringCalculator) CalculateAggregatedRecommendationFromEvents(
	ctx context.Context,
	ticker string,
	events []*stockModel.Stock,
) (*model.AggregatedRecommendation, error) {

	if len(events) == 0 {
		return nil, fmt.Errorf("no events provided for ticker %s", ticker)
	}

	// Filter to configurable retention period for consistency
	since := time.Now().AddDate(0, 0, -calc.config.DataRetentionDays)
	recentEvents := calc.filterRecentEvents(events, since)
	if len(recentEvents) == 0 {
		return nil, fmt.Errorf("no recent events for ticker %s", ticker)
	}

	// 2. Calculate aggregated statistics
	stats := calc.calculateEventStatistics(recentEvents)

	// 3. Calculate scores based on real data using configuration
	brokerFreqScore := calc.calculateBrokerFrequencyScore(ctx, recentEvents)
	targetScore := calc.calculateTargetMovementScore(recentEvents)
	ratingScore := calc.calculateRatingChangeScore(recentEvents)
	recencyScore := calc.calculateRecencyScore(recentEvents)
	consensusScore := calc.calculateConsensusScore(recentEvents)

	// 4. Combine scores with configurable weights
	scoringFactors := BasicScoringFactors{
		BrokerFrequencyScore: brokerFreqScore,
		TargetMovementScore:  targetScore,
		RatingChangeScore:    ratingScore,
		RecencyScore:         recencyScore,
		ConsensusScore:       consensusScore,
	}

	finalScore := calc.combineScores(scoringFactors, calc.config.Weights)

	// 5. Create aggregated recommendation
	recommendation := model.NewAggregatedRecommendation(
		ticker,
		recentEvents[0].Company,
		enums.RECOMMENDATION_TIER_BASIC,
	)

	recommendation.TotalEvents = len(recentEvents)
	recommendation.PositiveEvents = stats.PositiveCount
	recommendation.NegativeEvents = stats.NegativeCount
	recommendation.AvgTargetChange = stats.AvgTargetChange
	recommendation.LatestTargetPrice = stats.LatestTargetPrice
	recommendation.BrokerConsensus = consensusScore
	recommendation.BasicScore = finalScore
	recommendation.Confidence = calc.calculateConfidence(len(recentEvents))
	recommendation.LastEventTime = stats.LatestEventTime
	recommendation.ScoringFactors = calc.createScoringFactorsDetails(scoringFactors, calc.config.Weights)

	recommendation.DetermineType()

	return recommendation, nil
}

// filterRecentEvents filters events to only include those after the given time
func (calc *BasicScoringCalculator) filterRecentEvents(events []*stockModel.Stock, since time.Time) []*stockModel.Stock {
	var filtered []*stockModel.Stock
	for _, event := range events {
		if event.EventTime.After(since) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// calculateEventStatistics computes statistical data from events
func (calc *BasicScoringCalculator) calculateEventStatistics(events []*stockModel.Stock) model.EventStatistics {
	var stats model.EventStatistics
	var totalTargetChange float64
	var targetChangeCount int
	var latestEvent time.Time

	stats.TotalEvents = len(events)

	for _, event := range events {
		// Count positive/negative events
		if event.IsRecommendation() {
			stats.PositiveCount++
		} else {
			stats.NegativeCount++
		}

		// Calculate average target change
		targetChange := event.GetPriceTargetChange()
		if targetChange != 0 {
			totalTargetChange += targetChange
			targetChangeCount++
		}

		// Track latest target price and event time
		if event.TargetTo > 0 && (stats.LatestTargetPrice == 0 || event.EventTime.After(latestEvent)) {
			stats.LatestTargetPrice = event.TargetTo
			latestEvent = event.EventTime
		}

		if event.EventTime.After(stats.LatestEventTime) {
			stats.LatestEventTime = event.EventTime
		}
	}

	if targetChangeCount > 0 {
		stats.AvgTargetChange = totalTargetChange / float64(targetChangeCount)
	}

	return stats
}

// ⚡ OPTIMIZED CALCULATION WITH CACHING
// calculateBrokerFrequencyScore calculates broker credibility using cached stats
func (calc *BasicScoringCalculator) calculateBrokerFrequencyScore(
	ctx context.Context,
	events []*stockModel.Stock,
) float64 {
	// Get cached broker frequency map
	brokerFrequency, err := calc.getCachedBrokerStats(ctx)
	if err != nil {
		calc.logger.Warn("Failed to get cached brokerage stats", "error", err)
		return calc.config.NeutralScore // Use configured neutral score
	}

	if len(brokerFrequency) == 0 {
		return calc.config.NeutralScore // Use configured neutral score
	}

	// Calculate total reports for normalization
	var totalReports float64
	for _, freq := range brokerFrequency {
		totalReports += freq
	}

	if totalReports == 0 {
		return calc.config.NeutralScore
	}

	// Calculate weighted score based on frequency
	var weightedScore float64
	var totalWeight float64

	for _, event := range events {
		freq := brokerFrequency[event.Brokerage]
		if freq == 0 {
			freq = 1 // At least 1 for new brokers
		}

		// Normalize frequency (log to avoid outliers)
		normalizedFreq := math.Log(1+freq) / math.Log(1+totalReports)
		weight := normalizedFreq

		weightedScore += normalizedFreq * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return calc.config.NeutralScore
	}

	return math.Min(1.0, weightedScore/totalWeight)
}

// getCachedBrokerStats returns cached broker frequency map, refreshing if needed
func (calc *BasicScoringCalculator) getCachedBrokerStats(ctx context.Context) (map[string]float64, error) {
	calc.cacheMutex.RLock()
	if len(calc.brokerStatsCache) > 0 && time.Now().Before(calc.cacheExpiry) {
		// Return cached data
		result := make(map[string]float64)
		for k, v := range calc.brokerStatsCache {
			result[k] = v
		}
		calc.cacheMutex.RUnlock()
		return result, nil
	}
	calc.cacheMutex.RUnlock()

	// Need to refresh cache
	calc.cacheMutex.Lock()
	defer calc.cacheMutex.Unlock()

	// Double-check in case another goroutine updated while waiting for lock
	if len(calc.brokerStatsCache) > 0 && time.Now().Before(calc.cacheExpiry) {
		result := make(map[string]float64)
		for k, v := range calc.brokerStatsCache {
			result[k] = v
		}
		return result, nil
	}

	// Refresh from database
	brokerStats, err := calc.stockRepo.GetBrokerageStats(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	calc.brokerStatsCache = make(map[string]float64)
	for _, stat := range brokerStats {
		calc.brokerStatsCache[stat.Brokerage] = float64(stat.Count)
	}
	calc.cacheExpiry = time.Now().Add(calc.config.CacheTTL)

	// Return copy of cached data
	result := make(map[string]float64)
	for k, v := range calc.brokerStatsCache {
		result[k] = v
	}
	return result, nil
}

// calculateTargetMovementScore calculates score based on price target changes
func (calc *BasicScoringCalculator) calculateTargetMovementScore(events []*stockModel.Stock) float64 {
	var totalMovement float64
	var count int

	for _, event := range events {
		if event.TargetFrom > 0 && event.TargetTo > 0 {
			change := (event.TargetTo - event.TargetFrom) / event.TargetFrom
			totalMovement += change
			count++
		}
	}

	if count == 0 {
		return calc.config.NeutralScore // Use configured neutral score
	}

	avgMovement := totalMovement / float64(count)

	// Normalize using configurable cap instead of hardcoded 50%
	normalizedChange := math.Max(-calc.config.TargetMovementCap,
		math.Min(calc.config.TargetMovementCap, avgMovement))
	score := calc.config.NeutralScore + normalizedChange // Maps -cap to 0, 0% to neutral, +cap to 1.0

	return math.Max(0, math.Min(1, score))
}

// calculateRatingChangeScore calculates score based on rating changes
func (calc *BasicScoringCalculator) calculateRatingChangeScore(events []*stockModel.Stock) float64 {
	var totalRatingChange float64
	var count int

	for _, event := range events {
		if event.RatingFrom != "" && event.RatingTo != "" {
			ratingChange := event.GetRatingChangeScore()
			totalRatingChange += ratingChange
			count++
		}
	}

	if count == 0 {
		return calc.config.NeutralScore // Use configured neutral score
	}

	avgRatingChange := totalRatingChange / float64(count)

	// Normalize rating change to 0-1 scale using configured neutral score
	return math.Max(0, math.Min(1, calc.config.NeutralScore+avgRatingChange))
}

// calculateRecencyScore calculates score based on how recent the events are
func (calc *BasicScoringCalculator) calculateRecencyScore(events []*stockModel.Stock) float64 {
	if len(events) == 0 {
		return 0
	}

	// Find the most recent event
	var mostRecent time.Time
	for _, event := range events {
		if event.EventTime.After(mostRecent) {
			mostRecent = event.EventTime
		}
	}

	daysSince := time.Since(mostRecent).Hours() / 24
	thresholds := calc.config.RecencyThresholds

	// Use configurable thresholds and scores
	switch {
	case daysSince <= thresholds.VeryRecent:
		return thresholds.VeryRecentScore
	case daysSince <= thresholds.Recent:
		return thresholds.RecentScore
	case daysSince <= thresholds.Medium:
		// Linear decay using configurable rate
		return thresholds.MediumScore - (daysSince-thresholds.Recent)*thresholds.DecayRate
	default:
		return thresholds.OldScore
	}
}

// calculateConsensusScore calculates consensus based on positive vs negative events
func (calc *BasicScoringCalculator) calculateConsensusScore(events []*stockModel.Stock) float64 {
	if len(events) == 0 {
		return calc.config.NeutralScore
	}

	positiveCount := 0
	for _, event := range events {
		if event.IsRecommendation() {
			positiveCount++
		}
	}

	// Calculate consensus as ratio of positive events
	consensus := float64(positiveCount) / float64(len(events))
	return consensus
}

// combineScores combines all scoring factors with their weights
func (calc *BasicScoringCalculator) combineScores(factors BasicScoringFactors, weights model.ScoringWeights) float64 {
	score := factors.BrokerFrequencyScore*weights.BrokerFrequency +
		factors.TargetMovementScore*weights.TargetMovement +
		factors.RatingChangeScore*weights.RatingChange +
		factors.RecencyScore*weights.Recency +
		factors.ConsensusScore*weights.Consensus

	return math.Min(1.0, math.Max(0.0, score))
}

// calculateConfidence calculates confidence level based on number of events
func (calc *BasicScoringCalculator) calculateConfidence(eventCount int) float64 {
	thresholds := calc.config.ConfidenceThresholds

	if eventCount < thresholds.MinEvents {
		return thresholds.MinConfidence
	}

	switch {
	case eventCount <= thresholds.LowThreshold:
		// Linear interpolation between min and medium confidence
		ratio := float64(eventCount-thresholds.MinEvents) / float64(thresholds.LowThreshold-thresholds.MinEvents)
		return thresholds.MinConfidence + ratio*(0.6-thresholds.MinConfidence)
	case eventCount <= thresholds.HighThreshold:
		// Linear interpolation between medium and high confidence
		ratio := float64(eventCount-thresholds.LowThreshold) / float64(thresholds.HighThreshold-thresholds.LowThreshold)
		return 0.6 + ratio*(0.8-0.6)
	default:
		// High confidence with diminishing returns
		return math.Min(thresholds.MaxConfidence, 0.8+0.15*math.Log(float64(eventCount-thresholds.HighThreshold+1)))
	}
}

// createScoringFactorsDetails creates detailed scoring factors for API response
func (calc *BasicScoringCalculator) createScoringFactorsDetails(factors BasicScoringFactors, weights model.ScoringWeights) []model.ScoringFactor {
	return []model.ScoringFactor{
		{
			Name:        "Broker Frequency",
			Score:       factors.BrokerFrequencyScore,
			Weight:      weights.BrokerFrequency,
			Explanation: calc.explainBrokerFrequency(factors.BrokerFrequencyScore),
			Tier:        enums.RECOMMENDATION_TIER_BASIC,
		},
		{
			Name:        "Target Movement",
			Score:       factors.TargetMovementScore,
			Weight:      weights.TargetMovement,
			Explanation: calc.explainTargetMovement(factors.TargetMovementScore),
			Tier:        enums.RECOMMENDATION_TIER_BASIC,
		},
		{
			Name:        "Rating Change",
			Score:       factors.RatingChangeScore,
			Weight:      weights.RatingChange,
			Explanation: calc.explainRatingChange(factors.RatingChangeScore),
			Tier:        enums.RECOMMENDATION_TIER_BASIC,
		},
		{
			Name:        "Recency",
			Score:       factors.RecencyScore,
			Weight:      weights.Recency,
			Explanation: calc.explainRecency(factors.RecencyScore),
			Tier:        enums.RECOMMENDATION_TIER_BASIC,
		},
		{
			Name:        "Market Consensus",
			Score:       factors.ConsensusScore,
			Weight:      weights.Consensus,
			Explanation: calc.explainConsensus(factors.ConsensusScore),
			Tier:        enums.RECOMMENDATION_TIER_BASIC,
		},
	}
}

// Explanation methods for each factor
func (calc *BasicScoringCalculator) explainBrokerFrequency(score float64) string {
	switch {
	case score >= 0.8:
		return "High-frequency brokers with strong market presence"
	case score >= 0.6:
		return "Reputable brokers with regular market activity"
	case score >= 0.4:
		return "Moderate broker credibility based on activity"
	default:
		return "Limited broker activity data available"
	}
}

func (calc *BasicScoringCalculator) explainTargetMovement(score float64) string {
	switch {
	case score >= 0.7:
		return "Strong upward price target revisions"
	case score >= 0.5:
		return "Neutral to positive target price changes"
	case score >= 0.3:
		return "Mixed target price movements"
	default:
		return "Downward target price revisions"
	}
}

func (calc *BasicScoringCalculator) explainRatingChange(score float64) string {
	switch {
	case score >= 0.7:
		return "Positive rating upgrades and improvements"
	case score >= 0.5:
		return "Stable ratings with minor changes"
	case score >= 0.3:
		return "Mixed rating changes"
	default:
		return "Rating downgrades and concerns"
	}
}

func (calc *BasicScoringCalculator) explainRecency(score float64) string {
	switch {
	case score >= 0.8:
		return "Very recent market activity and analysis"
	case score >= 0.6:
		return "Recent market developments"
	case score >= 0.4:
		return "Moderately recent activity"
	default:
		return "Older market data with limited freshness"
	}
}

func (calc *BasicScoringCalculator) explainConsensus(score float64) string {
	switch {
	case score >= 0.7:
		return "Strong positive market consensus"
	case score >= 0.5:
		return "Balanced market sentiment"
	case score >= 0.3:
		return "Mixed market opinions"
	default:
		return "Negative market consensus"
	}
}
