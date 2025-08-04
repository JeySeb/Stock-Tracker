package usecases

import (
	"context"
	"fmt"
	"math"
	"strings"

	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/pkg/logger"
)

type BrokerUseCase interface {
	GetBrokersWithScores(ctx context.Context, limit *int, orderBy string) ([]*stockRepos.BrokerWithScore, error)
}

type BrokerQueryUseCase struct {
	stockRepo  stockRepos.StockRepository
	brokerRepo stockRepos.BrokerRepository
	logger     logger.Logger
}

func NewBrokerQueryUseCase(
	stockRepo stockRepos.StockRepository,
	brokerRepo stockRepos.BrokerRepository,
	logger logger.Logger,
) BrokerUseCase {
	return &BrokerQueryUseCase{
		stockRepo:  stockRepo,
		brokerRepo: brokerRepo,
		logger:     logger,
	}
}

// GetBrokersWithScores returns brokers with their calculated scores
func (uc *BrokerQueryUseCase) GetBrokersWithScores(
	ctx context.Context,
	limit *int,
	orderBy string,
) ([]*stockRepos.BrokerWithScore, error) {
	uc.logger.Info("Getting brokers with scores", "limit", limit, "orderBy", orderBy)

	// Validate orderBy parameter
	orderBy = uc.normalizeOrderBy(orderBy)

	// Get brokers with scores from repository
	brokers, err := uc.brokerRepo.GetBrokersWithScores(ctx, limit, orderBy)
	if err != nil {
		uc.logger.Error("Failed to get brokers with scores", "error", err)
		return nil, fmt.Errorf("failed to retrieve brokers with scores: %w", err)
	}

	// Calculate additional scores if needed
	for _, broker := range brokers {
		broker.CalculatedScore = uc.calculateBrokerScore(broker)
	}

	uc.logger.Info("Successfully retrieved brokers with scores", "count", len(brokers))
	return brokers, nil
}

// normalizeOrderBy validates and normalizes the orderBy parameter
func (uc *BrokerQueryUseCase) normalizeOrderBy(orderBy string) string {
	orderBy = strings.ToLower(strings.TrimSpace(orderBy))

	// Default to score descending if invalid
	if orderBy != "asc" && orderBy != "desc" {
		return "desc"
	}

	return orderBy
}

// calculateBrokerScore calculates a comprehensive score for a broker
// This uses the same logic as the BasicScoringCalculator but for individual brokers
func (uc *BrokerQueryUseCase) calculateBrokerScore(broker *stockRepos.BrokerWithScore) float64 {
	if broker.ReportCount == 0 {
		return 0.0
	}

	// Use logarithmic scaling to avoid outliers, similar to BasicScoringCalculator
	// Normalize the report count using log scaling
	logScore := math.Log(1 + float64(broker.ReportCount))

	// Combine with credibility score (weighted average)
	credibilityWeight := 0.0 // CURRENTLY NOT USED
	activityWeight := 1.0

	activityScore := math.Min(1.0, logScore/10.0) // Normalize to 0-1 range

	calculatedScore := (broker.CredibilityScore * credibilityWeight) + (activityScore * activityWeight)

	return math.Min(1.0, math.Max(0.0, calculatedScore))
}
