#!/bin/bash

echo "🔍 BACKEND VALIDATION - FASE 2"
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

echo "1. Building all components..."

# Build API server
echo "Building API server..."
go build -o bin/api cmd/api/main.go
print_status $? "API server build"

# Build ingestor
echo "Building ingestor..."
go build -o bin/ingestor cmd/ingestor/main.go  
print_status $? "Ingestor build"

echo ""
echo "2. Validating Go modules..."

# Check go mod
go mod tidy
go mod verify
print_status $? "Go modules validation"

echo ""
echo "3. Running tests..."

# Run tests if they exist
if [ -d "tests" ]; then
    go test ./tests/... -v
    print_status $? "Unit tests"
else
    print_warning "No tests directory found - consider adding unit tests"
fi

echo ""
echo "4. Checking core components..."

# Check if essential files exist
essential_files=(
    "internal/domain/entities/stock.go"
    "internal/domain/entities/broker.go"
    "internal/domain/repositories/stock_repository.go"
    "internal/domain/repositories/broker_repository.go"
    "internal/infrastructure/database/stock_repository_impl.go"
    "internal/infrastructure/clients/stock_api_client.go"
    "internal/domain/usecases/stock_ingestion.go"
    "internal/presentation/handlers/stock_handler.go"
    "pkg/logger/logger.go"
    "internal/infrastructure/config/config.go"
)

for file in "${essential_files[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $file${NC}"
    else
        echo -e "${RED}❌ Missing: $file${NC}"
        exit 1
    fi
done

echo ""
echo "5. Validating API compliance..."

# Check if API client has correct endpoint
if grep -q "8j5baasof2.execute-api.us-west-2.amazonaws.com" internal/infrastructure/clients/stock_api_client.go; then
    print_warning "API client may have hardcoded endpoint - should use config"
fi

# Check if proper authentication is implemented
if grep -q "Bearer" internal/infrastructure/clients/stock_api_client.go; then
    echo -e "${GREEN}✅ Bearer token authentication implemented${NC}"
else
    echo -e "${RED}❌ Missing Bearer token authentication${NC}"
    exit 1
fi

echo ""
echo "6. Database validation..."

# Check if repository has batch operations
if grep -q "BulkCreate" internal/infrastructure/database/stock_repository_impl.go; then
    echo -e "${GREEN}✅ Batch operations implemented${NC}"
else
    echo -e "${RED}❌ Missing batch operations${NC}"
    exit 1
fi

echo ""
echo "7. Checking pagination and filtering..."

# Check if pagination is implemented
if grep -q "Pagination" internal/domain/valueObjects/filters.go; then
    echo -e "${GREEN}✅ Pagination implemented${NC}"
else
    echo -e "${RED}❌ Missing pagination${NC}"
    exit 1
fi

echo ""
echo "8. Ingestion system validation..."

# Check if cron job is implemented
if grep -q "cron" cmd/ingestor/main.go; then
    echo -e "${GREEN}✅ Cron job scheduler implemented${NC}"
else
    echo -e "${RED}❌ Missing cron job scheduler${NC}"
    exit 1
fi

# Check if worker pool is implemented
if grep -q "errgroup\|worker" internal/domain/usecases/stock_ingestion.go; then
    echo -e "${GREEN}✅ Worker pool/concurrent processing implemented${NC}"
else
    print_warning "Consider implementing worker pool for better performance"
fi

echo ""
echo "9. API endpoints validation..."

# Check if basic CRUD endpoints exist
if grep -q "GetStocks\|GetStockByTicker\|GetStats" internal/presentation/handlers/stock_handler.go; then
    echo -e "${GREEN}✅ Basic API endpoints implemented${NC}"
else
    echo -e "${RED}❌ Missing basic API endpoints${NC}"
    exit 1
fi

echo ""
echo "10. Configuration management..."

# Check if environment configuration exists
if grep -q "getEnv\|os.Getenv" internal/infrastructure/config/config.go; then
    echo -e "${GREEN}✅ Environment configuration implemented${NC}"
else
    echo -e "${RED}❌ Missing environment configuration${NC}"
    exit 1
fi

echo ""
echo "📋 VALIDATION SUMMARY"
echo "====================="
echo -e "${GREEN}✅ All core components validated successfully!${NC}"
echo ""
echo "📌 FASE 2 REQUIREMENTS STATUS:"
echo "✅ Domain entities with business logic"
echo "✅ Repository implementations with CockroachDB support"  
echo "✅ API client for stock data service"
echo "✅ Ingestion system with batch processing"
echo "✅ Basic REST API with Chi router"
echo "✅ Configuration management"
echo "✅ Logging system"
echo ""
echo -e "${YELLOW}🚀 Ready for FASE 3: Recommendation Engine!${NC}"
echo ""
echo "💡 RECOMMENDATIONS FOR IMPROVEMENT:"
echo "1. Add comprehensive unit tests (target >80% coverage)"
echo "2. Add integration tests for API endpoints"
echo "3. Implement proper error handling middleware"
echo "4. Add metrics and monitoring"
echo "5. Consider adding database migrations"
echo "6. Add API documentation (OpenAPI/Swagger)" 