package usecases

import (
	"api/internal/domain/models"
	"api/internal/infrastructure/logging"
	"api/internal/infrastructure/rabbitmq/producers"
	"api/internal/infrastructure/repositories/cassandra"
	"api/internal/interfaces/dto"
	"api/internal/interfaces/events"
	"encoding/json"
	"fmt"
)

type ChargebackOpenedUseCase struct {
	ChargebackRepositoryCassandra cassandra.ChargebackRepositoryCassandra
	Producer                      producers.Producer
}

func NewChargebackOpenedUseCase(repository *cassandra.ChargebackRepositoryCassandra, producer *producers.Producer) *ChargebackOpenedUseCase {
	return &ChargebackOpenedUseCase{
		ChargebackRepositoryCassandra: *repository,
		Producer:                      *producer,
	}
}

func (uc *ChargebackOpenedUseCase) CreateChargeback(req dto.ChargebackRequest, traceID string) (*models.Chargeback, error) {
	// Verificamos se o chargeback já existe (CQRS - Query)
	logging.Info("Verifying if chargeback already exists")
	chargebackExists, err := uc.ChargebackRepositoryCassandra.GetChargeback(req.UserID, req.TransactionID)
	if err != nil {
		return nil, err
	}

	// Caso o chargeback não exista, criamos o evento para o processor criar (CQRS - Command)
	if chargebackExists != nil {
		logging.Info("Chargeback already exists")
		chargebackExists.Exists = true
		return chargebackExists, nil
	}

	logging.Info("Chargeback not exists, send event to processor")

	chargeback := models.NewChargeback(req.UserID, req.TransactionID, req.Reason)

	chargebackOpenedEventPayload, err := events.NewChargebackOpenedEvent(chargeback)
	if err != nil {
		return nil, fmt.Errorf("could not create chargeback opened event payload: %v", err)
	}

	message, err := json.Marshal(chargebackOpenedEventPayload)
	if err != nil {
		return nil, err
	}

	// Publicação do evento de chargeback opened
	err = uc.Producer.PublishChargebackOpenedEvent(message, req.UserID, req.TransactionID, traceID)
	if err != nil {
		return nil, err
	}

	return chargeback, nil
}
