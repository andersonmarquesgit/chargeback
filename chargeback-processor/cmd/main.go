package main

import (
	"log"
	"os"
	"os/signal"
	"processor/internal/application"
	"processor/internal/config"
	"syscall"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize application
	app := application.NewApplication(cfg)
	defer app.Close()

	log.Println("Processor started and listening for chargeback events...")

	// Wait for termination signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down processor...")
}
