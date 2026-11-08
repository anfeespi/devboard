// Package logger logger
package logger

import (
	"io"
	"log/slog"
	"os"
)

// Config struct configuration
type Config struct {
	Level  slog.Level
	JSON   bool
	Output io.Writer
}

// DefaultConfig configuration default
func DefaultConfig() Config {
	return Config{
		Level:  slog.LevelDebug,
		JSON:   false,
		Output: os.Stdout,
	}
}

// ProductionConfig configuration productive
func ProductionConfig() Config {
	return Config{
		Level:  slog.LevelInfo,
		JSON:   true,
		Output: os.Stdout,
	}
}

// New config
func New(config Config) *slog.Logger {
	out := config.Output
	if out == nil {
		out = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.Level == slog.LevelDebug,
	}

	var handler slog.Handler

	if config.JSON {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = slog.NewTextHandler(out, opts)
	}

	return slog.New(handler)
}
