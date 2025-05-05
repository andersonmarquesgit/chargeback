package respositories

import (
	"processor/internal/domain/models"
)

type ChargebackRepository interface {
	CreateChargeback(userID, transactionID, reason string) (*models.Chargeback, error)
}
