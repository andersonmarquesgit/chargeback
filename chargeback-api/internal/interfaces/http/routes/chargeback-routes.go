package routes

import (
	"api/internal/interfaces/http/handlers"
	"net/http"
)

func NewChargebackRoutes(handlers *handlers.ChargebackOpenedHandler) []Route {
	return []Route{
		{
			URI:          "/v1/chargebacks",
			Method:       http.MethodPost,
			Handler:      handlers.CreateChargeback,
			RequiresAuth: false,
		},
		// Add other routes as needed
	}
}
