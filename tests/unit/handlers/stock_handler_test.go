package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/valueObjects"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/tests/mocks"
)

type mockStockUseCase struct {
	mock.Mock
}

func (m *mockStockUseCase) GetStocks(ctx context.Context, filters valueObjects.StockFilters) (interface{}, *valueObjects.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0), args.Get(1).(*valueObjects.Pagination), args.Error(2)
}

func (m *mockStockUseCase) GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0), args.Error(1)
}

func (m *mockStockUseCase) GetStats(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func TestStockHandler_GetStocks_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	testStocks := []entities.Stock{
		{
			Ticker:  "AAPL",
			Company: "Apple Inc.",
			Action:  "upgraded by",
		},
	}

	pagination := &valueObjects.Pagination{
		Page:       1,
		Limit:      10,
		TotalItems: 1,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	mockUseCase.On("GetStocks", mock.Anything, mock.AnythingOfType("valueObjects.StockFilters")).
		Return(testStocks, pagination, nil)

	req := httptest.NewRequest("GET", "/stocks?ticker=AAPL&limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStocks(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.StockResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.Limit)

	mockUseCase.AssertExpectations(t)
}

func TestStockHandler_GetStocks_WithFilters(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	mockUseCase.On("GetStocks", mock.Anything, mock.MatchedBy(func(filters valueObjects.StockFilters) bool {
		return filters.Ticker == "AAPL" &&
			filters.Company == "Apple" &&
			filters.Brokerage == "Goldman" &&
			filters.Action == "upgraded" &&
			filters.SortBy == "event_time" &&
			filters.SortOrder == "desc" &&
			filters.Limit == 20 &&
			filters.Offset == 10
	})).Return([]entities.Stock{}, &valueObjects.Pagination{}, nil)

	req := httptest.NewRequest("GET", "/stocks?ticker=AAPL&company=Apple&brokerage=Goldman&action=upgraded&sort_by=event_time&sort_order=desc&limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStocks(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestStockHandler_GetStocks_DefaultFilters(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	mockUseCase.On("GetStocks", mock.Anything, mock.MatchedBy(func(filters valueObjects.StockFilters) bool {
		// Check that defaults are applied
		return filters.Limit == 50 && // Default limit
			filters.SortBy == "event_time" && // Default sort
			filters.SortOrder == "desc" // Default order
	})).Return([]entities.Stock{}, &valueObjects.Pagination{}, nil)

	req := httptest.NewRequest("GET", "/stocks", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStocks(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestStockHandler_GetStocks_UseCaseError(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	expectedError := errors.New("database connection failed")
	mockUseCase.On("GetStocks", mock.Anything, mock.AnythingOfType("valueObjects.StockFilters")).
		Return(nil, (*valueObjects.Pagination)(nil), expectedError)

	mockLogger.On("Error", "Failed to get stocks", "error", mock.Anything).Return()

	req := httptest.NewRequest("GET", "/stocks", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStocks(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse map[string]string
	err := json.NewDecoder(w.Body).Decode(&errorResponse)
	require.NoError(t, err)

	assert.Equal(t, "Failed to retrieve stocks", errorResponse["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestStockHandler_GetStockByTicker_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	testStocks := []entities.Stock{
		{
			Ticker:  "AAPL",
			Company: "Apple Inc.",
			Action:  "upgraded by",
		},
	}

	mockUseCase.On("GetStocksByTicker", mock.Anything, "AAPL").
		Return(testStocks, nil)

	// Create router to test URL parameters
	r := chi.NewRouter()
	r.Get("/stocks/{ticker}", handler.GetStockByTicker)

	req := httptest.NewRequest("GET", "/stocks/AAPL", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.StockResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response.Data)

	mockUseCase.AssertExpectations(t)
}

func TestStockHandler_GetStockByTicker_EmptyTicker(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	req := httptest.NewRequest("GET", "/stocks/", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStockByTicker(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]string
	err := json.NewDecoder(w.Body).Decode(&errorResponse)
	require.NoError(t, err)

	assert.Equal(t, "Ticker is required", errorResponse["error"])
}

func TestStockHandler_GetStockByTicker_UseCaseError(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	expectedError := errors.New("ticker not found")
	mockUseCase.On("GetStocksByTicker", mock.Anything, "INVALID").
		Return(nil, expectedError)

	mockLogger.On("Error", "Failed to get stocks by ticker", "ticker", "INVALID", "error", mock.Anything).Return()

	// Create router to test URL parameters
	r := chi.NewRouter()
	r.Get("/stocks/{ticker}", handler.GetStockByTicker)

	req := httptest.NewRequest("GET", "/stocks/INVALID", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse map[string]string
	err := json.NewDecoder(w.Body).Decode(&errorResponse)
	require.NoError(t, err)

	assert.Equal(t, "Failed to retrieve stocks", errorResponse["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestStockHandler_GetStats_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	testStats := map[string]interface{}{
		"total_stocks":   1000,
		"unique_tickers": 250,
		"total_brokers":  15,
		"last_updated":   "2024-01-15T10:30:00Z",
	}

	mockUseCase.On("GetStats", mock.Anything).
		Return(testStats, nil)

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStats(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.StockResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response.Data)

	statsData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1000), statsData["total_stocks"]) // JSON unmarshals numbers as float64

	mockUseCase.AssertExpectations(t)
}

func TestStockHandler_GetStats_UseCaseError(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	expectedError := errors.New("stats calculation failed")
	mockUseCase.On("GetStats", mock.Anything).
		Return(nil, expectedError)

	mockLogger.On("Error", "Failed to get stats", "error", mock.Anything).Return()

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStats(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse map[string]string
	err := json.NewDecoder(w.Body).Decode(&errorResponse)
	require.NoError(t, err)

	assert.Equal(t, "Failed to retrieve statistics", errorResponse["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestStockHandler_ParseFilters(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected valueObjects.StockFilters
	}{
		{
			name:  "All parameters",
			query: "ticker=AAPL&company=Apple&brokerage=Goldman&action=upgraded&sort_by=ticker&sort_order=asc&limit=25&offset=50",
			expected: valueObjects.StockFilters{
				Ticker:    "AAPL",
				Company:   "Apple",
				Brokerage: "Goldman",
				Action:    "upgraded",
				SortBy:    "ticker",
				SortOrder: "asc",
				Limit:     25,
				Offset:    50,
			},
		},
		{
			name:  "Invalid limit and offset",
			query: "limit=invalid&offset=invalid",
			expected: valueObjects.StockFilters{
				Limit:     50,           // Default
				Offset:    0,            // Default (invalid value ignored)
				SortBy:    "event_time", // Default
				SortOrder: "desc",       // Default
			},
		},
		{
			name:  "Empty query",
			query: "",
			expected: valueObjects.StockFilters{
				Limit:     50,           // Default
				Offset:    0,            // Default
				SortBy:    "event_time", // Default
				SortOrder: "desc",       // Default
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockUseCase := &mockStockUseCase{}
			mockLogger := &mocks.MockLogger{}
			handler := handlers.NewStockHandler(mockUseCase, mockLogger)

			mockUseCase.On("GetStocks", mock.Anything, mock.MatchedBy(func(filters valueObjects.StockFilters) bool {
				return filters.Ticker == tc.expected.Ticker &&
					filters.Company == tc.expected.Company &&
					filters.Brokerage == tc.expected.Brokerage &&
					filters.Action == tc.expected.Action &&
					filters.SortBy == tc.expected.SortBy &&
					filters.SortOrder == tc.expected.SortOrder &&
					filters.Limit == tc.expected.Limit &&
					filters.Offset == tc.expected.Offset
			})).Return([]entities.Stock{}, &valueObjects.Pagination{}, nil)

			req := httptest.NewRequest("GET", "/stocks?"+tc.query, nil)
			w := httptest.NewRecorder()

			// Act
			handler.GetStocks(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestStockHandler_NotImplementedMethods(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	testCases := []struct {
		name    string
		method  string
		path    string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:    "GetStockByID",
			method:  "GET",
			path:    "/stocks/123",
			handler: handler.GetStockByID,
		},
		{
			name:    "CreateStock",
			method:  "POST",
			path:    "/stocks",
			handler: handler.CreateStock,
		},
		{
			name:    "UpdateStock",
			method:  "PUT",
			path:    "/stocks/123",
			handler: handler.UpdateStock,
		},
		{
			name:    "DeleteStock",
			method:  "DELETE",
			path:    "/stocks/123",
			handler: handler.DeleteStock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			// Act
			tc.handler(w, req)

			// Assert
			assert.Equal(t, http.StatusNotImplemented, w.Code)

			var errorResponse map[string]string
			err := json.NewDecoder(w.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.Equal(t, "Not implemented", errorResponse["error"])
		})
	}
}

// Integration test with full router
func TestStockHandler_Integration(t *testing.T) {
	// Arrange
	mockUseCase := &mockStockUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewStockHandler(mockUseCase, mockLogger)

	r := chi.NewRouter()
	r.Route("/api/v1/stocks", func(r chi.Router) {
		r.Get("/", handler.GetStocks)
		r.Get("/{ticker}", handler.GetStockByTicker)
		r.Get("/stats", handler.GetStats)
	})

	testStocks := []entities.Stock{
		{Ticker: "AAPL", Company: "Apple Inc."},
	}

	mockUseCase.On("GetStocks", mock.Anything, mock.AnythingOfType("valueObjects.StockFilters")).
		Return(testStocks, &valueObjects.Pagination{}, nil)

	req := httptest.NewRequest("GET", "/api/v1/stocks?ticker=AAPL", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response handlers.StockResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response.Data)

	mockUseCase.AssertExpectations(t)
}
