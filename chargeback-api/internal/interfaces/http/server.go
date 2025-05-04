package http

import (
	"api/internal/application"
	"api/internal/config"
	"api/internal/infrastructure/middlewares"
	"api/internal/interfaces/http/handlers"
	"api/internal/interfaces/http/routes"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg        *config.Config
	handlers   *handlers.Handlers
	httpServer *http.Server
}

func NewServer(application *application.Application, handlers *handlers.Handlers) *Server {
	router := chi.NewRouter()
	cfg := application.Config

	// Add New Relic middlewares
	router.Use(middlewares.NewRelicMiddleware(application.NewRelicApp))

	// Global middlewares
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(cfg.Server.ReadTimeout))

	// Swagger configuration
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Routes configuration
	routes.SetupRoutes(router, handlers)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &Server{
		cfg:        cfg,
		handlers:   handlers,
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}
