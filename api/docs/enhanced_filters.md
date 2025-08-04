# Enhanced Stock Filters API

This document describes the enhanced filtering capabilities for the stock tracking API, which provides robust filtering options for retrieving stock data with multiple selections per categorical filter and advanced date/time filtering.

## Endpoint

```
GET /api/v1/stocks/enhanced
```

## Filter Parameters

### Multiple Selection Filters

These filters support multiple values and use array syntax in query parameters:

#### Tickers
- **Parameter**: `tickers[]`
- **Type**: Array of strings
- **Description**: Filter by multiple stock tickers
- **Example**: `?tickers[]=AAPL&tickers[]=GOOGL&tickers[]=MSFT`

#### Companies
- **Parameter**: `companies[]`
- **Type**: Array of strings
- **Description**: Filter by multiple company names (partial match)
- **Example**: `?companies[]=Apple&companies[]=Google`

#### Brokerages
- **Parameter**: `brokerages[]`
- **Type**: Array of strings
- **Description**: Filter by multiple brokerage names (partial match)
- **Example**: `?brokerages[]=Goldman&brokerages[]=Morgan`

#### Actions
- **Parameter**: `actions[]`
- **Type**: Array of strings
- **Description**: Filter by multiple action types (partial match)
- **Example**: `?actions[]=upgraded&actions[]=initiated`

### Rating Filters

#### Rating From
- **Parameter**: `rating_from`
- **Type**: String
- **Description**: Filter by rating from value (partial match)
- **Example**: `?rating_from=Buy`

#### Rating To
- **Parameter**: `rating_to`
- **Type**: String
- **Description**: Filter by rating to value (partial match)
- **Example**: `?rating_to=Strong Buy`

### Enhanced Date Filters

#### Basic Date Range
- **Parameter**: `date_from`
- **Type**: ISO 8601 datetime string
- **Description**: Filter stocks from this date
- **Example**: `?date_from=2024-01-01T00:00:00Z`

- **Parameter**: `date_to`
- **Type**: ISO 8601 datetime string
- **Description**: Filter stocks until this date
- **Example**: `?date_to=2024-12-31T23:59:59Z`

#### Multiple Date Ranges
- **Parameter**: `date_ranges`
- **Type**: Pipe-separated date range pairs
- **Description**: Filter by multiple specific date ranges
- **Format**: `from1,to1|from2,to2|...`
- **Example**: `?date_ranges=2024-01-01T00:00:00Z,2024-01-31T23:59:59Z|2024-03-01T00:00:00Z,2024-03-31T23:59:59Z`

### Time-Based Filters

These filters are relative to the current time:

#### Last Hours
- **Parameter**: `last_hours`
- **Type**: Integer
- **Description**: Filter stocks from the last N hours
- **Example**: `?last_hours=24`

#### Last Days
- **Parameter**: `last_days`
- **Type**: Integer
- **Description**: Filter stocks from the last N days
- **Example**: `?last_days=7`

#### Last Weeks
- **Parameter**: `last_weeks`
- **Type**: Integer
- **Description**: Filter stocks from the last N weeks
- **Example**: `?last_weeks=4`

#### Last Months
- **Parameter**: `last_months`
- **Type**: Integer
- **Description**: Filter stocks from the last N months
- **Example**: `?last_months=3`

### Target Price Filters

#### Target From
- **Parameter**: `target_from`
- **Type**: Float
- **Description**: Filter stocks with target price >= this value
- **Example**: `?target_from=100.0`

#### Target To
- **Parameter**: `target_to`
- **Type**: Float
- **Description**: Filter stocks with target price <= this value
- **Example**: `?target_to=200.0`

### Advanced Filters

#### Target Change Percentage
- **Parameter**: `min_target_change`
- **Type**: Float
- **Description**: Filter stocks with minimum target price change percentage
- **Example**: `?min_target_change=10.5`

- **Parameter**: `max_target_change`
- **Type**: Float
- **Description**: Filter stocks with maximum target price change percentage
- **Example**: `?max_target_change=50.0`

#### Data Availability Filters
- **Parameter**: `has_target_price`
- **Type**: Boolean
- **Description**: Filter stocks that have/don't have target prices
- **Example**: `?has_target_price=true`

- **Parameter**: `has_rating`
- **Type**: Boolean
- **Description**: Filter stocks that have/don't have ratings
- **Example**: `?has_rating=false`

### Brokerage Credibility Filters

#### Broker Score Range
- **Parameter**: `min_broker_score`
- **Type**: Float
- **Description**: Filter stocks from brokerages with minimum credibility score
- **Example**: `?min_broker_score=7.5`

- **Parameter**: `max_broker_score`
- **Type**: Float
- **Description**: Filter stocks from brokerages with maximum credibility score
- **Example**: `?max_broker_score=9.5`

### Pagination and Sorting

#### Limit
- **Parameter**: `limit`
- **Type**: Integer
- **Description**: Number of results per page (max 1000)
- **Default**: 50
- **Example**: `?limit=100`

#### Offset
- **Parameter**: `offset`
- **Type**: Integer
- **Description**: Number of results to skip
- **Default**: 0
- **Example**: `?offset=100`

#### Sort By
- **Parameter**: `sort_by`
- **Type**: String
- **Description**: Field to sort by
- **Default**: `event_time`
- **Options**: `event_time`, `ticker`, `company`, `action`, `target_from`, `target_to`
- **Example**: `?sort_by=ticker`

#### Sort Order
- **Parameter**: `sort_order`
- **Type**: String
- **Description**: Sort order
- **Default**: `desc`
- **Options**: `asc`, `desc`
- **Example**: `?sort_order=asc`

## Usage Examples

### Basic Multiple Ticker Filter
```bash
GET /api/v1/stocks/enhanced?tickers[]=AAPL&tickers[]=GOOGL&tickers[]=MSFT
```

### Complex Filter with Multiple Criteria
```bash
GET /api/v1/stocks/enhanced?tickers[]=AAPL&tickers[]=GOOGL&actions[]=upgraded&actions[]=initiated&brokerages[]=Goldman&last_days=30&target_from=100.0&limit=50&sort_by=event_time&sort_order=desc
```

### Time-Based Filtering
```bash
GET /api/v1/stocks/enhanced?last_hours=24&actions[]=upgraded&limit=100
```

### Multiple Date Ranges
```bash
GET /api/v1/stocks/enhanced?date_ranges=2024-01-01T00:00:00Z,2024-01-31T23:59:59Z|2024-03-01T00:00:00Z,2024-03-31T23:59:59Z&tickers[]=AAPL
```

### Company and Brokerage Filter
```bash
GET /api/v1/stocks/enhanced?companies[]=Apple&companies[]=Google&brokerages[]=Goldman&brokerages[]=Morgan&last_weeks=4
```

### Advanced Filtering with Target Changes
```bash
GET /api/v1/stocks/enhanced?tickers[]=AAPL&tickers[]=GOOGL&min_target_change=15.0&max_target_change=50.0&has_target_price=true&last_days=30
```

### High-Credibility Brokerage Filter
```bash
GET /api/v1/stocks/enhanced?min_broker_score=8.0&actions[]=upgraded&actions[]=initiated&has_rating=true&limit=100
```

### Complex Multi-Criteria Filter
```bash
GET /api/v1/stocks/enhanced?tickers[]=AAPL&tickers[]=MSFT&companies[]=Apple&brokerages[]=Goldman&actions[]=upgraded&min_target_change=10.0&max_broker_score=9.5&has_target_price=true&last_weeks=2&sort_by=event_time&sort_order=desc&limit=50
```

## Response Format

The response follows the same format as the basic filters endpoint:

```json
{
  "data": [
    {
      "id": "uuid",
      "ticker": "AAPL",
      "company": "Apple Inc.",
      "action": "upgraded",
      "rating_from": "Hold",
      "rating_to": "Buy",
      "target_from": 150.0,
      "target_to": 180.0,
      "event_time": "2024-01-15T10:30:00Z",
      "brokerage": "Goldman Sachs",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total_pages": 5,
    "total_items": 250,
    "has_next": true,
    "has_prev": false
  }
}
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful request
- `400 Bad Request`: Invalid filter parameters
- `500 Internal Server Error`: Server error

### Validation Errors

The API validates all filter parameters and returns descriptive error messages:

- Invalid date formats
- Invalid numeric values
- Conflicting date ranges
- Invalid sort fields
- Exceeded limit values

## Performance Considerations

- The enhanced filters are optimized for database performance
- Multiple selection filters use efficient SQL IN clauses
- Date range filters use indexed columns
- Time-based filters are calculated efficiently
- Pagination is handled at the database level

## Backward Compatibility

The original `/api/v1/stocks` endpoint remains unchanged and fully functional. The enhanced filters are available as a separate endpoint to ensure no breaking changes to existing integrations. 