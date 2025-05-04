package main

import (
	"api/internal/application"
	"api/internal/config"
	httpserver "api/internal/interfaces/http"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize application
	application := application.NewApplication(cfg)
	defer application.Close()

	// Initialize and start HTTP server
	server := httpserver.NewServer(application, application.Handlers)
	log.Printf("Starting chargeback api on port %s", cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
