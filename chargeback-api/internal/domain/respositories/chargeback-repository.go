package respositories

import (
	"api/internal/domain/models"
)

type ChargebackRepository interface {
	GetChargeback(userID, transactionID string) (*models.Chargeback, error)
}
