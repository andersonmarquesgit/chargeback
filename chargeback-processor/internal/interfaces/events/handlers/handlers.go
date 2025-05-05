package handlers

type EventHandlers struct {
	ChargebackOpenedEventHandler *ChargebackOpenedEventHandler
}

func NewEventHandlers(handler *ChargebackOpenedEventHandler) *EventHandlers {
	return &EventHandlers{
		ChargebackOpenedEventHandler: handler,
	}
}
