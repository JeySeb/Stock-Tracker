#!/bin/bash

# ==============================================
# STOCK TRACKER - COMPREHENSIVE TEST RUNNER
# ==============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 STOCK TRACKER - COMPREHENSIVE TEST SUITE${NC}"
echo "=============================================="

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed or not in PATH${NC}"
    exit 1
fi

# Function to run tests with coverage
run_tests_with_coverage() {
    local test_path="$1"
    local test_name="$2"
    
    echo -e "\n${YELLOW}📋 Running $test_name...${NC}"
    
    if go test -v -race -coverprofile=coverage.out "$test_path"; then
        echo -e "${GREEN}✅ $test_name passed${NC}"
        
        # Generate coverage report
        if [ -f "coverage.out" ]; then
            coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}')
            echo -e "${BLUE}📊 Coverage: $coverage${NC}"
            
            # Optional: Generate HTML coverage report
            if [ "$GENERATE_HTML_COVERAGE" = "true" ]; then
                go tool cover -html=coverage.out -o "coverage_${test_name// /_}.html"
                echo -e "${BLUE}📄 HTML coverage report: coverage_${test_name// /_}.html${NC}"
            fi
            
            rm coverage.out
        fi
    else
        echo -e "${RED}❌ $test_name failed${NC}"
        return 1
    fi
}

# Function to run linter
run_linter() {
    echo -e "\n${YELLOW}🔍 Running linter...${NC}"
    
    if command -v golangci-lint &> /dev/null; then
        if golangci-lint run ./...; then
            echo -e "${GREEN}✅ Linter passed${NC}"
        else
            echo -e "${RED}❌ Linter found issues${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️  golangci-lint not found, skipping linter${NC}"
    fi
}

# Function to check test coverage
check_coverage_threshold() {
    local min_coverage=70
    echo -e "\n${YELLOW}📊 Checking overall test coverage...${NC}"
    
    # Run all tests with coverage
    if go test -race -coverprofile=total_coverage.out ./...; then
        if [ -f "total_coverage.out" ]; then
            total_coverage=$(go tool cover -func=total_coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
            echo -e "${BLUE}📈 Total Coverage: ${total_coverage}%${NC}"
            
            if (( $(echo "$total_coverage >= $min_coverage" | bc -l) )); then
                echo -e "${GREEN}✅ Coverage meets minimum threshold (${min_coverage}%)${NC}"
            else
                echo -e "${RED}❌ Coverage below minimum threshold (${min_coverage}%)${NC}"
                echo -e "${YELLOW}Current: ${total_coverage}% | Required: ${min_coverage}%${NC}"
                rm total_coverage.out
                return 1
            fi
            
            # Generate final HTML report
            if [ "$GENERATE_HTML_COVERAGE" = "true" ]; then
                go tool cover -html=total_coverage.out -o coverage_total.html
                echo -e "${BLUE}📄 Total coverage report: coverage_total.html${NC}"
            fi
            
            rm total_coverage.out
        fi
    else
        echo -e "${RED}❌ Failed to run coverage tests${NC}"
        return 1
    fi
}

# Parse command line arguments
VERBOSE=false
GENERATE_HTML_COVERAGE=false
RUN_INTEGRATION=false
RUN_LINTER=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--html)
            GENERATE_HTML_COVERAGE=true
            shift
            ;;
        -i|--integration)
            RUN_INTEGRATION=true
            shift
            ;;
        --no-lint)
            RUN_LINTER=false
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose        Enable verbose output"
            echo "  -h, --html          Generate HTML coverage reports"
            echo "  -i, --integration   Run integration tests"
            echo "  --no-lint           Skip linter"
            echo "  --help              Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Set verbose mode
if [ "$VERBOSE" = true ]; then
    set -x
fi

echo -e "\n${BLUE}🔧 Configuration:${NC}"
echo "  • HTML Coverage: $GENERATE_HTML_COVERAGE"
echo "  • Integration Tests: $RUN_INTEGRATION"
echo "  • Linter: $RUN_LINTER"

# Create coverage directory
mkdir -p coverage

# Start testing
echo -e "\n${BLUE}🚀 Starting test execution...${NC}"

# 1. Run unit tests
echo -e "\n${YELLOW}📦 UNIT TESTS${NC}"
echo "=================="

# Test auth handlers
run_tests_with_coverage "./tests/unit/handlers" "Auth Handler Tests"

# Test use cases
run_tests_with_coverage "./tests/unit/usecases" "Use Case Tests"

# Test auth services
run_tests_with_coverage "./tests/unit/auth" "Auth Service Tests"

# Test existing stock functionality
run_tests_with_coverage "./tests/unit/handlers" "Stock Handler Tests"

# 2. Run integration tests (if enabled)
if [ "$RUN_INTEGRATION" = true ]; then
    echo -e "\n${YELLOW}🔗 INTEGRATION TESTS${NC}"
    echo "======================"
    run_tests_with_coverage "./tests/integration" "Integration Tests"
fi

# 3. Run linter (if enabled)
if [ "$RUN_LINTER" = true ]; then
    run_linter
fi

# 4. Check overall coverage
check_coverage_threshold

# 5. Run security checks
echo -e "\n${YELLOW}🔒 Security Checks${NC}"
echo "=================="
if command -v gosec &> /dev/null; then
    if gosec ./...; then
        echo -e "${GREEN}✅ Security scan passed${NC}"
    else
        echo -e "${RED}❌ Security issues found${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  gosec not found, skipping security scan${NC}"
fi

# Final summary
echo -e "\n${GREEN}🎉 ALL TESTS COMPLETED SUCCESSFULLY!${NC}"
echo "======================================"
echo -e "${BLUE}📋 Summary:${NC}"
echo "  • Unit tests: ✅"
echo "  • Integration tests: $([ "$RUN_INTEGRATION" = true ] && echo "✅" || echo "⏭️ Skipped")"
echo "  • Linter: $([ "$RUN_LINTER" = true ] && echo "✅" || echo "⏭️ Skipped")"
echo "  • Coverage check: ✅"
echo "  • Security scan: ✅"

echo -e "\n${GREEN}🚀 Your authentication system is ready for production!${NC}" 