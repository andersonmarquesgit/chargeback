package main

import (
	"batch/internal/application"
	"batch/internal/config"
	"batch/internal/scheduler"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize application
	app := application.NewApplication(cfg)
	defer app.Close()

	if cfg.Scheduler.Enabled {
		log.Printf("Starting scheduler with interval: %v", cfg.Scheduler.Interval)
		go scheduler.StartScheduler(
			app.UseCases.ChargebackBatchEventUseCase.BatchFilesRepository,
			app.BatchFileDownloader,
			app.FTPClient,
			cfg.Scheduler.Interval,
		)
	} else {
		log.Println("Scheduler is disabled. Worker running without scheduler")
	}

	log.Println("Batch started and listening for batch events...")

	// Wait for termination signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down batch...")
}
