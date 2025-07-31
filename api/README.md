# 🔧 Stock Tracker API

Go backend for the Stock Tracker platform with real-time data ingestion, AI recommendations, and user management.

## 🚀 Quick Start

```bash
# Install dependencies
go mod download

# Run API server
go run cmd/api/main.go

# Run data ingestion
go run cmd/ingestor/main.go

# Run migrations
go run cmd/migrator/main.go -direction=up
```

## 🏗️ Architecture

### Domain-Driven Design
```
internal/
├── domain/           # Business logic & use cases
│   ├── authentication/  # User auth & JWT
│   ├── stocks/         # Stock management
│   ├── recommendation/  # AI recommendations
│   ├── marketdata/     # Market data analysis
│   ├── subscription/   # Subscription management
│   └── chat/          # AI chat functionality
├── infrastructure/   # External dependencies
│   ├── database/      # CockroachDB repositories
│   ├── external/      # Yahoo Finance, Alpha Vantage APIs
│   ├── cache/         # In-memory caching
│   └── middleware/    # Auth, rate limiting, CORS
└── presentation/     # HTTP handlers & routes
```

## 📊 Data Ingestion

### Real-time Market Data
- **Yahoo Finance Integration**: Live stock prices and market data
- **Scheduled Ingestion**: Cron jobs every hour
- **Market Analysis**: Technical indicators and trends

### Ingestion Commands
```bash
# Run ingestion manually
go run cmd/ingestor/main.go

# Stock data ingestion
go run cmd/ingestor/main.go -type=stocks

# Market data ingestion  
go run cmd/ingestor/main.go -type=market
```

## 🗄️ Database Management

### Migrations
```bash
# Apply all migrations
go run cmd/migrator/main.go -direction=up

# Rollback migrations
go run cmd/migrator/main.go -direction=down

# Run specific migration
go run cmd/migrator/main.go -migration=004 -direction=up

# Reset database (down + up)
go run cmd/migrator/main.go -migration=004 -direction=reset
```

### Database Schema
- **Users**: Authentication and profiles
- **Stocks**: Stock data and metadata
- **Recommendations**: AI-generated stock recommendations
- **Market Data**: Real-time market analysis
- **Subscriptions**: User subscription tiers
- **Chat Sessions**: AI chat interactions

## 🔐 Authentication & Authorization

### JWT-based Auth
- **Access Tokens**: Short-lived (1 hour)
- **Refresh Tokens**: Long-lived (7 days)
- **User Tiers**: Guest, Basic, Premium

### User Tiers
- **Guest**: 100 req/hour, 10 recommendations
- **Basic**: 500 req/hour, 25 recommendations
- **Premium**: 2000 req/hour, 100 recommendations

## 🛠️ Development

### Environment Setup
```bash
# Required environment variables
DATABASE_URL=postgresql://user:pass@host:26257/db?sslmode=verify-full&sslrootcert=certs/cc-ca.crt
JWT_SECRET=your-jwt-secret
YAHOO_FINANCE_API_KEY=your-api-key
PORT=8080
LOG_LEVEL=info
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific tests
go test ./internal/domain/stocks/...

# Test with coverage
go test -cover ./...
```

### Linting & Formatting
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Security analysis
gosec ./...
```

## 📡 API Endpoints

### Core Endpoints
- `GET /api/v1/stocks` - List stocks
- `GET /api/v1/recommendations` - AI recommendations
- `GET /api/v1/market-data/analysis/{ticker}` - Market analysis
- `POST /api/v1/auth/login` - User authentication
- `GET /api/v1/subscriptions/plans` - Subscription plans

### Health & Monitoring
- `GET /health` - Health check
- `GET /ping` - Simple ping
- `GET /api/v1/stocks/stats` - System statistics

## 🔧 Configuration

### Server Settings
- **Port**: Configurable via PORT env var
- **Timeouts**: 15s read, 15s write, 60s idle
- **CORS**: Enabled for cross-origin requests
- **Compression**: Gzip enabled
- **Rate Limiting**: Tier-based limits

### External Services
- **CockroachDB Cloud**: Primary database
- **Yahoo Finance**: Market data source
- **Alpha Vantage**: Alternative data source
- **Redis**: Caching (optional)

## 📚 Documentation

For complete API documentation, see:
- **[API Documentation](../docs/API_Documentation.md)** - Full endpoint reference
- **[Postman Collection](../Stock-Tracker-API.postman_collection.json)** - API examples

---

**Built with Go 1.24.4 • Chi Router • CockroachDB • JWT Auth** 