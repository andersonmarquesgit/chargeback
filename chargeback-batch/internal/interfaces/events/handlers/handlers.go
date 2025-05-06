package handlers

type EventHandlers struct {
	ChargebackBatchEventHandler *ChargebackBatchEventHandler
}

func NewEventHandlers(handler *ChargebackBatchEventHandler) *EventHandlers {
	return &EventHandlers{
		ChargebackBatchEventHandler: handler,
	}
}
