package handlers

import (
	"net/http"
	"strconv"

	marketDataUsecases "stock-tracker/internal/domain/marketdata/usecase"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// MarketDataHandler handles market data analysis HTTP requests
type MarketDataHandler struct {
	marketDataUC marketDataUsecases.MarketDataAnalysisUseCase
	logger       logger.Logger
}

// NewMarketDataHandler creates a new market data handler
func NewMarketDataHandler(marketDataUC marketDataUsecases.MarketDataAnalysisUseCase, logger logger.Logger) *MarketDataHandler {
	return &MarketDataHandler{
		marketDataUC: marketDataUC,
		logger:       logger,
	}
}

// bindPeriodAndLimit extracts and validates period and limit from request
func (h *MarketDataHandler) bindPeriodAndLimit(r *http.Request) (period string, limit int) {
	limitStr := r.URL.Query().Get("limit")
	limit = 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Debug("Invalid limit parameter", "value", limitStr)
		} else {
			limit = h.marketDataUC.ValidateLimit(parsedLimit)
		}
	}

	period = r.URL.Query().Get("period")
	if period == "" {
		period = h.marketDataUC.GetDefaultPeriod()
	}

	return period, limit
}

// GetMarketDataAnalysis retrieves comprehensive analysis for a ticker
func (h *MarketDataHandler) GetMarketDataAnalysis(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	analysis, err := h.marketDataUC.GetMarketDataAnalysis(r.Context(), ticker)
	if err != nil {
		h.logger.Error("Failed to get market data analysis", "ticker", ticker, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve market data analysis"})
		return
	}

	if analysis == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "No analysis found for ticker"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(analysis, nil, "Market data analysis retrieved successfully")
	render.JSON(w, r, response)
}

// GetMarketDataTrend retrieves trend analysis for a ticker
func (h *MarketDataHandler) GetMarketDataTrend(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	period, _ := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	trend, err := h.marketDataUC.GetMarketDataTrend(r.Context(), ticker, period)
	if err != nil {
		h.logger.Error("Failed to get market data trend", "ticker", ticker, "period", period, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve market data trend"})
		return
	}

	if trend == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "No trend data found for ticker"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(trend, nil, "Market data trend retrieved successfully")
	render.JSON(w, r, response)
}

// GetMarketDataSummary retrieves summary statistics
func (h *MarketDataHandler) GetMarketDataSummary(w http.ResponseWriter, r *http.Request) {
	period, _ := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	summary, err := h.marketDataUC.GetMarketDataSummary(r.Context(), period)
	if err != nil {
		h.logger.Error("Failed to get market data summary", "period", period, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve market data summary"})
		return
	}

	if summary == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "No summary data found"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(summary, nil, "Market data summary retrieved successfully")
	render.JSON(w, r, response)
}

// GetTopPerformers retrieves top performing tickers
func (h *MarketDataHandler) GetTopPerformers(w http.ResponseWriter, r *http.Request) {
	period, limit := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	performers, err := h.marketDataUC.GetTopPerformers(r.Context(), limit, period)
	if err != nil {
		h.logger.Error("Failed to get top performers", "limit", limit, "period", period, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve top performers"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(performers, nil, "Top performers retrieved successfully")
	render.JSON(w, r, response)
}

// GetWorstPerformers retrieves worst performing tickers
func (h *MarketDataHandler) GetWorstPerformers(w http.ResponseWriter, r *http.Request) {
	period, limit := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	performers, err := h.marketDataUC.GetWorstPerformers(r.Context(), limit, period)
	if err != nil {
		h.logger.Error("Failed to get worst performers", "limit", limit, "period", period, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve worst performers"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(performers, nil, "Worst performers retrieved successfully")
	render.JSON(w, r, response)
}

// GetMarketDataByTicker retrieves market data for a ticker with filters
func (h *MarketDataHandler) GetMarketDataByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	// Parse query parameters into filters
	params := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	filters := h.marketDataUC.ParseFiltersFromRequest(params)
	marketData, err := h.marketDataUC.GetMarketDataByTicker(r.Context(), ticker, filters)
	if err != nil {
		h.logger.Error("Failed to get market data by ticker", "ticker", ticker, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve market data"})
		return
	}

	if marketData == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "No market data found for ticker"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(marketData, nil, "Market data retrieved successfully")
	render.JSON(w, r, response)
}

// GetLatestMarketData retrieves the latest market data for a ticker
func (h *MarketDataHandler) GetLatestMarketData(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	marketData, err := h.marketDataUC.GetLatestMarketData(r.Context(), ticker)
	if err != nil {
		h.logger.Error("Failed to get latest market data", "ticker", ticker, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve latest market data"})
		return
	}

	if marketData == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "No market data found for ticker"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(marketData, nil, "Latest market data retrieved successfully")
	render.JSON(w, r, response)
}

// GetMarketDataComparison compares two tickers
func (h *MarketDataHandler) GetMarketDataComparison(w http.ResponseWriter, r *http.Request) {
	ticker1 := r.URL.Query().Get("ticker1")
	ticker2 := r.URL.Query().Get("ticker2")

	if ticker1 == "" || ticker2 == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Both ticker1 and ticker2 are required"})
		return
	}

	// TODO: Implement comparison feature - TICKET-123
	w.Header().Set("Allow", "GET")
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Market data comparison not implemented yet"})
}

// GetMostVolatile retrieves most volatile tickers
func (h *MarketDataHandler) GetMostVolatile(w http.ResponseWriter, r *http.Request) {
	period, limit := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	volatile, err := h.marketDataUC.GetMostVolatile(r.Context(), limit, period)
	if err != nil {
		h.logger.Error("Failed to get most volatile tickers", "limit", limit, "period", period, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve most volatile tickers"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(volatile, nil, "Most volatile tickers retrieved successfully")
	render.JSON(w, r, response)
}

// GetMostActive retrieves most active tickers
func (h *MarketDataHandler) GetMostActive(w http.ResponseWriter, r *http.Request) {
	period, limit := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	active, err := h.marketDataUC.GetMostActive(r.Context(), limit, period)
	if err != nil {
		h.logger.Error("Failed to get most active tickers", "limit", limit, "period", period, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve most active tickers"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(active, nil, "Most active tickers retrieved successfully")
	render.JSON(w, r, response)
}

// GetHighRiskTickers retrieves high risk tickers
func (h *MarketDataHandler) GetHighRiskTickers(w http.ResponseWriter, r *http.Request) {
	_, limit := h.bindPeriodAndLimit(r)

	highRisk, err := h.marketDataUC.GetHighRiskTickers(r.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get high risk tickers", "limit", limit, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve high risk tickers"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(highRisk, nil, "High risk tickers retrieved successfully")
	render.JSON(w, r, response)
}

// GetLowRiskTickers retrieves low risk tickers
func (h *MarketDataHandler) GetLowRiskTickers(w http.ResponseWriter, r *http.Request) {
	_, limit := h.bindPeriodAndLimit(r)

	lowRisk, err := h.marketDataUC.GetLowRiskTickers(r.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get low risk tickers", "limit", limit, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve low risk tickers"})
		return
	}

	response := h.marketDataUC.CreateMarketDataResponse(lowRisk, nil, "Low risk tickers retrieved successfully")
	render.JSON(w, r, response)
}

// GetMarketDataWithStockAnalysis retrieves market data combined with stock analysis
func (h *MarketDataHandler) GetMarketDataWithStockAnalysis(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	// TODO: Implement combined analysis - TICKET-124
	w.Header().Set("Allow", "GET")
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Market data with stock analysis not implemented yet"})
}

// GetCorrelationWithBrokerActions retrieves correlation between market data and broker actions
func (h *MarketDataHandler) GetCorrelationWithBrokerActions(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	period, _ := h.bindPeriodAndLimit(r)

	if !h.marketDataUC.ValidatePeriod(period) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid period. Valid periods: 1d, 1w, 1m, 3m"})
		return
	}

	// TODO: Implement correlation analysis - TICKET-125
	w.Header().Set("Allow", "GET")
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Correlation with broker actions not implemented yet"})
}

// GetMarketDataImpactOnRecommendations retrieves impact of market data on recommendations
func (h *MarketDataHandler) GetMarketDataImpactOnRecommendations(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	// TODO: Implement impact analysis - TICKET-126
	w.Header().Set("Allow", "GET")
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Market data impact on recommendations not implemented yet"})
}
