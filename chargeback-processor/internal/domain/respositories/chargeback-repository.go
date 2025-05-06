package respositories

import (
	"processor/internal/domain/models"
	"time"
)

type ChargebackRepository interface {
	CreateChargeback(userID, transactionID, reason, fileID string, createdAt, updatedAt time.Time) (*models.Chargeback, error)
}
