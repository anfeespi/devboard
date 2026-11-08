// Package server creates and start server
package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anfeespi/devboard/internal/middlewares"
)

// Server struct
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	logger     *slog.Logger
}

// New creates a new server
func New(addr string, opts ...Option) *Server {
	cfg := defaultConfig()

	for _, o := range opts {
		o(&cfg)
	}

	mux := http.NewServeMux() // Enrutador http

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       cfg.readTimeout,
		WriteTimeout:      cfg.writeTimeout,
		IdleTimeout:       cfg.idleTimeout,
		ReadHeaderTimeout: cfg.writeTimeout,
	}

	return &Server{
		httpServer: httpServer,
		mux:        mux,
		logger:     cfg.logger,
	}
}

// Start server
func (s *Server) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(quit)

	serveErr := make(chan error, 1)
	go func() {
		s.logger.Info("Server initialized", slog.String("addr", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		return err
	case sig := <-quit:
		s.logger.Info("Signal received, shutting down", slog.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("error shutting down", slog.Any("error", err))
		return err
	}

	s.logger.Info("Server stopped")

	return nil
}

// RegisterRoutes register routes
func (s *Server) RegisterRoutes(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// ServeHTTP server
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Use use middleware (point the original handler to the middleware)
func (s *Server) Use(middleware func(http.Handler) http.Handler) {
	s.httpServer.Handler = middleware(s.httpServer.Handler)
}

// UseChain use the chain of middlewares
func (s *Server) UseChain(mws ...middlewares.Middleware) {
	s.httpServer.Handler = middlewares.Chain(s.httpServer.Handler, mws...)
}
