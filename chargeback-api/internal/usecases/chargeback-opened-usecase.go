package usecases

import (
	"api/internal/domain/models"
	"api/internal/infrastructure/logging"
	"api/internal/infrastructure/rabbitmq/producers"
	"api/internal/interfaces/dto"
	"encoding/json"
)

type ChargebackOpenedUseCase struct {
	Producer producers.Producer
}

func NewChargebackOpenedUseCase(producer *producers.Producer) *ChargebackOpenedUseCase {
	return &ChargebackOpenedUseCase{
		Producer: *producer,
	}
}

func (uc *ChargebackOpenedUseCase) CreateChargeback(req dto.ChargebackRequest, traceID string) (*models.Chargeback, error) {
	chargeback := models.NewChargeback(req.UserID, req.TransactionID, req.Reason)

	// Verificamos se o chargeback já existe (CQRS - Query)
	logging.Info("Verifying if chargeback already exists")
	//createdCustomer, err := uc.CustomerRepository.CreateCustomer(customer)
	//if err != nil {
	//	return nil, err
	//}

	// Caso o chargeback não exista, criamos o evento para o processor criar (CQRS - Command)

	// Criação do payload e encriptação da senha do cliente para publicação do evento de criação do usuário
	//chargebackOpenedEventPayload, err := events.NewUserUpsertEvent(createdCustomer, customerReq.Password)
	//if err != nil {
	//	return nil, fmt.Errorf("could not create user event payload: %v", err)
	//}
	//userEventPayload.EncryptPassword()

	message, err := json.Marshal("")
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
