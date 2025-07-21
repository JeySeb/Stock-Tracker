package handlers

import (
	"net/http"
	"strconv"

	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/domain/valueObjects"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type StockHandler struct {
	stockUC usecases.StockUseCase
	logger  logger.Logger
}

type StockResponse struct {
	Data       interface{}              `json:"data"`
	Pagination *valueObjects.Pagination `json:"pagination,omitempty"`
	Message    string                   `json:"message,omitempty"`
}

func NewStockHandler(stockUC usecases.StockUseCase, logger logger.Logger) *StockHandler {
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

func (h *StockHandler) parseFilters(r *http.Request) valueObjects.StockFilters {
	filters := valueObjects.StockFilters{
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
