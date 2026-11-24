// Package cache cache package
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config contain the config of the redis client
type Config struct {
	Addr     string
	Password string
	DB       int
	Timeout  time.Duration
}

// DefaultConfig returns config to local development
func DefaultConfig() Config {
	return Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Timeout:  0 * time.Second,
	}
}

// NewClient create a new redis client
func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        cfg.Addr,
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: cfg.Timeout,
	})

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("verifying connection to redis: %w", err)
	}

	return client, nil
}
