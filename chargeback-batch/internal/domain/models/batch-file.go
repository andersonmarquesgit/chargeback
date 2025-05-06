package models

import "time"

type BatchFile struct {
	FileID        string    `json:"file_id"`
	FileURL       string    `json:"file_url"`
	RecordCount   int       `json:"record_count"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	SentAt        time.Time `json:"sent_at"`
	RetryCount    int       `json:"retry_count"`
	LastAttemptAt time.Time `json:"last_attempt_at"`
}
