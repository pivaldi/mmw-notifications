// services/notifications/module.go
package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ovya/ogl/oglcore"
	"github.com/rotisserie/eris"
)

type Module struct {
	subscriber message.Subscriber
	logger     *slog.Logger
}

// Ensure Module implements oglcore.Module
var _ oglcore.Module = (*Module)(nil)

func (m *Module) RegisterRoutes(_ *http.ServeMux) {}

func Build(subscriber message.Subscriber, logger *slog.Logger) *Module {
	return &Module{
		subscriber: subscriber,
		logger:     logger,
	}
}

// StartWorkers boots up the background listeners for this module
// This mehod blocks until shutdown (ctx canceled)
func (m *Module) StartWorkers(ctx context.Context) error {
	messages, err := m.subscriber.Subscribe(ctx, "todo.created")
	if err != nil {
		return eris.Wrap(err, "events subscription failed")
	}

	m.logger.Info("notification worker started listening for events")

	// This blocks the errgroup correctly until shutdown
	for msg := range messages {
		m.logger.Info(fmt.Sprintf("Notification Service: received event new Todo! Payload: %s\n", string(msg.Payload)))

		msg.Ack()
	}

	m.logger.Info("notification worker stopped gracefully")

	return nil
}

func (m *Module) Close() error {
	m.logger.Info("shutting down notifications module internal resources")

	return nil
}
