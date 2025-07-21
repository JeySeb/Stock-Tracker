package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	StockAPIURL string
	StockAPIKey string
	LogLevel    string
	Port        string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/stock_system"),
		StockAPIURL: getEnv("STOCK_API_URL", "https://api.example.com/stocks"),
		StockAPIKey: getEnv("STOCK_API_KEY", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Port:        getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
