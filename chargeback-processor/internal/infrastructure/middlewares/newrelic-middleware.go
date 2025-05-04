package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"net/http"
	"processor/internal/infrastructure/logging"
)

func NewRelicMiddleware(newRelicApp *newrelic.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start the New Relic transaction
			newRelicTransaction := newRelicApp.StartTransaction(r.URL.Path)
			defer newRelicTransaction.End()

			// Add trace ID if not present
			traceID := r.Header.Get("request-trace-id")
			if traceID == "" {
				traceID = uuid.New().String()
			}

			// Add trace ID and country to transaction and context
			newRelicTransaction.AddAttribute("request-trace-id", traceID)
			ctx := newrelic.NewContext(r.Context(), newRelicTransaction)
			ctx = context.WithValue(ctx, "request-trace-id", traceID)
			r = r.WithContext(ctx)

			// Set response header
			w.Header().Set("request-trace-id", traceID)

			// Set trace ID and country in the logger
			logging.SetFields(logrus.Fields{
				"request-trace-id": traceID,
			})

			logging.Logger.AddHook(&logging.TraceIDHook{
				TraceIDKey: "request-trace-id",
				Context:    ctx,
			})

			// Serve the next handler
			next.ServeHTTP(w, r)
		})
	}
}
