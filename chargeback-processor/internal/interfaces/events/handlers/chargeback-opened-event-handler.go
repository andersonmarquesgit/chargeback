package handlers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"processor/internal/infrastructure/logging"
	"processor/internal/usecases"
)

type ChargebackOpenedEventHandler struct {
	ChargebackOpenedEventUseCase *usecases.ChargebackOpenedEventUseCase
}

func NewChargebackOpenedEventHandler(chargebackOpenedEventUseCase *usecases.ChargebackOpenedEventUseCase) *ChargebackOpenedEventHandler {
	return &ChargebackOpenedEventHandler{
		ChargebackOpenedEventUseCase: chargebackOpenedEventUseCase,
	}
}

func (h *ChargebackOpenedEventHandler) HandleChargebackOpenedEvent(msg amqp.Delivery) error {
	logging.Infof("Processing chargeback opened event: %s", string(msg.Body))

	// Process the message
	err := h.ChargebackOpenedEventUseCase.Process(msg)
	if err != nil {
		log.Printf("Failed to process send email event: %v", err)
		return err
	}
	return nil
}
