.PHONY: help setup dev-up dev-down build test lint clean migrate-up migrate-down

# Help
help: ## Show this help message
	@echo 'Usage: make [TARGET]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
setup: ## Initial project setup with CockroachDB Cloud
	@echo "🚀 Setting up development environment with CockroachDB Cloud..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Please create it with your DATABASE_URL"; \
		echo "Example: DATABASE_URL=postgresql://jeyseb:<password>@hiring-test-stock-cluster-13493.j77.aws-us-east-1.cockroachlabs.cloud:26257/stockdb?sslmode=verify-full&sslrootcert=certs/cc-ca.crt"; \
		exit 1; \
	fi
	@if [ ! -f certs/cc-ca.crt ]; then \
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
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	go run cmd/migrator/main.go -direction=up

migrate-down: ## Run database migrations down
	@echo "Running migrations down on CockroachDB Cloud..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	go run cmd/migrator/main.go -direction=down

migrate-reset: ## Reset database (down then up)
	make migrate-down
	make migrate-up

db-shell: ## Access CockroachDB Cloud shell
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | xargs) && psql "$$DATABASE_URL"

db-test-connection: ## Test CockroachDB Cloud connection
	@echo "🔍 Testing CockroachDB Cloud connection..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found"; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | xargs) && psql "$$DATABASE_URL" -c "SELECT version();" || (echo "❌ Connection failed" && exit 1)
	@echo "✅ Connection successful!"

##@ Backend
backend-deps: ## Install backend dependencies
	go mod tidy && go mod download

backend-run: ## Run backend locally
	go run cmd/api/main.go

backend-build: ## Build backend binary
	mkdir -p bin
	go build -o bin/api cmd/api/main.go

backend-test: ## Run backend tests
	go test -v ./...

backend-test-coverage: ## Run backend tests with coverage
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

backend-lint: ## Lint backend code
	golangci-lint run

##@ Frontend
frontend-deps: ## Install frontend dependencies
	cd webui && npm install

frontend-dev: ## Run frontend development server
	cd webui && npm run dev

frontend-build: ## Build frontend for production
	cd webui && npm run build

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