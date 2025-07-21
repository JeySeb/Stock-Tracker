package entities_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"stock-tracker/internal/domain/entities"
)

func TestNewStock(t *testing.T) {
	ticker := "AAPL"
	company := "Apple Inc."
	brokerage := "Goldman Sachs"
	action := "upgraded by"
	eventTime := time.Now()

	stock := entities.NewStock(ticker, company, brokerage, action, eventTime)

	assert.NotNil(t, stock)
	assert.NotEmpty(t, stock.ID)
	assert.Equal(t, ticker, stock.Ticker)
	assert.Equal(t, company, stock.Company)
	assert.Equal(t, brokerage, stock.Brokerage)
	assert.Equal(t, action, stock.Action)
	assert.Equal(t, eventTime, stock.EventTime)
	assert.NotZero(t, stock.CreatedAt)
	assert.NotZero(t, stock.UpdatedAt)
}

func TestStock_IsUpgrade(t *testing.T) {
	testCases := []struct {
		name     string
		action   string
		expected bool
	}{
		{
			name:     "Upgraded action",
			action:   "upgraded by Goldman Sachs",
			expected: true,
		},
		{
			name:     "Raised action",
			action:   "raised to Buy",
			expected: true,
		},
		{
			name:     "Initiated action",
			action:   "initiated by Morgan Stanley",
			expected: true,
		},
		{
			name:     "Downgraded action",
			action:   "downgraded by JP Morgan",
			expected: false,
		},
		{
			name:     "Lowered action",
			action:   "target lowered by Wells Fargo",
			expected: false,
		},
		{
			name:     "Reiterated action",
			action:   "reiterated by Deutsche Bank",
			expected: false,
		},
		{
			name:     "Case insensitive test",
			action:   "UPGRADED BY Citigroup",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &entities.Stock{Action: tc.action}
			result := stock.IsUpgrade()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStock_GetPriceTargetChange(t *testing.T) {
	testCases := []struct {
		name       string
		targetFrom float64
		targetTo   float64
		expected   float64
	}{
		{
			name:       "Positive change",
			targetFrom: 100.0,
			targetTo:   120.0,
			expected:   0.2, // 20% increase
		},
		{
			name:       "Negative change",
			targetFrom: 100.0,
			targetTo:   80.0,
			expected:   -0.2, // 20% decrease
		},
		{
			name:       "No change",
			targetFrom: 100.0,
			targetTo:   100.0,
			expected:   0.0,
		},
		{
			name:       "Zero target from",
			targetFrom: 0.0,
			targetTo:   100.0,
			expected:   0.0,
		},
		{
			name:       "Zero target to",
			targetFrom: 100.0,
			targetTo:   0.0,
			expected:   0.0,
		},
		{
			name:       "Both zero",
			targetFrom: 0.0,
			targetTo:   0.0,
			expected:   0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &entities.Stock{
				TargetFrom: tc.targetFrom,
				TargetTo:   tc.targetTo,
			}
			result := stock.GetPriceTargetChange()
			assert.InDelta(t, tc.expected, result, 0.001)
		})
	}
}

func TestStock_GetRatingScore(t *testing.T) {
	testCases := []struct {
		name         string
		ratingFrom   string
		ratingTo     string
		expectedFrom float64
		expectedTo   float64
	}{
		{
			name:         "Strong Buy to Buy",
			ratingFrom:   "Strong Buy",
			ratingTo:     "Buy",
			expectedFrom: 1.0,
			expectedTo:   0.8,
		},
		{
			name:         "Hold to Outperform",
			ratingFrom:   "Hold",
			ratingTo:     "Outperform",
			expectedFrom: 0.5,
			expectedTo:   0.75,
		},
		{
			name:         "Sell to Neutral",
			ratingFrom:   "Sell",
			ratingTo:     "Neutral",
			expectedFrom: 0.2,
			expectedTo:   0.4,
		},
		{
			name:         "Case insensitive",
			ratingFrom:   "strong buy",
			ratingTo:     "HOLD",
			expectedFrom: 1.0,
			expectedTo:   0.5,
		},
		{
			name:         "Unknown rating",
			ratingFrom:   "Unknown",
			ratingTo:     "Invalid",
			expectedFrom: 0.0,
			expectedTo:   0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &entities.Stock{
				RatingFrom: tc.ratingFrom,
				RatingTo:   tc.ratingTo,
			}
			fromScore, toScore := stock.GetRatingScore()
			assert.Equal(t, tc.expectedFrom, fromScore)
			assert.Equal(t, tc.expectedTo, toScore)
		})
	}
}

func TestStock_GetRatingChangeScore(t *testing.T) {
	testCases := []struct {
		name       string
		ratingFrom string
		ratingTo   string
		expected   float64
	}{
		{
			name:       "Upgrade from Hold to Buy",
			ratingFrom: "Hold",
			ratingTo:   "Buy",
			expected:   0.3, // 0.8 - 0.5
		},
		{
			name:       "Downgrade from Buy to Hold",
			ratingFrom: "Buy",
			ratingTo:   "Hold",
			expected:   -0.3, // 0.5 - 0.8
		},
		{
			name:       "No change",
			ratingFrom: "Buy",
			ratingTo:   "Buy",
			expected:   0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &entities.Stock{
				RatingFrom: tc.ratingFrom,
				RatingTo:   tc.ratingTo,
			}
			result := stock.GetRatingChangeScore()
			assert.InDelta(t, tc.expected, result, 0.001)
		})
	}
}

func TestStock_IsRecommendation(t *testing.T) {
	testCases := []struct {
		name       string
		action     string
		ratingFrom string
		ratingTo   string
		expected   bool
	}{
		{
			name:       "Positive action - upgraded",
			action:     "upgraded by Goldman Sachs",
			ratingFrom: "Hold",
			ratingTo:   "Hold",
			expected:   true,
		},
		{
			name:       "Positive action - initiated",
			action:     "initiated by JP Morgan",
			ratingFrom: "Neutral",
			ratingTo:   "Neutral",
			expected:   true,
		},
		{
			name:       "Positive action - reiterated",
			action:     "reiterated by Wells Fargo",
			ratingFrom: "Buy",
			ratingTo:   "Buy",
			expected:   true,
		},
		{
			name:       "Rating improvement",
			action:     "target lowered by Morgan Stanley",
			ratingFrom: "Hold",
			ratingTo:   "Buy",
			expected:   true,
		},
		{
			name:       "Negative action and rating",
			action:     "downgraded by Citigroup",
			ratingFrom: "Buy",
			ratingTo:   "Hold",
			expected:   false,
		},
		{
			name:       "No positive indicators",
			action:     "target lowered by Deutsche Bank",
			ratingFrom: "Buy",
			ratingTo:   "Sell",
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &entities.Stock{
				Action:     tc.action,
				RatingFrom: tc.ratingFrom,
				RatingTo:   tc.ratingTo,
			}
			result := stock.IsRecommendation()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStock_GetPriceChange(t *testing.T) {
	testCases := []struct {
		name       string
		priceClose *float64
		expected   float64
	}{
		{
			name:       "With price close",
			priceClose: func() *float64 { p := 150.50; return &p }(),
			expected:   150.50,
		},
		{
			name:       "No price close",
			priceClose: nil,
			expected:   0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &entities.Stock{
				PriceClose: tc.priceClose,
			}
			result := stock.GetPriceChange()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStock_Validate(t *testing.T) {
	// Test valid stock
	validStock := entities.NewStock("AAPL", "Apple Inc.", "Goldman Sachs", "upgraded by", time.Now())
	err := validStock.Validate()
	assert.NoError(t, err)

	// Test invalid stock - empty ticker
	invalidStock := &entities.Stock{
		Ticker:  "",
		Company: "Apple Inc.",
		Action:  "upgraded by",
	}
	err = invalidStock.Validate()
	assert.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "ticker")

	// Test invalid stock - empty company
	invalidStock2 := &entities.Stock{
		Ticker:  "AAPL",
		Company: "",
		Action:  "upgraded by",
	}
	err = invalidStock2.Validate()
	assert.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "company")
}

// Benchmark tests for performance critical methods
func BenchmarkStock_IsUpgrade(b *testing.B) {
	stock := &entities.Stock{Action: "upgraded by Goldman Sachs"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stock.IsUpgrade()
	}
}

func BenchmarkStock_GetRatingScore(b *testing.B) {
	stock := &entities.Stock{
		RatingFrom: "Hold",
		RatingTo:   "Buy",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stock.GetRatingScore()
	}
}
