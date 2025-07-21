package testhelpers

import (
	"time"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/valueObjects"
)

// CreateTestStock creates a stock entity for testing
func CreateTestStock(ticker, company, brokerage, action string) *entities.Stock {
	stock := entities.NewStock(ticker, company, brokerage, action, time.Now())
	stock.RatingFrom = "Hold"
	stock.RatingTo = "Buy"
	stock.TargetFrom = 100.0
	stock.TargetTo = 120.0
	return stock
}

// CreateTestBroker creates a broker entity for testing
func CreateTestBroker(name string, credibility float64) *entities.Broker {
	return entities.NewBroker(name, credibility)
}

// CreateTestStockFilters creates stock filters for testing
func CreateTestStockFilters() valueObjects.StockFilters {
	filters := valueObjects.StockFilters{
		Ticker:    "AAPL",
		Company:   "Apple",
		Brokerage: "Goldman Sachs",
		Action:    "upgraded",
		SortBy:    "event_time",
		SortOrder: "desc",
		Limit:     10,
		Offset:    0,
	}
	filters.SetDefaults()
	return filters
}

// CreateTestPagination creates pagination for testing
func CreateTestPagination() *valueObjects.Pagination {
	return &valueObjects.Pagination{
		Page:       1,
		Limit:      10,
		TotalItems: 100,
		TotalPages: 10,
		HasNext:    true,
		HasPrev:    false,
	}
}

// TestStockData represents common test stock data
var TestStockData = []struct {
	Ticker     string
	Company    string
	Brokerage  string
	Action     string
	RatingFrom string
	RatingTo   string
	TargetFrom float64
	TargetTo   float64
}{
	{
		Ticker:     "AAPL",
		Company:    "Apple Inc.",
		Brokerage:  "Goldman Sachs",
		Action:     "upgraded by",
		RatingFrom: "Hold",
		RatingTo:   "Buy",
		TargetFrom: 150.0,
		TargetTo:   180.0,
	},
	{
		Ticker:     "GOOGL",
		Company:    "Alphabet Inc.",
		Brokerage:  "Morgan Stanley",
		Action:     "initiated by",
		RatingFrom: "Neutral",
		RatingTo:   "Outperform",
		TargetFrom: 2800.0,
		TargetTo:   3000.0,
	},
	{
		Ticker:     "MSFT",
		Company:    "Microsoft Corporation",
		Brokerage:  "JP Morgan",
		Action:     "reiterated by",
		RatingFrom: "Buy",
		RatingTo:   "Buy",
		TargetFrom: 350.0,
		TargetTo:   400.0,
	},
}

// CreateTestStocks creates multiple test stocks
func CreateTestStocks(count int) []*entities.Stock {
	stocks := make([]*entities.Stock, 0, count)

	for i := 0; i < count; i++ {
		dataIndex := i % len(TestStockData)
		data := TestStockData[dataIndex]

		stock := CreateTestStock(
			data.Ticker,
			data.Company,
			data.Brokerage,
			data.Action,
		)
		stock.RatingFrom = data.RatingFrom
		stock.RatingTo = data.RatingTo
		stock.TargetFrom = data.TargetFrom
		stock.TargetTo = data.TargetTo

		stocks = append(stocks, stock)
	}

	return stocks
}

// AssertStockEquals compares two stocks for equality in tests
func AssertStockEquals(expected, actual *entities.Stock) bool {
	return expected.Ticker == actual.Ticker &&
		expected.Company == actual.Company &&
		expected.Brokerage == actual.Brokerage &&
		expected.Action == actual.Action &&
		expected.RatingFrom == actual.RatingFrom &&
		expected.RatingTo == actual.RatingTo
}

// MockTime returns a fixed time for consistent testing
func MockTime() time.Time {
	return time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
}
