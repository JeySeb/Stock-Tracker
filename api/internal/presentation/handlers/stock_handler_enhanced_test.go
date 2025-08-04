package handlers

import (
	"net/http/httptest"
	"testing"

	stockValidation "stock-tracker/internal/domain/stocks/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a simple mock implementation for testing
type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...interface{})  {}
func (m *MockLogger) Error(msg string, args ...interface{}) {}
func (m *MockLogger) Debug(msg string, args ...interface{}) {}
func (m *MockLogger) Warn(msg string, args ...interface{})  {}

func TestEnhancedFiltersAdvancedFeatures(t *testing.T) {
	handler := &StockHandler{}

	tests := []struct {
		name     string
		query    string
		expected stockValidation.EnhancedStockFilters
	}{
		{
			name:  "Target change percentage filters",
			query: "min_target_change=10.5&max_target_change=50.0",
			expected: stockValidation.EnhancedStockFilters{
				MinTargetChange: float64Ptr(10.5),
				MaxTargetChange: float64Ptr(50.0),
				Limit:           50,
				SortBy:          "event_time",
				SortOrder:       "desc",
			},
		},
		{
			name:  "Data availability filters",
			query: "has_target_price=true&has_rating=false",
			expected: stockValidation.EnhancedStockFilters{
				HasTargetPrice: boolPtr(true),
				HasRating:      boolPtr(false),
				Limit:          50,
				SortBy:         "event_time",
				SortOrder:      "desc",
			},
		},
		{
			name:  "Brokerage credibility filters",
			query: "min_broker_score=7.5&max_broker_score=9.5",
			expected: stockValidation.EnhancedStockFilters{
				MinBrokerScore: float64Ptr(7.5),
				MaxBrokerScore: float64Ptr(9.5),
				Limit:          50,
				SortBy:         "event_time",
				SortOrder:      "desc",
			},
		},
		{
			name:  "Complex advanced filter combination",
			query: "tickers[]=AAPL&min_target_change=15.0&max_target_change=50.0&has_target_price=true&min_broker_score=8.0&last_days=30",
			expected: stockValidation.EnhancedStockFilters{
				Tickers:         []string{"AAPL"},
				MinTargetChange: float64Ptr(15.0),
				MaxTargetChange: float64Ptr(50.0),
				HasTargetPrice:  boolPtr(true),
				MinBrokerScore:  float64Ptr(8.0),
				LastDays:        intPtr(30),
				Limit:           50,
				SortBy:          "event_time",
				SortOrder:       "desc",
			},
		},
		{
			name:  "All advanced filters combined",
			query: "tickers[]=AAPL&tickers[]=GOOGL&companies[]=Apple&brokerages[]=Goldman&actions[]=upgraded&min_target_change=10.0&max_target_change=40.0&has_target_price=true&has_rating=true&min_broker_score=7.0&max_broker_score=9.5&last_weeks=2&sort_by=event_time&sort_order=desc&limit=100",
			expected: stockValidation.EnhancedStockFilters{
				Tickers:         []string{"AAPL", "GOOGL"},
				Companies:       []string{"Apple"},
				Brokerages:      []string{"Goldman"},
				Actions:         []string{"upgraded"},
				MinTargetChange: float64Ptr(10.0),
				MaxTargetChange: float64Ptr(40.0),
				HasTargetPrice:  boolPtr(true),
				HasRating:       boolPtr(true),
				MinBrokerScore:  float64Ptr(7.0),
				MaxBrokerScore:  float64Ptr(9.5),
				LastWeeks:       intPtr(2),
				SortBy:          "event_time",
				SortOrder:       "desc",
				Limit:           100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?"+tt.query, nil)
			result := handler.parseEnhancedFilters(req)

			// Compare basic filters
			assert.Equal(t, tt.expected.Tickers, result.Tickers)
			assert.Equal(t, tt.expected.Companies, result.Companies)
			assert.Equal(t, tt.expected.Brokerages, result.Brokerages)
			assert.Equal(t, tt.expected.Actions, result.Actions)
			assert.Equal(t, tt.expected.Limit, result.Limit)
			assert.Equal(t, tt.expected.SortBy, result.SortBy)
			assert.Equal(t, tt.expected.SortOrder, result.SortOrder)

			// Compare advanced filters
			if tt.expected.MinTargetChange != nil {
				assert.Equal(t, *tt.expected.MinTargetChange, *result.MinTargetChange)
			}
			if tt.expected.MaxTargetChange != nil {
				assert.Equal(t, *tt.expected.MaxTargetChange, *result.MaxTargetChange)
			}
			if tt.expected.HasTargetPrice != nil {
				assert.Equal(t, *tt.expected.HasTargetPrice, *result.HasTargetPrice)
			}
			if tt.expected.HasRating != nil {
				assert.Equal(t, *tt.expected.HasRating, *result.HasRating)
			}
			if tt.expected.MinBrokerScore != nil {
				assert.Equal(t, *tt.expected.MinBrokerScore, *result.MinBrokerScore)
			}
			if tt.expected.MaxBrokerScore != nil {
				assert.Equal(t, *tt.expected.MaxBrokerScore, *result.MaxBrokerScore)
			}

			// Compare time-based filters
			if tt.expected.LastDays != nil {
				assert.Equal(t, *tt.expected.LastDays, *result.LastDays)
			}
			if tt.expected.LastWeeks != nil {
				assert.Equal(t, *tt.expected.LastWeeks, *result.LastWeeks)
			}
		})
	}
}

func TestEnhancedFiltersValidationAdvanced(t *testing.T) {
	tests := []struct {
		name        string
		filters     stockValidation.EnhancedStockFilters
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid advanced filters",
			filters: stockValidation.EnhancedStockFilters{
				Tickers:         []string{"AAPL", "GOOGL"},
				MinTargetChange: float64Ptr(10.0),
				MaxTargetChange: float64Ptr(50.0),
				HasTargetPrice:  boolPtr(true),
				MinBrokerScore:  float64Ptr(7.0),
				MaxBrokerScore:  float64Ptr(9.5),
				Limit:           50,
			},
			expectError: false,
		},
		{
			name: "Invalid target change range",
			filters: stockValidation.EnhancedStockFilters{
				MinTargetChange: float64Ptr(50.0),
				MaxTargetChange: float64Ptr(10.0),
				Limit:           50,
			},
			expectError: true,
			errorMsg:    "min_target_change must be less than or equal to max_target_change",
		},
		{
			name: "Invalid brokerage score range",
			filters: stockValidation.EnhancedStockFilters{
				MinBrokerScore: float64Ptr(9.5),
				MaxBrokerScore: float64Ptr(7.0),
				Limit:          50,
			},
			expectError: true,
			errorMsg:    "min_broker_score must be less than or equal to max_broker_score",
		},
		{
			name: "Valid complex filter combination",
			filters: stockValidation.EnhancedStockFilters{
				Tickers:         []string{"AAPL", "GOOGL", "MSFT"},
				Companies:       []string{"Apple", "Google"},
				Brokerages:      []string{"Goldman", "Morgan"},
				Actions:         []string{"upgraded", "initiated"},
				MinTargetChange: float64Ptr(5.0),
				MaxTargetChange: float64Ptr(30.0),
				HasTargetPrice:  boolPtr(true),
				HasRating:       boolPtr(false),
				MinBrokerScore:  float64Ptr(6.0),
				MaxBrokerScore:  float64Ptr(9.0),
				LastDays:        intPtr(7),
				Limit:           100,
				SortBy:          "event_time",
				SortOrder:       "desc",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filters.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnhancedFiltersHandlerIntegration(t *testing.T) {
	mockUseCase := &MockStockUseCase{}
	mockLogger := &MockLogger{}
	handler := &StockHandler{
		stockUC: mockUseCase,
		logger:  mockLogger,
	}

	// Test successful enhanced filters request
	t.Run("Successful enhanced filters request", func(t *testing.T) {
		expectedPagination := &stockValidation.Pagination{
			Page:       1,
			Limit:      50,
			TotalItems: 100,
			TotalPages: 2,
			HasNext:    true,
			HasPrev:    false,
		}

		mockUseCase.On("GetStocksWithEnhancedFilters", mock.Anything, mock.AnythingOfType("validation.EnhancedStockFilters")).Return([]interface{}{}, expectedPagination, nil)

		req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?tickers[]=AAPL&tickers[]=GOOGL&min_target_change=10.0&has_target_price=true", nil)
		w := httptest.NewRecorder()

		handler.GetStocksWithEnhancedFilters(w, req)

		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	// Test validation error
	t.Run("Validation error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?min_target_change=50.0&max_target_change=10.0", nil)
		w := httptest.NewRecorder()

		handler.GetStocksWithEnhancedFilters(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid filters")
	})
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}
