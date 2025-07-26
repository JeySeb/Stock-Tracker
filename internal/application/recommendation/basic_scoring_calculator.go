package recommendation

import (
	"context"
	"fmt"
	"math"
	"time"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
	"stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/pkg/logger"
)

// BasicScoringCalculator calculates recommendation scores without hardcoded values
type BasicScoringCalculator struct {
	stockRepo repositories.StockRepository
	logger    logger.Logger
}

// BasicScoringFactors represents the different factors used in scoring
type BasicScoringFactors struct {
	BrokerFrequencyScore float64 `json:"broker_frequency_score"` // Calculated dynamically
	TargetMovementScore  float64 `json:"target_movement_score"`  // Based on TargetFrom/TargetTo
	RatingChangeScore    float64 `json:"rating_change_score"`    // Based on RatingFrom/RatingTo
	RecencyScore         float64 `json:"recency_score"`          // Based on EventTime
	ConsensusScore       float64 `json:"consensus_score"`        // Aggregation of events by ticker
}

// NewBasicScoringCalculator creates a new instance of the scoring calculator
func NewBasicScoringCalculator(
	stockRepo repositories.StockRepository,
	logger logger.Logger,
) *BasicScoringCalculator {
	return &BasicScoringCalculator{
		stockRepo: stockRepo,
		logger:    logger,
	}
}

// CalculateAggregatedRecommendation calculates a recommendation based on all events for a ticker
func (calc *BasicScoringCalculator) CalculateAggregatedRecommendation(
	ctx context.Context,
	ticker string,
) (*model.AggregatedRecommendation, error) {

	// 1. Get ALL events for this ticker (last 90 days)
	since := time.Now().AddDate(0, 0, -90)
	events, err := calc.stockRepo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for ticker %s: %w", ticker, err)
	}

	// Filter by date
	recentEvents := calc.filterRecentEvents(events, since)
	if len(recentEvents) == 0 {
		return nil, fmt.Errorf("no recent events for ticker %s", ticker)
	}

	// 2. Calculate aggregated statistics
	stats := calc.calculateEventStatistics(recentEvents)

	// 3. Calculate scores based on real data
	brokerFreqScore := calc.calculateBrokerFrequencyScore(ctx, recentEvents)
	targetScore := calc.calculateTargetMovementScore(recentEvents)
	ratingScore := calc.calculateRatingChangeScore(recentEvents)
	recencyScore := calc.calculateRecencyScore(recentEvents)
	consensusScore := calc.calculateConsensusScore(recentEvents)

	// 4. Combine scores with weights
	weights := model.ScoringWeights{
		BrokerFrequency: 0.25,
		TargetMovement:  0.25,
		RatingChange:    0.25,
		Recency:         0.15,
		Consensus:       0.10,
	}

	scoringFactors := BasicScoringFactors{
		BrokerFrequencyScore: brokerFreqScore,
		TargetMovementScore:  targetScore,
		RatingChangeScore:    ratingScore,
		RecencyScore:         recencyScore,
		ConsensusScore:       consensusScore,
	}

	finalScore := calc.combineScores(scoringFactors, weights)

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
	recommendation.ScoringFactors = calc.createScoringFactorsDetails(scoringFactors, weights)

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

// ⚡ DYNAMIC CALCULATION - NO HARDCODING
// calculateBrokerFrequencyScore calculates broker credibility dynamically from database
func (calc *BasicScoringCalculator) calculateBrokerFrequencyScore(
	ctx context.Context,
	events []*stockModel.Stock,
) float64 {
	// Get broker statistics from database
	brokerStats, err := calc.stockRepo.GetBrokerageStats(ctx)
	if err != nil {
		calc.logger.Warn("Failed to get brokerage stats", "error", err)
		return 0.5 // neutral score if no data available
	}

	// Create frequency map
	brokerFrequency := make(map[string]float64)
	totalReports := 0

	for _, stat := range brokerStats {
		brokerFrequency[stat.Brokerage] = float64(stat.Count)
		totalReports += stat.Count
	}

	if totalReports == 0 {
		return 0.5 // neutral if no historical data
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
		normalizedFreq := math.Log(1+freq) / math.Log(1+float64(totalReports))
		weight := normalizedFreq

		weightedScore += normalizedFreq * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0.5
	}

	return math.Min(1.0, weightedScore/totalWeight)
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
		return 0.5 // neutral if no target data
	}

	avgMovement := totalMovement / float64(count)

	// Normalize: negative changes get lower scores, positive get higher
	// Cap at ±50% change for normalization
	normalizedChange := math.Max(-0.5, math.Min(0.5, avgMovement))
	score := 0.5 + normalizedChange // Maps -50% to 0, 0% to 0.5, +50% to 1.0

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
		return 0.5 // neutral if no rating data
	}

	avgRatingChange := totalRatingChange / float64(count)

	// Normalize rating change to 0-1 scale
	// Positive changes get higher scores
	return math.Max(0, math.Min(1, 0.5+avgRatingChange))
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

	// More recent events get higher scores
	switch {
	case daysSince <= 1:
		return 1.0 // Very recent
	case daysSince <= 7:
		return 0.8 // Within a week
	case daysSince <= 30:
		return 0.6 - (daysSince-7)*0.01 // Linear decay
	default:
		return 0.2 // Old news
	}
}

// calculateConsensusScore calculates consensus based on positive vs negative events
func (calc *BasicScoringCalculator) calculateConsensusScore(events []*stockModel.Stock) float64 {
	if len(events) == 0 {
		return 0.5
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

// calculateConfidence calculates confidence based on the number of events
func (calc *BasicScoringCalculator) calculateConfidence(eventCount int) float64 {
	switch {
	case eventCount >= 10:
		return 0.9 // High confidence
	case eventCount >= 5:
		return 0.7 // Medium confidence
	case eventCount >= 2:
		return 0.5 // Low confidence
	default:
		return 0.3 // Very low confidence
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
