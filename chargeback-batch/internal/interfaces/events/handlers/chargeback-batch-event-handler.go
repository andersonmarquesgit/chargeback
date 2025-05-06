package handlers

import (
	"batch/internal/infrastructure/logging"
	"batch/internal/usecases"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type ChargebackBatchEventHandler struct {
	ChargebackBatchEventUseCase *usecases.ChargebackBatchEventUseCase
}

func NewChargebackBatchEventHandler(chargebackBatchEventUseCase *usecases.ChargebackBatchEventUseCase) *ChargebackBatchEventHandler {
	return &ChargebackBatchEventHandler{
		ChargebackBatchEventUseCase: chargebackBatchEventUseCase,
	}
}

func (h *ChargebackBatchEventHandler) HandleChargebackOpenedEvent(msg amqp.Delivery) error {
	logging.Infof("Processing batch file event: %s", string(msg.Body))

	// Process the message
	err := h.ChargebackBatchEventUseCase.Process(msg)
	if err != nil {
		log.Printf("Failed to process batch file event: %v", err)
		return err
	}
	return nil
}
