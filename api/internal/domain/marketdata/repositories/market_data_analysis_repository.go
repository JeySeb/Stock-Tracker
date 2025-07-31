package repositories

import (
	"context"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
)

// MarketDataAnalysisRepository defines the interface for market data analysis operations
type MarketDataAnalysisRepository interface {
	// Basic market data operations
	GetMarketDataByTicker(ctx context.Context, ticker string, filters *model.MarketDataFilters) ([]*model.MarketData, error)
	GetLatestMarketDataByTicker(ctx context.Context, ticker string) (*model.MarketData, error)
	GetMarketDataByDateRange(ctx context.Context, startDate, endDate time.Time, filters *model.MarketDataFilters) ([]*model.MarketData, error)
	
	// Analysis operations
	GetMarketDataAnalysis(ctx context.Context, ticker string) (*model.MarketDataAnalysis, error)
	GetMarketDataTrend(ctx context.Context, ticker string, period string) (*model.MarketDataTrend, error)
	GetMarketDataComparison(ctx context.Context, ticker1, ticker2 string, date time.Time) (*model.MarketDataComparison, error)
	GetMarketDataSummary(ctx context.Context, period string) (*model.MarketDataSummary, error)
	
	// Advanced analytics
	GetTopPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetWorstPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetMostVolatile(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetMostActive(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	
	// Correlation and patterns
	GetCorrelatedTickers(ctx context.Context, ticker string, threshold float64, period string) ([]string, error)
	GetPriceBreakouts(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetVolumeSurges(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	
	// Risk analysis
	GetHighRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error)
	GetLowRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error)
	GetRiskDistribution(ctx context.Context) (map[string]int, error)
	
	// Trend analysis
	GetTrendingTickers(ctx context.Context, direction string, limit int, period string) ([]*model.MarketDataAnalysis, error)
	GetTrendStrength(ctx context.Context, ticker string, period string) (string, error)
	GetTrendReversal(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error)
	
	// Statistical operations
	GetPriceStatistics(ctx context.Context, ticker string, period string) (map[string]float64, error)
	GetVolumeStatistics(ctx context.Context, ticker string, period string) (map[string]int64, error)
	GetVolatilityStatistics(ctx context.Context, ticker string, period string) (map[string]float64, error)
	
	// Alert system
	GetMarketAlerts(ctx context.Context, alertType string, severity string, limit int) ([]*model.MarketDataAlert, error)
	CreateMarketAlert(ctx context.Context, alert *model.MarketDataAlert) error
	UpdateMarketAlert(ctx context.Context, alert *model.MarketDataAlert) error
	DeleteMarketAlert(ctx context.Context, alertID string) error
	
	// Data quality and monitoring
	GetDataQualityStats(ctx context.Context) (map[string]interface{}, error)
	GetDataSourceStats(ctx context.Context) (map[string]interface{}, error)
	GetDataFreshness(ctx context.Context) (map[string]time.Time, error)
	
	// Cross-domain operations (market data + stocks)
	GetMarketDataWithStockAnalysis(ctx context.Context, ticker string) (map[string]interface{}, error)
	GetCorrelationWithBrokerActions(ctx context.Context, ticker string, period string) (map[string]interface{}, error)
	GetMarketDataImpactOnRecommendations(ctx context.Context, ticker string) (map[string]interface{}, error)
} 