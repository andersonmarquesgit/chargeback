package handlers

// Handlers struct to store the customer handler and other handlers
type Handlers struct {
	ChargebackOpenedHandler *ChargebackOpenedHandler
}

// NewHandlers struct
func NewHandlers(chargebackOpenedHandler *ChargebackOpenedHandler) *Handlers {
	return &Handlers{
		ChargebackOpenedHandler: chargebackOpenedHandler,
	}
}
