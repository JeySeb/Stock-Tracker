package entities

import (
	"time"
	"github.com/google/uuid"
)

type IngestionStatus string

const (
	IngestionStatusRunning   IngestionStatus = "running"
	IngestionStatusCompleted IngestionStatus = "completed"
	IngestionStatusFailed    IngestionStatus = "failed"
)

type IngestionLog struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	BatchID           string          `json:"batch_id" db:"batch_id"`
	TotalRecords      int             `json:"total_records" db:"total_records"`
	SuccessfulRecords int             `json:"successful_records" db:"successful_records"`
	FailedRecords     int             `json:"failed_records" db:"failed_records"`
	Status            IngestionStatus `json:"status" db:"status"`
	ErrorDetails      map[string]interface{} `json:"error_details,omitempty" db:"error_details"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	CompletedAt       *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
}

func NewIngestionLog( batchID string, totalRecords int ) *IngestionLog {
	return &IngestionLog{
		ID:                uuid.New(),
		BatchID:           batchID,
		TotalRecords:      totalRecords,
		SuccessfulRecords: 0,
		FailedRecords:     0,
		Status:            IngestionStatusRunning,
		CreatedAt:         time.Now(),
	}
}

func (il *IngestionLog) Complete() {
	il.Status = IngestionStatusCompleted
	now := time.Now()
	il.CompletedAt = &now
}

func (il *IngestionLog) Fail(errorDetails map[string]interface{}) {
	il.Status = IngestionStatusFailed
	il.ErrorDetails = errorDetails
	now := time.Now()
	il.CompletedAt = &now
}