package models

import (
	"time"
)

type Chargeback struct {
	Status        string    `json:"status"`
	UserID        string    `json:"user_id"`
	TransactionID string    `json:"transaction_id"`
	Reason        string    `json:"reason"`
	Exists        bool      `json:"exists"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewChargeback(userID, transactionID, reason string) *Chargeback {
	return &Chargeback{
		Status:        "opened",
		UserID:        userID,
		TransactionID: transactionID,
		Reason:        reason,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
