package application

import (
	"github.com/gocql/gocql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"processor/internal/config"
	"processor/internal/infrastructure/filewriter"
	"processor/internal/infrastructure/logging"
	"processor/internal/infrastructure/objectstorage/minio"
	"processor/internal/infrastructure/rabbitmq"
	"processor/internal/infrastructure/rabbitmq/consumers"
	"processor/internal/infrastructure/rabbitmq/producers"
	"processor/internal/infrastructure/repositories/cassandra"
	eventshandlers "processor/internal/interfaces/events/handlers"
	"processor/internal/usecases"
	"time"
)

type Application struct {
	Config        *config.Config
	Cassandra     *gocql.Session
	Producers     *producers.Producers
	Consumers     *consumers.Consumers
	UseCases      *usecases.UseCases
	EventHandlers *eventshandlers.EventHandlers
	NewRelicApp   *newrelic.Application
}

func NewApplication(cfg *config.Config) *Application {
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("chargeback-processor"),
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

	// Initialize file uploader
	chargebackUploader, err := minio.NewChargebackUploader(
		cfg.Minio.Endpoint,
		cfg.Minio.AccessKey,
		cfg.Minio.SecretKey,
		cfg.Minio.BucketName,
		cfg.Minio.UseSSL,
	)
	if err != nil {
		log.Fatalf("Could not initialize uploader: %v", err)
	}

	// Initialize repositories
	chargebackRepository := cassandra.NewChargebackRepositoryCassandra(cassandraSession)

	// Initialize writer
	chargebackWriter := filewriter.NewChargebackWriter("/tmp/chargebacks", cfg.Chargeback.MaxRecords, cfg.Chargeback.MaxDuration, chargebackUploader, producers.ChargebackBatchProducer)

	// Start background flush monitor for last files lost
	go func() {
		ticker := time.NewTicker(cfg.Chargeback.MaxDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				chargebackWriter.MaybeFlush()
			}
		}
	}()

	// Initializer use cases
	chargebackOpenedEventUseCase := usecases.NewChargebackOpenedEventUseCase(chargebackRepository, chargebackWriter)
	useCases := usecases.NewUseCases(chargebackOpenedEventUseCase)

	// Initializer specific handlers
	chargebackOpenedEventHandler := eventshandlers.NewChargebackOpenedEventHandler(chargebackOpenedEventUseCase)

	// Initialize event handlers
	eventHandlers := eventshandlers.NewEventHandlers(chargebackOpenedEventHandler)

	// Initialize consumers
	consumers, err := consumers.NewConsumers(rabbitConn, eventHandlers)
	if err != nil {
		log.Fatalf("Could not create RabbitMQ consumers: %v", err)
	}

	return &Application{
		Config:        cfg,
		Cassandra:     cassandraSession,
		Producers:     producers,
		Consumers:     consumers,
		UseCases:      useCases,
		EventHandlers: eventHandlers,
		NewRelicApp:   newRelicApp,
	}
}

func (app *Application) Close() {
	if app.Cassandra != nil {
		app.Cassandra.Close()
	}
	if app.Producers != nil {
		app.Producers.ChargebackBatchProducer.Connection.Close()
		// Close other producers as needed
	}
	// Close RabbitMQ connection
	if app.Consumers != nil {
		app.Consumers.ChargebackOpenedEventConsumer.Connection.Close()
		// Close other consumers as needed
	}
	if app.NewRelicApp != nil {
		app.NewRelicApp.Shutdown(10 * time.Second)
	}
}
