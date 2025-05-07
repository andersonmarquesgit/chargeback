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

	go scheduler.StartScheduler(
		app.UseCases.ChargebackBatchEventUseCase.BatchFilesRepository,
		app.BatchFileDownloader,
		app.FTPClient,
	)

	log.Println("Batch started and listening for batch events...")

	// Wait for termination signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down batch...")
}
