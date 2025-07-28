package usecase

import (
	"context"

	recommendationModel "stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/domain/shared/enums"
	stockModel "stock-tracker/internal/domain/stocks/model"
)

// RecommendationUseCase defines the interface for recommendation business logic
type RecommendationUseCase interface {
	GetRecommendations(ctx context.Context, request RecommendationRequest) (*RecommendationResponse, error)
	GetRecommendationByTicker(ctx context.Context, ticker string, userTier enums.UserTier) (*recommendationModel.AggregatedRecommendation, error)
}

// BasicScoringCalculatorInterface defines the interface for basic scoring calculator
type BasicScoringCalculatorInterface interface {
	CalculateAggregatedRecommendation(ctx context.Context, ticker string) (*recommendationModel.AggregatedRecommendation, error)
	CalculateAggregatedRecommendationFromEvents(ctx context.Context, ticker string, events []*stockModel.Stock) (*recommendationModel.AggregatedRecommendation, error)
}

// ExternalDataEnricherInterface defines the interface for external data enricher
type ExternalDataEnricherInterface interface {
	EnrichRecommendation(ctx context.Context, baseRecommendation *recommendationModel.AggregatedRecommendation) (*recommendationModel.AggregatedRecommendation, error)
}
