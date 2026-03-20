// services/notifications/module.go
package notifications

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	oglcore "github.com/ovya/ogl/platform/core"
	"github.com/rotisserie/eris"
)

type App struct {
	subscriber message.Subscriber
	logger     *slog.Logger
}

// Ensure Module implements oglcore.Module
var _ oglcore.Module = (*App)(nil)

func New(subscriber message.Subscriber, logger *slog.Logger) *App {
	return &App{
		subscriber: subscriber,
		logger:     logger,
	}
}

func (m *App) GetName() string {
	return "notification"
}

// StartWorkers boots up the background listeners for this module
// This mehod blocks until shutdown (ctx canceled)
func (m *App) Start(ctx context.Context) error {
	messages, err := m.subscriber.Subscribe(ctx, "todo.created")
	if err != nil {
		return eris.Wrap(err, "events subscription failed")
	}

	m.logger.Info("notification worker started listening for events")

	// This blocks the errgroup correctly until shutdown
	for msg := range messages {
		m.logger.Info(fmt.Sprintf("Notification Service: received event. Payload: %s\n", string(msg.Payload)))

		msg.Ack()
	}

	m.logger.Info("notification worker stopped gracefully")

	return nil
}

func (m *App) Close() error {
	m.logger.Info("shutting down notifications module internal resources")

	return nil
}
