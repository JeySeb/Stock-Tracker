package handlers

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	stockValidation "stock-tracker/internal/domain/stocks/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStockUseCase is a mock implementation for testing
type MockStockUseCase struct {
	mock.Mock
}

func (m *MockStockUseCase) GetStocks(ctx context.Context, filters stockValidation.StockFilters) (interface{}, *stockValidation.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0), args.Get(1).(*stockValidation.Pagination), args.Error(2)
}

func (m *MockStockUseCase) GetStocksWithEnhancedFilters(ctx context.Context, filters stockValidation.EnhancedStockFilters) (interface{}, *stockValidation.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0), args.Get(1).(*stockValidation.Pagination), args.Error(2)
}

func (m *MockStockUseCase) GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0), args.Error(1)
}

func (m *MockStockUseCase) GetStats(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func TestParseEnhancedFilters(t *testing.T) {
	handler := &StockHandler{}

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
					{
						From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						To:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
					},
					{
						From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
					},
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
			req := httptest.NewRequest("GET", "/api/v1/stocks/enhanced?"+tt.query, nil)
			result := handler.parseEnhancedFilters(req)

			// Compare the result with expected
			assert.Equal(t, tt.expected.Tickers, result.Tickers)
			assert.Equal(t, tt.expected.Companies, result.Companies)
			assert.Equal(t, tt.expected.Brokerages, result.Brokerages)
			assert.Equal(t, tt.expected.Actions, result.Actions)
			assert.Equal(t, tt.expected.RatingFrom, result.RatingFrom)
			assert.Equal(t, tt.expected.RatingTo, result.RatingTo)
			assert.Equal(t, tt.expected.Limit, result.Limit)
			assert.Equal(t, tt.expected.Offset, result.Offset)
			assert.Equal(t, tt.expected.SortBy, result.SortBy)
			assert.Equal(t, tt.expected.SortOrder, result.SortOrder)

			// Compare time-based filters
			if tt.expected.LastHours != nil {
				assert.Equal(t, *tt.expected.LastHours, *result.LastHours)
			}
			if tt.expected.LastDays != nil {
				assert.Equal(t, *tt.expected.LastDays, *result.LastDays)
			}
			if tt.expected.LastWeeks != nil {
				assert.Equal(t, *tt.expected.LastWeeks, *result.LastWeeks)
			}
			if tt.expected.LastMonths != nil {
				assert.Equal(t, *tt.expected.LastMonths, *result.LastMonths)
			}

			// Compare target price filters
			if tt.expected.TargetFrom != nil {
				assert.Equal(t, *tt.expected.TargetFrom, *result.TargetFrom)
			}
			if tt.expected.TargetTo != nil {
				assert.Equal(t, *tt.expected.TargetTo, *result.TargetTo)
			}

			// Compare date filters
			if tt.expected.DateFrom != nil {
				assert.Equal(t, tt.expected.DateFrom.Format(time.RFC3339), result.DateFrom.Format(time.RFC3339))
			}
			if tt.expected.DateTo != nil {
				assert.Equal(t, tt.expected.DateTo.Format(time.RFC3339), result.DateTo.Format(time.RFC3339))
			}

			// Compare date ranges
			if len(tt.expected.DateRanges) > 0 {
				assert.Equal(t, len(tt.expected.DateRanges), len(result.DateRanges))
				for i, expectedRange := range tt.expected.DateRanges {
					assert.Equal(t, expectedRange.From.Format(time.RFC3339), result.DateRanges[i].From.Format(time.RFC3339))
					assert.Equal(t, expectedRange.To.Format(time.RFC3339), result.DateRanges[i].To.Format(time.RFC3339))
				}
			}
		})
	}
}

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

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}
