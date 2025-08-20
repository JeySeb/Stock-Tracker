package handlers_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	stockValidation "stock-tracker/internal/domain/stocks/validation"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/tests/mocks"
)

// Extend mockStockUseCase (declared in stock_handler_test.go) with additional methods

func (m *mockStockUseCase) GetUniqueTickers(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func (m *mockStockUseCase) GetUniqueCompanies(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

// Helper pointer builders
func intPtr(i int) *int              { return &i }
func float64Ptr(f float64) *float64  { return &f }
func boolPtr(b bool) *bool           { return &b }
func timePtr(t time.Time) *time.Time { return &t }

// Enhanced filter parsing via handler (ensures unexported parser behavior through the public endpoint)
func TestEnhancedFilters_ParseVariousScenarios(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected stockValidation.EnhancedStockFilters
	}{
		{
			name:  "Multiple tickers",
			query: "tickers[]=AAPL&tickers[]=GOOGL&tickers[]=MSFT",
			expected: stockValidation.EnhancedStockFilters{
				Tickers:   []string{"AAPL", "GOOGL", "MSFT"},
				Limit:     50,
				SortBy:    "event_time",
				SortOrder: "desc",
			},
		},
		{
			name:  "Multiple companies",
			query: "companies[]=Apple&companies[]=Google",
			expected: stockValidation.EnhancedStockFilters{
				Companies: []string{"Apple", "Google"},
				Limit:     50,
				SortBy:    "event_time",
				SortOrder: "desc",
			},
		},
		{
			name:  "Multiple brokerages",
			query: "brokerages[]=Goldman&brokerages[]=Morgan",
			expected: stockValidation.EnhancedStockFilters{
				Brokerages: []string{"Goldman", "Morgan"},
				Limit:      50,
				SortBy:     "event_time",
				SortOrder:  "desc",
			},
		},
		{
			name:  "Multiple actions",
			query: "actions[]=upgraded&actions[]=initiated",
			expected: stockValidation.EnhancedStockFilters{
				Actions:   []string{"upgraded", "initiated"},
				Limit:     50,
				SortBy:    "event_time",
				SortOrder: "desc",
			},
		},
		{
			name:  "Time-based filters",
			query: "last_hours=24&last_days=7&last_weeks=2&last_months=1",
			expected: stockValidation.EnhancedStockFilters{
				LastHours:  intPtr(24),
				LastDays:   intPtr(7),
				LastWeeks:  intPtr(2),
				LastMonths: intPtr(1),
				Limit:      50,
				SortBy:     "event_time",
				SortOrder:  "desc",
			},
		},
		{
			name:  "Target price filters",
			query: "target_from=100.0&target_to=200.0",
			expected: stockValidation.EnhancedStockFilters{
				TargetFrom: float64Ptr(100.0),
				TargetTo:   float64Ptr(200.0),
				Limit:      50,
				SortBy:     "event_time",
				SortOrder:  "desc",
			},
		},
		{
			name:  "Date filters",
			query: "date_from=2024-01-01T00:00:00Z&date_to=2024-12-31T23:59:59Z",
			expected: stockValidation.EnhancedStockFilters{
				DateFrom:  timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				DateTo:    timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
				Limit:     50,
				SortBy:    "event_time",
				SortOrder: "desc",
			},
		},
		{
			name:  "Date ranges",
			query: "date_ranges=2024-01-01T00:00:00Z,2024-01-31T23:59:59Z|2024-03-01T00:00:00Z,2024-03-31T23:59:59Z",
			expected: stockValidation.EnhancedStockFilters{
				DateRanges: []stockValidation.DateRange{
					{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), To: time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)},
					{From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), To: time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC)},
				},
				Limit:     50,
				SortBy:    "event_time",
				SortOrder: "desc",
			},
		},
		{
			name:  "Pagination and sorting",
			query: "limit=100&offset=50&sort_by=ticker&sort_order=asc",
			expected: stockValidation.EnhancedStockFilters{
				Limit:     100,
				Offset:    50,
				SortBy:    "ticker",
				SortOrder: "asc",
			},
		},
		{
			name:  "Complex filter combination",
			query: "tickers[]=AAPL&tickers[]=GOOGL&actions[]=upgraded&brokerages[]=Goldman&last_days=30&target_from=100.0&limit=50&sort_by=event_time&sort_order=desc",
			expected: stockValidation.EnhancedStockFilters{
				Tickers:    []string{"AAPL", "GOOGL"},
				Actions:    []string{"upgraded"},
				Brokerages: []string{"Goldman"},
				LastDays:   intPtr(30),
				TargetFrom: float64Ptr(100.0),
				Limit:      50,
				SortBy:     "event_time",
				SortOrder:  "desc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockStockUseCase{}
			mockLogger := &mocks.MockLogger{}
			h := handlers.NewStockHandler(mockUC, mockLogger)

			// Expectation: the parsed filters must match what we expect
			mockUC.On("GetStocksWithEnhancedFilters", mock.Anything, mock.MatchedBy(func(f stockValidation.EnhancedStockFilters) bool {
				// Compare slices
				if len(tt.expected.Tickers) != len(f.Tickers) {
					return false
				}
				for i := range tt.expected.Tickers {
					if tt.expected.Tickers[i] != f.Tickers[i] {
						return false
					}
				}
				if len(tt.expected.Companies) != len(f.Companies) {
					return false
				}
				for i := range tt.expected.Companies {
					if tt.expected.Companies[i] != f.Companies[i] {
						return false
					}
				}
				if len(tt.expected.Brokerages) != len(f.Brokerages) {
					return false
				}
				for i := range tt.expected.Brokerages {
					if tt.expected.Brokerages[i] != f.Brokerages[i] {
						return false
					}
				}
				if len(tt.expected.Actions) != len(f.Actions) {
					return false
				}
				for i := range tt.expected.Actions {
					if tt.expected.Actions[i] != f.Actions[i] {
						return false
					}
				}

				// Compare pointers safely
				if (tt.expected.LastHours == nil) != (f.LastHours == nil) {
					return false
				}
				if tt.expected.LastHours != nil && *tt.expected.LastHours != *f.LastHours {
					return false
				}
				if (tt.expected.LastDays == nil) != (f.LastDays == nil) {
					return false
				}
				if tt.expected.LastDays != nil && *tt.expected.LastDays != *f.LastDays {
					return false
				}
				if (tt.expected.LastWeeks == nil) != (f.LastWeeks == nil) {
					return false
				}
				if tt.expected.LastWeeks != nil && *tt.expected.LastWeeks != *f.LastWeeks {
					return false
				}
				if (tt.expected.LastMonths == nil) != (f.LastMonths == nil) {
					return false
				}
				if tt.expected.LastMonths != nil && *tt.expected.LastMonths != *f.LastMonths {
					return false
				}

				if (tt.expected.TargetFrom == nil) != (f.TargetFrom == nil) {
					return false
				}
				if tt.expected.TargetFrom != nil && *tt.expected.TargetFrom != *f.TargetFrom {
					return false
				}
				if (tt.expected.TargetTo == nil) != (f.TargetTo == nil) {
					return false
				}
				if tt.expected.TargetTo != nil && *tt.expected.TargetTo != *f.TargetTo {
					return false
				}

				if (tt.expected.DateFrom == nil) != (f.DateFrom == nil) {
					return false
				}
				if tt.expected.DateFrom != nil && !tt.expected.DateFrom.Equal(*f.DateFrom) {
					return false
				}
				if (tt.expected.DateTo == nil) != (f.DateTo == nil) {
					return false
				}
				if tt.expected.DateTo != nil && !tt.expected.DateTo.Equal(*f.DateTo) {
					return false
				}

				if len(tt.expected.DateRanges) != len(f.DateRanges) {
					return false
				}
				for i := range tt.expected.DateRanges {
					if !tt.expected.DateRanges[i].From.Equal(f.DateRanges[i].From) {
						return false
					}
					if !tt.expected.DateRanges[i].To.Equal(f.DateRanges[i].To) {
						return false
					}
				}

				// Scalars
				if tt.expected.Limit != f.Limit {
					return false
				}
				if tt.expected.Offset != f.Offset {
					return false
				}
				if tt.expected.SortBy != f.SortBy {
					return false
				}
				if tt.expected.SortOrder != f.SortOrder {
					return false
				}
				return true
			})).Return([]interface{}{}, &stockValidation.Pagination{Page: 1, Limit: 50}, nil)

			req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?"+tt.query, nil)
			w := httptest.NewRecorder()
			h.GetStocksWithEnhancedFilters(w, req)

			assert.Equal(t, 200, w.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

// Advanced enhanced filters parsing
func TestEnhancedFilters_AdvancedFeatures(t *testing.T) {
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
			mockUC := &mockStockUseCase{}
			mockLogger := &mocks.MockLogger{}
			h := handlers.NewStockHandler(mockUC, mockLogger)

			mockUC.On("GetStocksWithEnhancedFilters", mock.Anything, mock.MatchedBy(func(f stockValidation.EnhancedStockFilters) bool {
				// Only compare fields used in each test
				if len(tt.expected.Tickers) != len(f.Tickers) {
					return false
				}
				for i := range tt.expected.Tickers {
					if tt.expected.Tickers[i] != f.Tickers[i] {
						return false
					}
				}
				if len(tt.expected.Companies) != len(f.Companies) {
					return false
				}
				for i := range tt.expected.Companies {
					if tt.expected.Companies[i] != f.Companies[i] {
						return false
					}
				}
				if len(tt.expected.Brokerages) != len(f.Brokerages) {
					return false
				}
				for i := range tt.expected.Brokerages {
					if tt.expected.Brokerages[i] != f.Brokerages[i] {
						return false
					}
				}
				if len(tt.expected.Actions) != len(f.Actions) {
					return false
				}
				for i := range tt.expected.Actions {
					if tt.expected.Actions[i] != f.Actions[i] {
						return false
					}
				}

				if (tt.expected.MinTargetChange == nil) != (f.MinTargetChange == nil) {
					return false
				}
				if tt.expected.MinTargetChange != nil && *tt.expected.MinTargetChange != *f.MinTargetChange {
					return false
				}
				if (tt.expected.MaxTargetChange == nil) != (f.MaxTargetChange == nil) {
					return false
				}
				if tt.expected.MaxTargetChange != nil && *tt.expected.MaxTargetChange != *f.MaxTargetChange {
					return false
				}
				if (tt.expected.HasTargetPrice == nil) != (f.HasTargetPrice == nil) {
					return false
				}
				if tt.expected.HasTargetPrice != nil && *tt.expected.HasTargetPrice != *f.HasTargetPrice {
					return false
				}
				if (tt.expected.HasRating == nil) != (f.HasRating == nil) {
					return false
				}
				if tt.expected.HasRating != nil && *tt.expected.HasRating != *f.HasRating {
					return false
				}
				if (tt.expected.MinBrokerScore == nil) != (f.MinBrokerScore == nil) {
					return false
				}
				if tt.expected.MinBrokerScore != nil && *tt.expected.MinBrokerScore != *f.MinBrokerScore {
					return false
				}
				if (tt.expected.MaxBrokerScore == nil) != (f.MaxBrokerScore == nil) {
					return false
				}
				if tt.expected.MaxBrokerScore != nil && *tt.expected.MaxBrokerScore != *f.MaxBrokerScore {
					return false
				}
				if (tt.expected.LastDays == nil) != (f.LastDays == nil) {
					return false
				}
				if tt.expected.LastDays != nil && *tt.expected.LastDays != *f.LastDays {
					return false
				}
				if (tt.expected.LastWeeks == nil) != (f.LastWeeks == nil) {
					return false
				}
				if tt.expected.LastWeeks != nil && *tt.expected.LastWeeks != *f.LastWeeks {
					return false
				}
				if tt.expected.Limit != f.Limit {
					return false
				}
				if tt.expected.SortBy != f.SortBy {
					return false
				}
				if tt.expected.SortOrder != f.SortOrder {
					return false
				}
				return true
			})).Return([]interface{}{}, &stockValidation.Pagination{Page: 1, Limit: 50}, nil)

			req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?"+tt.query, nil)
			w := httptest.NewRecorder()
			h.GetStocksWithEnhancedFilters(w, req)

			assert.Equal(t, 200, w.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

// Enhanced filters validation tests (direct struct validation)
func TestEnhancedFiltersValidation(t *testing.T) {
	tests := []struct {
		name        string
		filters     stockValidation.EnhancedStockFilters
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid filters",
			filters: stockValidation.EnhancedStockFilters{
				Tickers: []string{"AAPL", "GOOGL"},
				Limit:   50,
			},
			expectError: false,
		},
		{
			name: "Invalid date range",
			filters: stockValidation.EnhancedStockFilters{
				DateFrom: timePtr(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
				DateTo:   timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				Limit:    50,
			},
			expectError: true,
			errorMsg:    "date_from must be before date_to",
		},
		{
			name: "Invalid last_hours",
			filters: stockValidation.EnhancedStockFilters{
				LastHours: intPtr(0),
				Limit:     50,
			},
			expectError: true,
			errorMsg:    "last_hours must be greater than 0",
		},
		{
			name: "Invalid target price range",
			filters: stockValidation.EnhancedStockFilters{
				TargetFrom: float64Ptr(200.0),
				TargetTo:   float64Ptr(100.0),
				Limit:      50,
			},
			expectError: true,
			errorMsg:    "target_from must be less than or equal to target_to",
		},
		{
			name: "Invalid limit",
			filters: stockValidation.EnhancedStockFilters{
				Limit: 0,
			},
			expectError: true,
			errorMsg:    "limit must be greater than 0",
		},
		{
			name: "Invalid offset",
			filters: stockValidation.EnhancedStockFilters{
				Offset: -1,
				Limit:  50,
			},
			expectError: true,
			errorMsg:    "offset must be greater than or equal to 0",
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

// Handler integration tests for enhanced endpoint
func TestEnhancedFilters_HandlerIntegration(t *testing.T) {
	t.Run("Successful enhanced filters request", func(t *testing.T) {
		mockUC := &mockStockUseCase{}
		mockLogger := &mocks.MockLogger{}
		h := handlers.NewStockHandler(mockUC, mockLogger)

		expectedPagination := &stockValidation.Pagination{
			Page:       1,
			Limit:      50,
			TotalItems: 100,
			TotalPages: 2,
			HasNext:    true,
			HasPrev:    false,
		}

		mockUC.On("GetStocksWithEnhancedFilters", mock.Anything, mock.AnythingOfType("validation.EnhancedStockFilters")).
			Return([]interface{}{}, expectedPagination, nil)

		req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?tickers[]=AAPL&tickers[]=GOOGL&min_target_change=10.0&has_target_price=true", nil)
		w := httptest.NewRecorder()
		h.GetStocksWithEnhancedFilters(w, req)

		assert.Equal(t, 200, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		mockUC := &mockStockUseCase{}
		mockLogger := &mocks.MockLogger{}
		h := handlers.NewStockHandler(mockUC, mockLogger)

		// Expect error logging for invalid filters
		mockLogger.On("Error", "Invalid enhanced filters", "error", mock.Anything).Return()

		req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?min_target_change=50.0&max_target_change=10.0", nil)
		w := httptest.NewRecorder()
		h.GetStocksWithEnhancedFilters(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid filters")
	})
}
