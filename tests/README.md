# 🧪 Comprehensive Test Suite

This directory contains a robust unit test suite for the stock recommendation system backend, designed to ensure high code quality and reliability.

## 📁 Test Structure

```
tests/
├── mocks/                    # Mock implementations for testing
│   ├── repositories.go       # Repository mocks
│   ├── clients.go            # API client mocks
│   └── logger.go             # Logger mocks
├── unit/                     # Unit tests
│   ├── entities/             # Domain entity tests
│   │   └── stock_test.go     # Stock business logic tests
│   ├── usecases/             # Use case tests  
│   │   └── stock_ingestion_test.go  # Ingestion logic tests
│   ├── clients/              # External client tests
│   │   └── stock_api_client_test.go # API client tests
│   └── handlers/             # HTTP handler tests
│       └── stock_handler_test.go    # REST API tests
├── testhelpers/              # Test utilities and helpers
│   └── test_helpers.go       # Common test data creators
└── README.md                 # This file
```

## 🎯 Test Coverage Goals

- **Target Coverage**: ≥80% (industry standard)
- **Current Focus**: Unit tests for all critical business logic
- **Quality Gates**: All tests must pass before deployment

## 🔧 Test Categories

### 1. Entity Tests (`entities/`)
- **Focus**: Domain business logic validation
- **Coverage**: All entity methods and edge cases
- **Examples**:
  - Stock rating calculations
  - Price target analysis
  - Recommendation scoring
  - Input validation

### 2. Use Case Tests (`usecases/`)
- **Focus**: Application service logic
- **Coverage**: Complete workflow testing with mocks
- **Examples**:
  - Stock ingestion orchestration
  - Error handling scenarios
  - Concurrent processing validation
  - Context cancellation

### 3. Client Tests (`clients/`)
- **Focus**: External API integration
- **Coverage**: HTTP client behavior and error handling
- **Examples**:
  - API pagination handling
  - Rate limiting validation
  - Network error scenarios
  - Data transformation

### 4. Handler Tests (`handlers/`)
- **Focus**: HTTP API layer
- **Coverage**: Request/response handling and validation
- **Examples**:
  - Parameter parsing
  - Error response formatting
  - Status code validation
  - Content type handling

## 🚀 Running Tests

### Quick Commands

```bash
# Run all tests with coverage
go test -coverprofile=coverage.out ./tests/unit/...

# Run specific test suite
go test -v ./tests/unit/entities/

# Run single test
go test -v ./tests/unit/entities/ -run TestStock_IsUpgrade

# Run with race detection
go test -race ./tests/unit/...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Comprehensive Test Runner

```bash
# Run complete test suite with reporting
./scripts/run-tests.sh
```

This script provides:
- ✅ Unit test execution
- 📊 Coverage analysis (HTML reports)
- 🏃‍♂️ Benchmark testing
- 🔍 Race condition detection
- 📋 Comprehensive reporting

## 📊 Test Metrics & Quality Gates

### Coverage Targets by Layer
- **Entities**: ≥95% (critical business logic)
- **Use Cases**: ≥90% (application orchestration)
- **Handlers**: ≥85% (API interface)
- **Clients**: ≥80% (external integration)

### Performance Benchmarks
- **Entity Operations**: <1ms per operation
- **API Responses**: <100ms for simple queries
- **Batch Processing**: <10s for 1000 records

## 🛠️ Test Patterns & Best Practices

### 1. Arrange-Act-Assert (AAA) Pattern
```go
func TestStock_IsUpgrade(t *testing.T) {
    // Arrange
    stock := &entities.Stock{Action: "upgraded by Goldman Sachs"}
    
    // Act
    result := stock.IsUpgrade()
    
    // Assert
    assert.True(t, result)
}
```

### 2. Table-Driven Tests
```go
testCases := []struct {
    name     string
    action   string
    expected bool
}{
    {"Upgraded action", "upgraded by Goldman Sachs", true},
    {"Downgraded action", "downgraded by JP Morgan", false},
}
```

### 3. Mock Usage with Testify
```go
mockRepo := &mocks.MockStockRepository{}
mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(stocks, nil)
```

### 4. HTTP Testing with httptest
```go
req := httptest.NewRequest("GET", "/stocks?ticker=AAPL", nil)
w := httptest.NewRecorder()
handler.GetStocks(w, req)
```

## 📈 Continuous Integration Integration

### GitHub Actions Example
```yaml
- name: Run Tests
  run: |
    go test -race -coverprofile=coverage.out ./tests/unit/...
    go tool cover -func=coverage.out
```

### Coverage Requirements
- **Pull Requests**: Must maintain ≥80% coverage
- **New Features**: Must include comprehensive tests
- **Bug Fixes**: Must include regression tests

## 🔍 Test Data Management

### Test Helpers (`testhelpers/`)
- `CreateTestStock()` - Standard test stock creation
- `CreateTestBroker()` - Standard test broker creation
- `TestStockData` - Realistic test data sets
- `MockTime()` - Deterministic time for tests

### Mock Strategy
- **Repository Layer**: Full mocking for isolation
- **External APIs**: HTTP mocking with httptest
- **Logging**: Mock to verify log calls
- **Time**: Fixed time for deterministic tests

## 🚨 Troubleshooting

### Common Issues

1. **Import Errors**
   ```bash
   go mod tidy
   ```

2. **Mock Assertion Failures**
   ```go
   // Always call at end of test
   mockRepo.AssertExpectations(t)
   ```

3. **Race Conditions**
   ```bash
   go test -race ./...
   ```

4. **Coverage Not Generated**
   ```bash
   go test -coverprofile=coverage.out ./tests/unit/...
   ```

## 📝 Writing New Tests

### Entity Tests
1. Test all public methods
2. Include edge cases and error conditions
3. Validate business rules
4. Test input validation

### Use Case Tests  
1. Mock all dependencies
2. Test success and failure paths
3. Verify error handling
4. Test concurrent scenarios

### Handler Tests
1. Test all HTTP methods
2. Validate request parsing
3. Check response formatting
4. Test error responses

### Client Tests
1. Mock HTTP responses
2. Test network error handling
3. Validate data transformation
4. Test timeout scenarios

## 🎉 Success Metrics

- ✅ **100% Test Execution**: All tests passing
- ✅ **≥80% Coverage**: Meeting industry standards
- ✅ **0 Race Conditions**: Concurrent safety validated
- ✅ **Performance Benchmarks**: Meeting latency targets
- ✅ **Clean Architecture**: Well-isolated unit tests

---

## 📞 Support & Resources

- **Test Framework**: [Testify](https://github.com/stretchr/testify)
- **Coverage Tools**: Built-in Go tools
- **HTTP Testing**: Standard library `httptest`
- **Benchmarking**: Built-in Go benchmarking

For questions about testing patterns or adding new tests, refer to existing test files as examples of best practices. 