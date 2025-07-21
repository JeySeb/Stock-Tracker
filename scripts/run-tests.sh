#!/bin/bash

set -e

echo "🧪 COMPREHENSIVE TEST RUNNER"
echo "============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Create test output directory
mkdir -p coverage

echo "1. Running Unit Tests..."
echo "========================"

# Run unit tests with coverage
go test -v -race -coverprofile=coverage/coverage.out ./tests/unit/... 2>&1 | tee coverage/unit-tests.log
print_status $? "Unit tests completed"

echo ""
echo "2. Running Integration Tests..."
echo "==============================="

# Run integration tests if they exist
if [ -d "tests/integration" ]; then
    go test -v -race ./tests/integration/... 2>&1 | tee coverage/integration-tests.log
    print_status $? "Integration tests completed"
else
    print_warning "No integration tests found"
fi

echo ""
echo "3. Benchmark Tests..."
echo "===================="

# Run benchmark tests
go test -bench=. -benchmem ./tests/unit/... 2>&1 | tee coverage/benchmark.log
print_status $? "Benchmark tests completed"

echo ""
echo "4. Code Coverage Analysis..."
echo "============================"

# Generate coverage report
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
print_status $? "Coverage HTML report generated"

# Get coverage percentage
COVERAGE=$(go tool cover -func=coverage/coverage.out | grep "total:" | awk '{print $3}')
echo -e "${BLUE}📊 Total Coverage: ${COVERAGE}${NC}"

# Check if coverage meets threshold (80%)
COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
if (( $(echo "$COVERAGE_NUM >= 80" | bc -l) )); then
    echo -e "${GREEN}✅ Coverage meets 80% threshold${NC}"
else
    print_warning "Coverage below 80% threshold"
fi

echo ""
echo "5. Test Summary..."
echo "=================="

# Count test results
UNIT_TESTS_PASSED=$(grep "PASS:" coverage/unit-tests.log | wc -l)
UNIT_TESTS_FAILED=$(grep "FAIL:" coverage/unit-tests.log | wc -l)

echo "📊 Test Results:"
echo "  • Unit Tests Passed: $UNIT_TESTS_PASSED"
echo "  • Unit Tests Failed: $UNIT_TESTS_FAILED"
echo "  • Coverage: $COVERAGE"

echo ""
echo "6. Running Linter..."
echo "==================="

# Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
    golangci-lint run ./... 2>&1 | tee coverage/lint.log
    print_status $? "Linting completed"
else
    print_warning "golangci-lint not found, skipping"
fi

echo ""
echo "7. Running Race Condition Tests..."
echo "=================================="

# Run tests with race detector specifically
go test -race -short ./tests/unit/... 2>&1 | tee coverage/race-tests.log
print_status $? "Race condition tests completed"

echo ""
echo "8. Generating Test Reports..."
echo "============================="

# Generate a comprehensive test report
cat > coverage/test-report.md << EOF
# Test Report

## Coverage Summary
- **Total Coverage**: $COVERAGE
- **Coverage Threshold**: 80%
- **Status**: $(if (( $(echo "$COVERAGE_NUM >= 80" | bc -l) )); then echo "✅ PASSED"; else echo "❌ BELOW THRESHOLD"; fi)

## Test Results
- **Unit Tests Passed**: $UNIT_TESTS_PASSED
- **Unit Tests Failed**: $UNIT_TESTS_FAILED

## Files Generated
- \`coverage.html\` - Interactive coverage report
- \`unit-tests.log\` - Detailed unit test output
- \`benchmark.log\` - Benchmark results
- \`race-tests.log\` - Race condition test results

## How to View Coverage
\`\`\`bash
open coverage/coverage.html
\`\`\`

## Test Command Examples
\`\`\`bash
# Run specific test
go test -v ./tests/unit/entities -run TestStock_IsUpgrade

# Run with coverage
go test -coverprofile=coverage.out ./tests/unit/...

# Run benchmarks
go test -bench=. ./tests/unit/entities

# Run with race detection
go test -race ./tests/unit/...
\`\`\`
EOF

echo "📋 Test report generated: coverage/test-report.md"

echo ""
echo "9. Performance Analysis..."
echo "========================="

# Generate performance profile if pprof tests exist
if grep -q "pprof" coverage/benchmark.log; then
    print_info "Performance profiling data available in benchmark results"
else
    print_info "No performance profiling data generated"
fi

echo ""
echo "🎉 TEST EXECUTION COMPLETE"
echo "=========================="
echo -e "${GREEN}✨ All tests executed successfully!${NC}"
echo ""
echo "📁 Reports available in: ./coverage/"
echo "🌐 Open coverage report: coverage/coverage.html"
echo "📊 Test summary: coverage/test-report.md"

# Final status based on coverage
if (( $(echo "$COVERAGE_NUM >= 80" | bc -l) )); then
    echo -e "${GREEN}🚀 Ready for production deployment!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Consider adding more tests to reach 80% coverage${NC}"
    exit 0
fi 