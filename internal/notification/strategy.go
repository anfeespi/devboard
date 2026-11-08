// Package notification notification
package notification

import (
	"context"
	"log/slog"
)

// Notifier interface for the notifications
type Notifier interface {
	Notify(ctx context.Context, userID, message string) error
}

// LogNotifier registry the notification in the logs
type LogNotifier struct {
	logger *slog.Logger
}

// NewLogNotifier constructor
func NewLogNotifier(logger *slog.Logger) *LogNotifier {
	return &LogNotifier{
		logger: logger,
	}
}

// Notify notify method from logs
func (n *LogNotifier) Notify(ctx context.Context, userID, message string) error {
	n.logger.InfoContext(ctx, "notification sended",
		slog.String("user_id", userID),
		slog.String("message", message),
	)

	return nil
}

// NoOpNotifier dont notify
type NoOpNotifier struct{}

// Notify dont notify
func (n *NoOpNotifier) Notify(_ context.Context, _, _ string) error {
	return nil
}
