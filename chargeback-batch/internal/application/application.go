package application

import (
	"batch/internal/config"
	"batch/internal/infrastructure/logging"
	"batch/internal/infrastructure/objectstorage"
	"batch/internal/infrastructure/objectstorage/minio"
	"batch/internal/infrastructure/transfer/ftp"

	//"batch/internal/infrastructure/objectstorage/minio"
	"batch/internal/infrastructure/rabbitmq"
	"batch/internal/infrastructure/rabbitmq/consumers"
	"batch/internal/infrastructure/repositories/postgres"
	eventshandlers "batch/internal/interfaces/events/handlers"
	"batch/internal/usecases"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"time"
)

type Application struct {
	Config              *config.Config
	DBConn              *sql.DB
	Consumers           *consumers.Consumers
	UseCases            *usecases.UseCases
	EventHandlers       *eventshandlers.EventHandlers
	NewRelicApp         *newrelic.Application
	FTPClient           ftp.Client
	BatchFileDownloader objectstorage.Downloader
}

func NewApplication(cfg *config.Config) *Application {
	var newRelicApp *newrelic.Application
	var err error

	if cfg.NewRelic.Enabled && cfg.NewRelic.LicenseKey != "" {
		newRelicApp, err = newrelic.NewApplication(
			newrelic.ConfigAppName("chargeback-api"),
			newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
			newrelic.ConfigDistributedTracerEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			log.Fatalf("Failed to create New Relic application: %v", err)
		}

		logging.InitializeLogger(newRelicApp)
		logging.Logger.Info("Application started with New Relic integration")

	} else {
		logging.InitializeLogger(nil)
		logging.Logger.Info("New Relic is disabled. Starting without New Relic integration")
	}

	// Initialize the global logger
	logging.InitializeLogger(newRelicApp)
	logging.Logger.Info("Application started with New Relic integration")

	// Connect to database
	dbConn, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Connect to RabbitMQ
	rabbitConn, err := rabbitmq.Connect(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}

	// Initialize producers

	// Initialize file downloader
	batchFileDownloader, err := minio.NewBatchFileDownloader(
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
	batchFileRepo := postgres.NewBatchFilesRepositoryPostgres(dbConn)

	// Initializer use cases
	chargebackBatchEventUseCase := usecases.NewChargebackBatchEventUseCase(batchFileRepo)
	useCases := usecases.NewUseCases(chargebackBatchEventUseCase)

	// Initializer specific handlers
	chargebackBatchEventHandler := eventshandlers.NewChargebackBatchEventHandler(chargebackBatchEventUseCase)

	// Initialize event handlers
	eventHandlers := eventshandlers.NewEventHandlers(chargebackBatchEventHandler)

	// Initialize consumers
	consumers, err := consumers.NewConsumers(rabbitConn, eventHandlers)
	if err != nil {
		log.Fatalf("Could not create RabbitMQ consumers: %v", err)
	}

	ftpClient, err := ftp.NewFTPClient(cfg.FTP.Host, cfg.FTP.Port, cfg.FTP.Username, cfg.FTP.Password)
	if err != nil {
		log.Fatalf("Could not initialize FTP client: %v", err)
	}

	return &Application{
		Config:              cfg,
		DBConn:              dbConn,
		Consumers:           consumers,
		UseCases:            useCases,
		EventHandlers:       eventHandlers,
		FTPClient:           ftpClient,
		BatchFileDownloader: batchFileDownloader,
		NewRelicApp:         newRelicApp,
	}
}

func (app *Application) Close() {
	if app.DBConn != nil {
		app.DBConn.Close()
	}
	// Close RabbitMQ connection
	if app.Consumers != nil {
		app.Consumers.ChargebackBatchEventConsumer.Connection.Close()
		// Close other consumers as needed
	}
	if app.NewRelicApp != nil {
		app.NewRelicApp.Shutdown(10 * time.Second)
	}
}
