package usecases

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"processor/internal/infrastructure/logging"
	"processor/internal/infrastructure/repositories/cassandra"
)

type ChargebackOpenedEventUseCase struct {
	ChargebackRepositoryCassandra *cassandra.ChargebackRepositoryCassandra
}

func NewChargebackOpenedEventUseCase(chargebackRepositoryCassandra *cassandra.ChargebackRepositoryCassandra) *ChargebackOpenedEventUseCase {
	return &ChargebackOpenedEventUseCase{
		ChargebackRepositoryCassandra: chargebackRepositoryCassandra,
	}
}

func (uc *ChargebackOpenedEventUseCase) Process(msg amqp.Delivery) error {

	logging.Info("Chargeback receive successfully")

	return nil
}
