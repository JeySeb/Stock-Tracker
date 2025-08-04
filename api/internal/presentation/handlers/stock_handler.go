package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	stockUsecases "stock-tracker/internal/domain/stocks/usecase"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type StockHandler struct {
	stockUC stockUsecases.StockUseCase
	logger  logger.Logger
}

type StockResponse struct {
	Data       interface{}                 `json:"data"`
	Pagination *stockValidation.Pagination `json:"pagination,omitempty"`
	Message    string                      `json:"message,omitempty"`
}

func NewStockHandler(stockUC stockUsecases.StockUseCase, logger logger.Logger) *StockHandler {
	return &StockHandler{
		stockUC: stockUC,
		logger:  logger,
	}
}

func (h *StockHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	filters := h.parseFilters(r)

	stocks, pagination, err := h.stockUC.GetStocks(r.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to get stocks", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve stocks"})
		return
	}

	response := StockResponse{
		Data:       stocks,
		Pagination: pagination,
	}

	render.JSON(w, r, response)
}

func (h *StockHandler) GetStockByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticker is required"})
		return
	}

	stocks, err := h.stockUC.GetStocksByTicker(r.Context(), ticker)
	if err != nil {
		h.logger.Error("Failed to get stocks by ticker", "ticker", ticker, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve stocks"})
		return
	}

	response := StockResponse{
		Data: stocks,
	}

	render.JSON(w, r, response)
}

func (h *StockHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.stockUC.GetStats(r.Context())
	if err != nil {
		h.logger.Error("Failed to get stats", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve statistics"})
		return
	}

	response := StockResponse{
		Data: stats,
	}

	render.JSON(w, r, response)
}

func (h *StockHandler) GetStocksWithEnhancedFilters(w http.ResponseWriter, r *http.Request) {
	filters := h.parseEnhancedFilters(r)

	// Validate filters
	if err := filters.Validate(); err != nil {
		h.logger.Error("Invalid enhanced filters", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid filters: " + err.Error()})
		return
	}

	stocks, pagination, err := h.stockUC.GetStocksWithEnhancedFilters(r.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to get stocks with enhanced filters", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve stocks"})
		return
	}

	response := StockResponse{
		Data:       stocks,
		Pagination: pagination,
	}

	render.JSON(w, r, response)
}

func (h *StockHandler) parseFilters(r *http.Request) stockValidation.StockFilters {
	filters := stockValidation.StockFilters{
		Ticker:    r.URL.Query().Get("ticker"),
		Company:   r.URL.Query().Get("company"),
		Brokerage: r.URL.Query().Get("brokerage"),
		Action:    r.URL.Query().Get("action"),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	filters.SetDefaults()
	return filters
}

func (h *StockHandler) parseEnhancedFilters(r *http.Request) stockValidation.EnhancedStockFilters {
	filters := stockValidation.EnhancedStockFilters{
		RatingFrom: r.URL.Query().Get("rating_from"),
		RatingTo:   r.URL.Query().Get("rating_to"),
		SortBy:     r.URL.Query().Get("sort_by"),
		SortOrder:  r.URL.Query().Get("sort_order"),
	}

	// Parse array parameters - handle both formats: tickers[]=value and tickers=value
	if tickers := r.URL.Query()["tickers[]"]; len(tickers) > 0 {
		filters.Tickers = tickers
	} else if tickers := r.URL.Query()["tickers"]; len(tickers) > 0 {
		filters.Tickers = tickers
	}
	if companies := r.URL.Query()["companies[]"]; len(companies) > 0 {
		filters.Companies = companies
	} else if companies := r.URL.Query()["companies"]; len(companies) > 0 {
		filters.Companies = companies
	}
	if brokerages := r.URL.Query()["brokerages[]"]; len(brokerages) > 0 {
		filters.Brokerages = brokerages
	} else if brokerages := r.URL.Query()["brokerages"]; len(brokerages) > 0 {
		filters.Brokerages = brokerages
	}
	if actions := r.URL.Query()["actions[]"]; len(actions) > 0 {
		filters.Actions = actions
	} else if actions := r.URL.Query()["actions"]; len(actions) > 0 {
		filters.Actions = actions
	}

	// Parse numeric parameters
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// Parse time-based filters
	if lastHoursStr := r.URL.Query().Get("last_hours"); lastHoursStr != "" {
		if lastHours, err := strconv.Atoi(lastHoursStr); err == nil {
			filters.LastHours = &lastHours
		}
	}

	if lastDaysStr := r.URL.Query().Get("last_days"); lastDaysStr != "" {
		if lastDays, err := strconv.Atoi(lastDaysStr); err == nil {
			filters.LastDays = &lastDays
		}
	}

	if lastWeeksStr := r.URL.Query().Get("last_weeks"); lastWeeksStr != "" {
		if lastWeeks, err := strconv.Atoi(lastWeeksStr); err == nil {
			filters.LastWeeks = &lastWeeks
		}
	}

	if lastMonthsStr := r.URL.Query().Get("last_months"); lastMonthsStr != "" {
		if lastMonths, err := strconv.Atoi(lastMonthsStr); err == nil {
			filters.LastMonths = &lastMonths
		}
	}

	// Parse target price filters
	if targetFromStr := r.URL.Query().Get("target_from"); targetFromStr != "" {
		if targetFrom, err := strconv.ParseFloat(targetFromStr, 64); err == nil {
			filters.TargetFrom = &targetFrom
		}
	}

	if targetToStr := r.URL.Query().Get("target_to"); targetToStr != "" {
		if targetTo, err := strconv.ParseFloat(targetToStr, 64); err == nil {
			filters.TargetTo = &targetTo
		}
	}

	// Parse advanced filters
	if minTargetChangeStr := r.URL.Query().Get("min_target_change"); minTargetChangeStr != "" {
		if minTargetChange, err := strconv.ParseFloat(minTargetChangeStr, 64); err == nil {
			filters.MinTargetChange = &minTargetChange
		}
	}

	if maxTargetChangeStr := r.URL.Query().Get("max_target_change"); maxTargetChangeStr != "" {
		if maxTargetChange, err := strconv.ParseFloat(maxTargetChangeStr, 64); err == nil {
			filters.MaxTargetChange = &maxTargetChange
		}
	}

	if hasTargetPriceStr := r.URL.Query().Get("has_target_price"); hasTargetPriceStr != "" {
		if hasTargetPrice, err := strconv.ParseBool(hasTargetPriceStr); err == nil {
			filters.HasTargetPrice = &hasTargetPrice
		}
	}

	if hasRatingStr := r.URL.Query().Get("has_rating"); hasRatingStr != "" {
		if hasRating, err := strconv.ParseBool(hasRatingStr); err == nil {
			filters.HasRating = &hasRating
		}
	}

	// Parse brokerage score filters
	if minBrokerScoreStr := r.URL.Query().Get("min_broker_score"); minBrokerScoreStr != "" {
		if minBrokerScore, err := strconv.ParseFloat(minBrokerScoreStr, 64); err == nil {
			filters.MinBrokerScore = &minBrokerScore
		}
	}

	if maxBrokerScoreStr := r.URL.Query().Get("max_broker_score"); maxBrokerScoreStr != "" {
		if maxBrokerScore, err := strconv.ParseFloat(maxBrokerScoreStr, 64); err == nil {
			filters.MaxBrokerScore = &maxBrokerScore
		}
	}

	// Parse date filters
	if dateFromStr := r.URL.Query().Get("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			filters.DateFrom = &dateFrom
		}
	}

	if dateToStr := r.URL.Query().Get("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			filters.DateTo = &dateTo
		}
	}

	// Parse date ranges (comma-separated pairs)
	if dateRangesStr := r.URL.Query().Get("date_ranges"); dateRangesStr != "" {
		dateRanges := strings.Split(dateRangesStr, "|")
		for _, dateRangeStr := range dateRanges {
			parts := strings.Split(dateRangeStr, ",")
			if len(parts) == 2 {
				if from, err := time.Parse(time.RFC3339, parts[0]); err == nil {
					if to, err := time.Parse(time.RFC3339, parts[1]); err == nil {
						filters.DateRanges = append(filters.DateRanges, stockValidation.DateRange{
							From: from,
							To:   to,
						})
					}
				}
			}
		}
	}

	filters.SetDefaults()
	return filters
}

// GetStockByID retrieves a stock by its ID
func (h *StockHandler) GetStockByID(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Not implemented"})
}

// CreateStock creates a new stock
func (h *StockHandler) CreateStock(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Not implemented"})
}

// UpdateStock updates an existing stock
func (h *StockHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Not implemented"})
}

// DeleteStock deletes a stock
func (h *StockHandler) DeleteStock(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Not implemented"})
}
