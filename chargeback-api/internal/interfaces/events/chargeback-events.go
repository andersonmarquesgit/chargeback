package events

import (
	"api/internal/domain/models"
	"time"
)

type Event struct {
}

type ChargebackOpenedEvent struct {
	Status        string    `json:"status"`
	UserID        string    `json:"user_id"`
	TransactionID string    `json:"transaction_id"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewUserUpsertEvent cria um novo ChargebackOpenedEvent a partir de um Chargeback model
func NewChargebackOpenedEvent(chargeback *models.Chargeback) (*ChargebackOpenedEvent, error) {
	return &ChargebackOpenedEvent{
		UserID:        chargeback.UserID,
		TransactionID: chargeback.TransactionID,
		Reason:        chargeback.Reason,
		Status:        chargeback.Status,
		CreatedAt:     chargeback.CreatedAt,
		UpdatedAt:     chargeback.UpdatedAt,
	}, nil
}
