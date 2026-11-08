// Package main
package main

import (
	"log/slog"
	"os"
	"time"

	_ "github.com/anfeespi/devboard/docs"
	"github.com/anfeespi/devboard/internal/handler"
	"github.com/anfeespi/devboard/internal/logger"
	"github.com/anfeespi/devboard/internal/middlewares"
	"github.com/anfeespi/devboard/internal/notification"
	"github.com/anfeespi/devboard/internal/server"
	"github.com/anfeespi/devboard/internal/validator"
	httpSwagger "github.com/swaggo/http-swagger"
)

const PORT = ":8080"

// @title Devboard API
// @version 0.1
// @description API REST para gestion de proyectos
// @contact.name Soporte Devboard Andres
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Config of the environment
	env := os.Getenv("GO_ENV")
	slog.Info("env charged", "GO_ENV", env)

	// 2. Logger
	var log *slog.Logger
	if env == "production" {
		log = logger.New(logger.ProductionConfig())
	} else {
		log = logger.New(logger.DefaultConfig())
	}

	// 3. Shared dependencies
	validate := validator.New()
	_ = validate
	notifier := notification.NewLogNotifier(log)
	_ = notifier

	// 4. Server with Functional Options
	srv := server.New(PORT,
		server.WithLogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(40*time.Second),
	)

	// 5. Handlers
	healthHandler := handler.NewHealthHandler()

	// 6. Register routes
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

	// 7. Middleware chain
	srv.UseChain(
		middlewares.Recovery(log),
		middlewares.Logger(log),
	)

	// 8. Start the app
	if err := srv.Start(); err != nil {
		log.Error("error fatal", "error", err)
		os.Exit(1)
	}

}
