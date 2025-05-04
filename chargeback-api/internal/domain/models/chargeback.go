package models

import (
	"time"
)

type Chargeback struct {
	ChargebackID  string    `json:"chargeback_id"`
	Status        string    `json:"status"`
	UserID        string    `json:"user_id"`
	TransactionID string    `json:"transaction_id"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewChargeback(userID, transactionID, reason string) *Chargeback {
	return &Chargeback{
		ChargebackID:  userID + transactionID,
		Status:        "pending",
		UserID:        userID,
		TransactionID: transactionID,
		Reason:        reason,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
