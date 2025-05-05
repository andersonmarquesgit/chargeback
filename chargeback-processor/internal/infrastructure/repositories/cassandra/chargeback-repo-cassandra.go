package cassandra

import (
	"github.com/gocql/gocql"
	"processor/internal/domain/models"
	"time"
)

type ChargebackRepositoryCassandra struct {
	session *gocql.Session
}

func NewChargebackRepositoryCassandra(session *gocql.Session) *ChargebackRepositoryCassandra {
	return &ChargebackRepositoryCassandra{
		session: session,
	}
}

func (r *ChargebackRepositoryCassandra) CreateChargeback(userID, transactionID, reason string) (*models.Chargeback, error) {
	now := time.Now()

	cb := &models.Chargeback{
		UserID:        userID,
		TransactionID: transactionID,
		Status:        "opened",
		Reason:        reason,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	query := `INSERT INTO chargebacks_by_user_transaction 
		(user_id, transaction_id, status, reason, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`

	err := r.session.Query(query,
		cb.UserID,
		cb.TransactionID,
		cb.Status,
		cb.Reason,
		cb.CreatedAt,
		cb.UpdatedAt,
	).Exec()

	if err != nil {
		return nil, err
	}

	return cb, nil
}
