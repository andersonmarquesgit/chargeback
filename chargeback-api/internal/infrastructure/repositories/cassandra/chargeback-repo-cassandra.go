package cassandra

import (
	"api/internal/domain/models"
	"github.com/gocql/gocql"
)

type ChargebackRepositoryCassandra struct {
	session *gocql.Session
}

func NewChargebackRepositoryCassandra(session *gocql.Session) *ChargebackRepositoryCassandra {
	return &ChargebackRepositoryCassandra{
		session: session,
	}
}

func (r *ChargebackRepositoryCassandra) GetChargeback(userID, transactionID string) (*models.Chargeback, error) {
	query := `SELECT status, reason, created_at, updated_at 
			  FROM chargebacks_by_user_transaction 
			  WHERE user_id = ? AND transaction_id = ? LIMIT 1`

	var cb models.Chargeback
	err := r.session.Query(query, userID, transactionID).Consistency(gocql.One).Scan(
		&cb.Status,
		&cb.Reason,
		&cb.CreatedAt,
		&cb.UpdatedAt,
	)
	if err == gocql.ErrNotFound {
		return nil, nil // não existe
	}
	if err != nil {
		return nil, err
	}

	cb.UserID = userID
	cb.TransactionID = transactionID
	return &cb, nil
}
