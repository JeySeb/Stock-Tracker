package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Database
	DatabaseURL string

	// External APIs
	StockAPIURL string
	StockAPIKey string

	// Server
	LogLevel string
	Port     string

	// JWT Configuration
	JWTSecret          string
	JWTAccessTokenTTL  time.Duration
	JWTRefreshTokenTTL time.Duration
	JWTIssuer          string

	// Security
	BCryptCost       int
	RateLimitEnabled bool

	// Environment
	Environment string
}

func LoadConfig() (*Config, error) {
	return &Config{
		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/stock_system"),

		// External APIs
		StockAPIURL: getEnv("STOCK_API_URL", "https://api.example.com/stocks"),
		StockAPIKey: getEnv("STOCK_API_KEY", ""),

		// Server
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Port:     getEnv("PORT", "8080"),

		// JWT Configuration
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTAccessTokenTTL:  getDurationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
		JWTRefreshTokenTTL: getDurationEnv("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		JWTIssuer:          getEnv("JWT_ISSUER", "stock-tracker"),

		// Security
		BCryptCost:       getIntEnv("BCRYPT_COST", 12),
		RateLimitEnabled: getBoolEnv("RATE_LIMIT_ENABLED", true),

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
