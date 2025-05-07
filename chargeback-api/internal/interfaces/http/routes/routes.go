// @title Chargeback API
// @version 1.0
// @description This is Chargeback API. Is possible to open a new chargeback.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
package routes

import (
	_ "api/docs"
	"api/internal/interfaces/http/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Route struct {
	URI          string
	Method       string
	Handler      http.HandlerFunc
	RequiresAuth bool
}

// SetupRoutes configure routes and middlewares
func SetupRoutes(router chi.Router, handlers *handlers.Handlers) {
	routes := NewChargebackRoutes(handlers.ChargebackOpenedHandler)

	for _, route := range routes {
		router.Method(route.Method, route.URI, route.Handler)
	}
}
