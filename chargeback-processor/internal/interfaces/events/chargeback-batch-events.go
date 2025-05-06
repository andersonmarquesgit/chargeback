package events

import (
	"time"
)

type Event struct {
}

type ChargebackBatchEvent struct {
	FileID      string    `json:"file_id"`
	FileURL     string    `json:"file_url"`
	CreatedAt   time.Time `json:"created_at"`
	RecordCount int       `json:"record_count"`
}

// NewChargebackOpenedEvent cria um novo ChargebackOpenedEvent
func NewChargebackOpenedEvent(fileID, fileURL string, createdAt time.Time, recordCount int) (*ChargebackBatchEvent, error) {
	return &ChargebackBatchEvent{
		FileID:      fileID,
		FileURL:     fileURL,
		CreatedAt:   createdAt,
		RecordCount: recordCount,
	}, nil
}
