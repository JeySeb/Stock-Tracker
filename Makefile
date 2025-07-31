.PHONY: help setup dev-up dev-down build test lint clean migrate-up migrate-down

# Help
help: ## Show this help message
	@echo 'Usage: make [TARGET]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
api-setup: ## Initial project setup with CockroachDB Cloud
	@echo "🚀 Setting up development environment with CockroachDB Cloud..."
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found. Please create it with your DATABASE_URL"; \
		echo "Example: DATABASE_URL=postgresql://jeyseb:<password>@hiring-test-stock-cluster-13493.j77.aws-us-east-1.cockroachlabs.cloud:26257/stockdb?sslmode=verify-full&sslrootcert=certs/cc-ca.crt"; \
		exit 1; \
	fi
	cd api && \
	if [ ! -f certs/cc-ca.crt ]; then \
		echo "📥 Downloading CockroachDB Cloud SSL certificate..."; \
		mkdir -p certs; \
		curl -o certs/cc-ca.crt https://cockroachlabs.cloud/clusters/hiring-test-stock-cluster-13493/cert; \
	fi
	docker compose up -d redis localstack
	sleep 5
	make migrate-up
	@echo "✅ Setup complete!"

dev-up: ## Start development environment
	docker compose up -d redis localstack
	@echo "🔧 Development services are running"
	@echo "LocalStack: http://localhost:4566"
	@echo "Redis: localhost:6379"

dev-down: ## Stop development environment
	docker compose down

dev-logs: ## Follow development logs
	docker compose logs -f

##@ Database
migrate-up: ## Run database migrations up
	@echo "Running migrations up on CockroachDB Cloud..."
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	go run cmd/migrator/main.go -direction=up

migrate-down: ## Run database migrations down
	@echo "Running migrations down on CockroachDB Cloud..."
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	go run cmd/migrator/main.go -direction=down

migrate-reset: ## Reset database (down then up)
	make migrate-down
	make migrate-up

migrate-status: ## Show current migration status
	@echo "🔍 Checking migration status..."
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | xargs) && psql "$$DATABASE_URL" -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" 2>/dev/null || echo "❌ Could not get migration status"

migrate-specific: ## Run specific migration (usage: make migrate-specific MIGRATION=004 DIRECTION=reset)
	@if [ -z "$(MIGRATION)" ]; then \
		echo "❌ MIGRATION parameter is required"; \
		echo "Usage: make migrate-specific MIGRATION=004 DIRECTION=reset"; \
		echo "DIRECTION options: up, down, reset (default: reset)"; \
		exit 1; \
	fi
	@echo "🔧 Running migration $(MIGRATION) with direction $(or $(DIRECTION),reset)..."
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	go run cmd/migrator/main.go -migration=$(MIGRATION) -direction=$(or $(DIRECTION),reset)

db-reset: ## ⚠️ DANGER: Complete database reset - drops ALL tables and runs migrations fresh
	@echo "🚨 COMPLETE DATABASE RESET - This will destroy ALL data!"
	@echo "Are you sure? This action cannot be undone."
	@read -p "Type 'RESET' to confirm: " confirm && [ "$$confirm" = "RESET" ]
	./scripts/reset_and_migrate.sh

db-shell: ## Access CockroachDB Cloud shell
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | xargs) && psql "$$DATABASE_URL"

db-test-connection: ## Test CockroachDB Cloud connection
	@echo "🔍 Testing CockroachDB Cloud connection..."
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | xargs) && psql "$$DATABASE_URL" -c "SELECT version();" || (echo "❌ Connection failed" && exit 1)
	@echo "✅ Connection successful!"

##@ Backend
backend-deps: ## Install backend dependencies
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	go mod tidy && go mod download

backend-run: ## Run backend locally
	cd api && \
	go run cmd/api/main.go

backend-build: ## Build backend binary
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi && \
	mkdir -p bin && \
	go build -o bin/api cmd/api/main.go

backend-test: ## Run backend tests
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🧪 Running complete test suite..."
	./scripts/run-tests.sh

backend-lint: ## Lint backend code
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	golangci-lint run

backend-test-unit: ## Run unit tests only
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "📦 Running unit tests..."
	go test -v -race ./tests/unit/...

backend-test-api: ## Run API tests (auth, handlers, endpoints)
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🌐 Running API tests..."
	go test -v -race ./tests/unit/handlers/...
	@echo "🔐 Running auth service tests..."
	go test -v -race ./tests/unit/auth/...

backend-test-integration: ## Run integration tests
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🔗 Running integration tests..."
	go test -v -race ./tests/integration/...

backend-test-usecases: ## Run use case tests
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "⚙️ Running use case tests..."
	go test -v -race ./tests/unit/usecases/...

backend-test-coverage: ## Generate detailed coverage report
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "📊 Generating coverage report..."
	go test -race -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@coverage=$$(go tool cover -func=coverage/coverage.out | tail -1 | awk '{print $$3}'); \
	echo "📈 Total Coverage: $$coverage"

backend-test-coverage-html: ## Generate and open HTML coverage report
	make backend-test-coverage
	@echo "🌐 Opening coverage report..."
	@command -v xdg-open >/dev/null 2>&1 && xdg-open coverage/coverage.html || \
	command -v open >/dev/null 2>&1 && open coverage/coverage.html || \
	echo "📄 Coverage report: coverage/coverage.html"

backend-test-security: ## Run security tests
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🔒 Running security tests..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "❌ gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

backend-test-clean: ## Clean test artifacts
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🧹 Cleaning test artifacts..."
	rm -rf coverage/
	rm -f *.out *.html

backend-test-quick: ## Run quick tests (no race detection, no coverage)
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "⚡ Running quick tests..."
	go test -short ./tests/unit/...

##@ API Testing
api-test-auth: ## Test authentication endpoints specifically
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🔐 Testing authentication endpoints..."
	go test -v -run "TestAuth" ./tests/unit/handlers/...

api-test-subscription: ## Test subscription endpoints
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "💳 Testing subscription endpoints..."
	go test -v -run "TestSubscription" ./tests/unit/handlers/...

api-test-stocks: ## Test stock endpoints
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "📈 Testing stock endpoints..."
	go test -v -run "TestStock" ./tests/unit/handlers/...

api-validate: ## Validate all API endpoints are working
	cd api && \
	if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@echo "🌐 Validating API endpoints..."
	@if ! pgrep -f "cmd/api/main.go" > /dev/null; then \
		echo "⚠️  API server not running. Starting in background..."; \
		make backend-run & \
		sleep 3; \
		API_STARTED=true; \
	fi; \
	echo "🔍 Testing health endpoint..."; \
	curl -f http://localhost:8080/health || (echo "❌ Health check failed" && exit 1); \
	echo "✅ API endpoints validated"; \
	if [ "$$API_STARTED" = "true" ]; then \
		echo "🛑 Stopping test API server..."; \
		pkill -f "cmd/api/main.go"; \
	fi

##@ Frontend
frontend-deps: ## Install frontend dependencies
	cd webui && npm install

frontend-build: ## Build frontend for production
	cd webui && npm run build

frontend-dev: ## Run frontend development server
	cd webui && npm run dev

frontend-test: ## Run frontend tests
	cd webui && npm run test

frontend-lint: ## Lint frontend code
	cd webui && npm run lint

##@ Infrastructure
infra-plan-local: ## Plan Terraform for local environment
	cd infra/terraform/environments/local && terraform plan

infra-apply-local: ## Apply Terraform for local environment
	cd infra/terraform/environments/local && terraform apply

infra-destroy-local: ## Destroy Terraform local infrastructure
	cd infra/terraform/environments/local && terraform destroy

##@ Docker
docker-build-all: ## Build all Docker images
	docker compose build

docker-up-full: ## Start full application stack
	docker compose --profile backend --profile frontend up -d

docker-clean: ## Clean Docker resources
	docker compose down -v
	docker system prune -f

##@ Utilities
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf webui/dist/
	rm -f coverage.out coverage.html

check-deps: ## Check for required dependencies
	@command -v docker >/dev/null 2>&1 || (echo "❌ Docker is required" && exit 1)
	@docker compose version >/dev/null 2>&1 || (echo "❌ Docker Compose is required" && exit 1)
	@command -v go >/dev/null 2>&1 || (echo "❌ Go is required" && exit 1)
	@command -v node >/dev/null 2>&1 || (echo "❌ Node.js is required" && exit 1)
	@command -v psql >/dev/null 2>&1 || (echo "❌ PostgreSQL client is required" && exit 1)
	@echo "✅ All dependencies are installed"

status: ## Show status of all services
	@echo "📊 Service Status:"
	@docker compose ps
	@echo ""
	@echo "🔍 CockroachDB Cloud Status:"
	@make db-test-connection 2>/dev/null || echo "❌ Cannot connect to CockroachDB Cloud"