package model

import "time"

type AIGeneratedInsights struct {
	MarketSentiment     string    `json:"market_sentiment"`      // "Bullish", "Bearish", "Neutral"
	SentimentScore      float64   `json:"sentiment_score"`       // -1 to 1
	RiskAssessment      string    `json:"risk_assessment"`       // "Low", "Medium", "High"
	KeyDrivers          []string  `json:"key_drivers"`           // Main factors affecting the stock
	CompetitorAnalysis  []string  `json:"competitor_analysis"`   // Comparative insights
	NewsImpact          *float64  `json:"news_impact,omitempty"` // News sentiment impact
	TechnicalIndicators []string  `json:"technical_indicators"`  // Technical analysis signals
	GeneratedAt         time.Time `json:"generated_at"`
}
