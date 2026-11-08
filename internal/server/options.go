package server

import (
	"log/slog"
	"time"
)

type config struct {
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
	logger            *slog.Logger
}

// defaultConfig return default values
func defaultConfig() config {
	return config{
		readTimeout:       10 * time.Second,
		writeTimeout:      30 * time.Second,
		idleTimeout:       60 * time.Second,
		readHeaderTimeout: 5 * time.Second,
		logger:            slog.Default(),
	}
}

// Option is a func that change the configuration
type Option func(*config)

// WithReadTimeout configure the timeout of read
func WithReadTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.readTimeout = duration
	}
}

// WithWriteTimeout configure the timeout of write
func WithWriteTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.writeTimeout = duration
	}
}

// WithIdleTimeout configure the timeout of idle
func WithIdleTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.idleTimeout = duration
	}
}

// WithReadHeaderTimeout configure the timeout of read header
func WithReadHeaderTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.readHeaderTimeout = duration
	}
}

// WithLogger configure the logger
func WithLogger(logger *slog.Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}
