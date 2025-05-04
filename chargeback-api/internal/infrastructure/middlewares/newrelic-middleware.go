package middlewares

import (
	"api/internal/infrastructure/logging"
	"context"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"net/http"
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

			// Get country from request header
			country := r.Header.Get("country")
			if country == "" {
				country = "unknown"
			}

			// Add trace ID and country to transaction and context
			newRelicTransaction.AddAttribute("request-trace-id", traceID)
			newRelicTransaction.AddAttribute("country", country)
			ctx := newrelic.NewContext(r.Context(), newRelicTransaction)
			ctx = context.WithValue(ctx, "request-trace-id", traceID)
			ctx = context.WithValue(ctx, "country", country)
			r = r.WithContext(ctx)

			// Set response header
			w.Header().Set("request-trace-id", traceID)
			w.Header().Set("country", country)

			// Set trace ID and country in the logger
			logging.SetFields(logrus.Fields{
				"request-trace-id": traceID,
				"country":          country,
			})

			logging.Logger.AddHook(&logging.TraceIDHook{
				TraceIDKey: "request-trace-id",
				Context:    ctx,
			})

			logging.Logger.AddHook(&logging.CountryHook{
				CountryKey: "country",
				Context:    ctx,
			})

			// Serve the next handler
			next.ServeHTTP(w, r)
		})
	}
}
