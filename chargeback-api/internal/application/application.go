package application

import (
	"api/internal/config"
	"api/internal/infrastructure/logging"
	"api/internal/infrastructure/rabbitmq"
	"api/internal/infrastructure/rabbitmq/producers"
	"api/internal/infrastructure/repositories/cassandra"
	"api/internal/interfaces/http/handlers"
	"api/internal/usecases"
	"github.com/gocql/gocql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"time"
)

type Application struct {
	Config      *config.Config
	Cassandra   *gocql.Session
	Producers   *producers.Producers
	UseCases    *usecases.UseCases
	Handlers    *handlers.Handlers
	NewRelicApp *newrelic.Application
}

func NewApplication(cfg *config.Config) *Application {
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("chargeback-api"),
		newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		log.Fatalf("Failed to create New Relic application: %v", err)
	}

	// Initialize the global logger
	logging.InitializeLogger(newRelicApp)
	logging.Logger.Info("Application started with New Relic integration")

	// Connect to database
	cassandraSession := config.NewCassandraSession(cfg.Database)

	// Connect to RabbitMQ
	rabbitConn, err := rabbitmq.Connect(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}

	// Initialize producers
	producers, err := producers.NewProducers(rabbitConn)
	if err != nil {
		log.Fatalf("Could not create RabbitMQ producers: %v", err)
	}

	// Initialize repositories
	chargebackRepository := cassandra.NewChargebackRepositoryCassandra(cassandraSession)

	// Initialize use cases
	chargebackUseCase := usecases.NewChargebackOpenedUseCase(chargebackRepository, producers.ChargebackOpenedProducer)
	useCases := usecases.NewUseCases(chargebackUseCase)

	// Initialize specific handlers
	chargebackOpenedHandler := handlers.NewChargebackOpenedHandler(useCases.ChargebackOpenedUseCase)

	// Initialize handlers
	handlers := handlers.NewHandlers(chargebackOpenedHandler)

	return &Application{
		Config:      cfg,
		Cassandra:   cassandraSession,
		Producers:   producers,
		UseCases:    useCases,
		Handlers:    handlers,
		NewRelicApp: newRelicApp,
	}
}

func (app *Application) Close() {
	if app.Cassandra != nil {
		app.Cassandra.Close()
	}
	if app.Producers != nil {
		app.Producers.ChargebackOpenedProducer.Connection.Close()
		// Close other producers as needed
	}
	if app.NewRelicApp != nil {
		app.NewRelicApp.Shutdown(10 * time.Second)
	}
}
