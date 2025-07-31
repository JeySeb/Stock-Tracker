package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
	"stock-tracker/internal/domain/marketdata/repositories"
	"stock-tracker/pkg/logger"
)

// MarketDataAnalysisUseCase defines the interface for market data analysis business logic
type MarketDataAnalysisUseCase interface {
	// Basic market data operations
	GetMarketDataAnalysis(ctx context.Context, ticker string) (*model.MarketDataAnalysis, error)
	GetMarketDataTrend(ctx context.Context, ticker string, period string) (*model.MarketDataTrend, error)
	GetMarketDataSummary(ctx context.Context, period string) (*model.MarketDataSummary, error)
	GetMarketDataByTicker(ctx context.Context, ticker string, filters *model.MarketDataFilters) ([]*model.MarketData, error)
	GetLatestMarketData(ctx context.Context, ticker string) (*model.MarketData, error)

	// Advanced analytics
	GetTopPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetWorstPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetMostVolatile(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetMostActive(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)

	// Risk analysis
	GetHighRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error)
	GetLowRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error)

	// Utility methods
	ParseFiltersFromRequest(params map[string]string) *model.MarketDataFilters
	ValidatePeriod(period string) bool
	GetDefaultPeriod() string
	ValidateLimit(limit int) int
	CreateMarketDataResponse(data interface{}, pagination *model.Pagination, message string) *model.MarketDataResponse
	CalculatePagination(total, page, perPage int) *model.Pagination
}

// marketDataAnalysisUseCaseImpl handles market data analysis operations
type marketDataAnalysisUseCaseImpl struct {
	marketDataRepo repositories.MarketDataAnalysisRepository
	logger         logger.Logger
}

// NewMarketDataAnalysisUseCase creates a new market data analysis use case
func NewMarketDataAnalysisUseCase(
	marketDataRepo repositories.MarketDataAnalysisRepository,
	logger logger.Logger,
) MarketDataAnalysisUseCase {
	return &marketDataAnalysisUseCaseImpl{
		marketDataRepo: marketDataRepo,
		logger:         logger,
	}
}

// GetMarketDataAnalysis retrieves comprehensive analysis for a ticker
func (uc *marketDataAnalysisUseCaseImpl) GetMarketDataAnalysis(ctx context.Context, ticker string) (*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting market data analysis", "ticker", ticker)

	analysis, err := uc.marketDataRepo.GetMarketDataAnalysis(ctx, ticker)
	if err != nil {
		uc.logger.Error("Failed to get market data analysis", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to get market data analysis: %w", err)
	}

	return analysis, nil
}

// GetMarketDataTrend retrieves trend analysis for a ticker
func (uc *marketDataAnalysisUseCaseImpl) GetMarketDataTrend(ctx context.Context, ticker string, period string) (*model.MarketDataTrend, error) {
	uc.logger.Info("Getting market data trend", "ticker", ticker, "period", period)

	trend, err := uc.marketDataRepo.GetMarketDataTrend(ctx, ticker, period)
	if err != nil {
		uc.logger.Error("Failed to get market data trend", "ticker", ticker, "period", period, "error", err)
		return nil, fmt.Errorf("failed to get market data trend: %w", err)
	}

	return trend, nil
}

// GetMarketDataSummary retrieves summary statistics
func (uc *marketDataAnalysisUseCaseImpl) GetMarketDataSummary(ctx context.Context, period string) (*model.MarketDataSummary, error) {
	uc.logger.Info("Getting market data summary", "period", period)

	summary, err := uc.marketDataRepo.GetMarketDataSummary(ctx, period)
	if err != nil {
		uc.logger.Error("Failed to get market data summary", "period", period, "error", err)
		return nil, fmt.Errorf("failed to get market data summary: %w", err)
	}

	return summary, nil
}

// GetTopPerformers retrieves top performing tickers
func (uc *marketDataAnalysisUseCaseImpl) GetTopPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting top performers", "limit", limit, "period", period)

	performers, err := uc.marketDataRepo.GetTopPerformers(ctx, limit, period)
	if err != nil {
		uc.logger.Error("Failed to get top performers", "limit", limit, "period", period, "error", err)
		return nil, fmt.Errorf("failed to get top performers: %w", err)
	}

	return performers, nil
}

// GetWorstPerformers retrieves worst performing tickers
func (uc *marketDataAnalysisUseCaseImpl) GetWorstPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting worst performers", "limit", limit, "period", period)

	performers, err := uc.marketDataRepo.GetWorstPerformers(ctx, limit, period)
	if err != nil {
		uc.logger.Error("Failed to get worst performers", "limit", limit, "period", period, "error", err)
		return nil, fmt.Errorf("failed to get worst performers: %w", err)
	}

	return performers, nil
}

// GetMarketDataByTicker retrieves market data for a ticker with filters
func (uc *marketDataAnalysisUseCaseImpl) GetMarketDataByTicker(ctx context.Context, ticker string, filters *model.MarketDataFilters) ([]*model.MarketData, error) {
	uc.logger.Info("Getting market data by ticker", "ticker", ticker)

	filters.SetDefaults()
	marketData, err := uc.marketDataRepo.GetMarketDataByTicker(ctx, ticker, filters)
	if err != nil {
		uc.logger.Error("Failed to get market data by ticker", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to get market data by ticker: %w", err)
	}

	return marketData, nil
}

// GetLatestMarketData retrieves the latest market data for a ticker
func (uc *marketDataAnalysisUseCaseImpl) GetLatestMarketData(ctx context.Context, ticker string) (*model.MarketData, error) {
	uc.logger.Info("Getting latest market data", "ticker", ticker)

	marketData, err := uc.marketDataRepo.GetLatestMarketDataByTicker(ctx, ticker)
	if err != nil {
		uc.logger.Error("Failed to get latest market data", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	return marketData, nil
}

// ParseFiltersFromRequest parses filters from HTTP request parameters
func (uc *marketDataAnalysisUseCaseImpl) ParseFiltersFromRequest(params map[string]string) *model.MarketDataFilters {
	filters := &model.MarketDataFilters{}

	// Parse basic filters
	if ticker := params["ticker"]; ticker != "" {
		filters.Ticker = ticker
	}
	if dataSource := params["data_source"]; dataSource != "" {
		filters.DataSource = dataSource
	}
	if dataQuality := params["data_quality"]; dataQuality != "" {
		filters.DataQuality = dataQuality
	}
	if trendDirection := params["trend_direction"]; trendDirection != "" {
		filters.TrendDirection = trendDirection
	}
	if riskLevel := params["risk_level"]; riskLevel != "" {
		filters.RiskLevel = riskLevel
	}

	// Parse numeric filters
	if minPriceStr := params["min_price"]; minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}
	if maxPriceStr := params["max_price"]; maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}
	if minChangeStr := params["min_change"]; minChangeStr != "" {
		if minChange, err := strconv.ParseFloat(minChangeStr, 64); err == nil {
			filters.MinChange = &minChange
		}
	}
	if maxChangeStr := params["max_change"]; maxChangeStr != "" {
		if maxChange, err := strconv.ParseFloat(maxChangeStr, 64); err == nil {
			filters.MaxChange = &maxChange
		}
	}
	if minVolumeStr := params["min_volume"]; minVolumeStr != "" {
		if minVolume, err := strconv.ParseInt(minVolumeStr, 10, 64); err == nil {
			filters.MinVolume = &minVolume
		}
	}
	if maxVolumeStr := params["max_volume"]; maxVolumeStr != "" {
		if maxVolume, err := strconv.ParseInt(maxVolumeStr, 10, 64); err == nil {
			filters.MaxVolume = &maxVolume
		}
	}

	// Parse pagination
	if limitStr := params["limit"]; limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}
	if offsetStr := params["offset"]; offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// Parse sorting
	if sortBy := params["sort_by"]; sortBy != "" {
		filters.SortBy = sortBy
	}
	if sortOrder := params["sort_order"]; sortOrder != "" {
		filters.SortOrder = sortOrder
	}

	// Parse date filters
	if startDateStr := params["start_date"]; startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = startDate
		}
	}
	if endDateStr := params["end_date"]; endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = endDate
		}
	}

	filters.SetDefaults()
	return filters
}

// ValidatePeriod validates the analysis period
func (uc *marketDataAnalysisUseCaseImpl) ValidatePeriod(period string) bool {
	validPeriods := []string{"1d", "1w", "1m", "3m"}
	for _, validPeriod := range validPeriods {
		if period == validPeriod {
			return true
		}
	}
	return false
}

// GetDefaultPeriod returns the default period if none is specified
func (uc *marketDataAnalysisUseCaseImpl) GetDefaultPeriod() string {
	return "1d"
}

// ValidateLimit validates and returns a safe limit value
func (uc *marketDataAnalysisUseCaseImpl) ValidateLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// CreateMarketDataResponse creates a standardized API response
func (uc *marketDataAnalysisUseCaseImpl) CreateMarketDataResponse(data interface{}, pagination *model.Pagination, message string) *model.MarketDataResponse {
	response := &model.MarketDataResponse{
		Data:    data,
		Message: message,
	}

	if pagination != nil {
		response.Pagination = pagination
	}

	// Add metadata
	response.Metadata = &model.Metadata{
		GeneratedAt: time.Now(),
		DataPoints:  0, // Will be calculated based on data type
	}

	return response
}

// CalculatePagination calculates pagination information
func (uc *marketDataAnalysisUseCaseImpl) CalculatePagination(total, page, perPage int) *model.Pagination {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 50
	}

	totalPages := (total + perPage - 1) / perPage
	hasNext := page < totalPages
	hasPrevious := page > 1

	return &model.Pagination{
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}
}

// GetMostVolatile retrieves the most volatile tickers
func (uc *marketDataAnalysisUseCaseImpl) GetMostVolatile(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting most volatile tickers", "limit", limit, "period", period)

	limit = uc.ValidateLimit(limit)
	volatile, err := uc.marketDataRepo.GetMostVolatile(ctx, limit, period)
	if err != nil {
		uc.logger.Error("Failed to get most volatile tickers", "limit", limit, "period", period, "error", err)
		return nil, fmt.Errorf("failed to get most volatile tickers: %w", err)
	}

	return volatile, nil
}

// GetMostActive retrieves the most active tickers by volume
func (uc *marketDataAnalysisUseCaseImpl) GetMostActive(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting most active tickers", "limit", limit, "period", period)

	limit = uc.ValidateLimit(limit)
	active, err := uc.marketDataRepo.GetMostActive(ctx, limit, period)
	if err != nil {
		uc.logger.Error("Failed to get most active tickers", "limit", limit, "period", period, "error", err)
		return nil, fmt.Errorf("failed to get most active tickers: %w", err)
	}

	return active, nil
}

// GetHighRiskTickers retrieves high risk tickers
func (uc *marketDataAnalysisUseCaseImpl) GetHighRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting high risk tickers", "limit", limit)

	limit = uc.ValidateLimit(limit)
	highRisk, err := uc.marketDataRepo.GetHighRiskTickers(ctx, limit)
	if err != nil {
		uc.logger.Error("Failed to get high risk tickers", "limit", limit, "error", err)
		return nil, fmt.Errorf("failed to get high risk tickers: %w", err)
	}

	return highRisk, nil
}

// GetLowRiskTickers retrieves low risk tickers
func (uc *marketDataAnalysisUseCaseImpl) GetLowRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error) {
	uc.logger.Info("Getting low risk tickers", "limit", limit)

	limit = uc.ValidateLimit(limit)
	lowRisk, err := uc.marketDataRepo.GetLowRiskTickers(ctx, limit)
	if err != nil {
		uc.logger.Error("Failed to get low risk tickers", "limit", limit, "error", err)
		return nil, fmt.Errorf("failed to get low risk tickers: %w", err)
	}

	return lowRisk, nil
}
