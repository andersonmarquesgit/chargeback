package usecases

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"processor/internal/domain/models"
	"processor/internal/infrastructure/filewriter"
	"processor/internal/infrastructure/logging"
	"processor/internal/infrastructure/repositories/cassandra"
)

type ChargebackOpenedEventUseCase struct {
	ChargebackRepositoryCassandra *cassandra.ChargebackRepositoryCassandra
	ChargebackWriter              *filewriter.ChargebackWriter
}

func NewChargebackOpenedEventUseCase(chargebackRepositoryCassandra *cassandra.ChargebackRepositoryCassandra, chargebackWriter *filewriter.ChargebackWriter) *ChargebackOpenedEventUseCase {
	return &ChargebackOpenedEventUseCase{
		ChargebackRepositoryCassandra: chargebackRepositoryCassandra,
		ChargebackWriter:              chargebackWriter,
	}
}

func (uc *ChargebackOpenedEventUseCase) Process(msg amqp.Delivery) error {
	var cb models.Chargeback
	if err := json.Unmarshal(msg.Body, &cb); err != nil {
		logging.Infof("failed to unmarshal chargeback: %v", err)
		return err
	}

	fileID, err := uc.ChargebackWriter.Write(cb)
	if err != nil {
		logging.Infof("failed to write chargeback to file: %v", err)
		return err
	}

	// Persistir no Cassandra
	_, err = uc.ChargebackRepositoryCassandra.CreateChargeback(cb.UserID, cb.TransactionID, cb.Reason, fileID, cb.CreatedAt, cb.UpdatedAt)
	if err != nil {
		logging.Infof("failed to persist chargeback: %v", err)
		return err
	}

	logging.Infof("Chargeback written to file: %s and saved to Cassandra", fileID)
	return nil
}
