package model

import (
	"time"

	"github.com/google/uuid"
)

type BrokerCredibility struct {
	ID               uuid.UUID `json:"id" db:"id"`
	BrokerID         uuid.UUID `json:"broker_id" db:"broker_id"`
	BrokerName       string    `json:"broker_name" db:"broker_name"`
	TotalReports     int       `json:"total_reports" db:"total_reports"`
	AccuracyScore    float64   `json:"accuracy_score" db:"accuracy_score"`       // 0-1
	CredibilityScore float64   `json:"credibility_score" db:"credibility_score"` // 0-1
	LastUpdated      time.Time `json:"last_updated" db:"last_updated"`
}
