// Package main
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/anfeespi/devboard/docs"
	"github.com/anfeespi/devboard/internal/cache"
	"github.com/anfeespi/devboard/internal/handler"
	"github.com/anfeespi/devboard/internal/logger"
	"github.com/anfeespi/devboard/internal/middlewares"
	"github.com/anfeespi/devboard/internal/notification"
	pgrepo "github.com/anfeespi/devboard/internal/repository/postgres"
	"github.com/anfeespi/devboard/internal/server"
	"github.com/anfeespi/devboard/internal/usecase"
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
	// 1. Environment configuration
	// 2. Logger
	log := setupLogger()

	dbURL := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")

	// Create connection pool
	pool, err := pgrepo.NewPool(context.Background(), pgrepo.Config{
		URL:             dbURL,
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	})
	if err != nil {
		log.Error("cannot connect to the bd", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	log.Info("connected to the db!")

	redisClient, err := cache.NewClient(context.Background(), cache.Config{
		Addr:    redisAddr,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("cannot connect to redus", slog.Any("error", err))
		os.Exit(1)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error(
				"cannot close the connection to redis",
				slog.Any("error", err),
			)
		}
	}()

	log.Info("connected to Redis!")

	// 3. Shared dependencies

	// Infrastructure: adapters
	userRepo := pgrepo.NewUserRepository(pool)
	taskRepo := pgrepo.NewTaskRepository(pool)

	// Use cases: business logic
	userUC := usecase.NewUserUseCase(userRepo)
	taskTC := usecase.NewTaskUseCase(taskRepo, taskRepo)

	// Shared validation
	validate := validator.New()

	// 4. Handlers
	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(userUC, validate, log)
	taskHandler := handler.NewTaskHandler(taskTC, validate, log)

	notifier := notification.NewLogNotifier(log)
	_ = notifier

	// 5. Server with Functional Options
	srv := server.New(PORT,
		server.WithLogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(40*time.Second),
	)

	// 6. Register routes
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))
	srv.RegisterRoutes("POST /api/v1/tasks", http.HandlerFunc(taskHandler.Create))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}", http.HandlerFunc(taskHandler.Get))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/status", http.HandlerFunc(taskHandler.UpdateStatus))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/assign", http.HandlerFunc(taskHandler.Assign))
	srv.RegisterRoutes("GET /api/v1/projects/{id}/tasks", http.HandlerFunc(taskHandler.ListByProject))

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

func setupLogger() *slog.Logger {
	if os.Getenv("GO_ENV") == "production" {
		return logger.New(logger.ProductionConfig())
	}

	return logger.New(logger.DefaultConfig())
}
