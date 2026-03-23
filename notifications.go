// services/notifications/module.go
package notifications

import (
	"context"
	"log/slog"
	"os"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/rocketchat"

	"github.com/ThreeDotsLabs/watermill/message"
	oglcore "github.com/ovya/ogl/platform/core"
	"github.com/rotisserie/eris"
	"golang.org/x/sync/errgroup"
)

const AppName = "Notifications"

var notifier *notify.Notify

type App struct {
	subscriber message.Subscriber
	topics     []string
	logger     *slog.Logger
}

// Ensure Module implements oglcore.Module
var _ oglcore.App = (*App)(nil)

func New(subscriber message.Subscriber, logger *slog.Logger, topics ...string) *App {
	return &App{
		subscriber: subscriber,
		topics:     topics,
		logger:     logger,
	}
}

// TODO: make a DI service for that.
func initRocketNotifier() error {
	host := os.Getenv("ROCKET_HOST")
	if host == "" {
		return nil
	}

	user := os.Getenv("ROCKET_USER")
	key := os.Getenv("ROCKET_KEY")
	rocketChatSvc, err := rocketchat.New(host, "https", user, key)
	if err != nil {
		return eris.Wrap(err, "failed to connect to rocketchat")
	}
	rocketChatSvc.AddReceivers("POC_CHANNEL")
	notifier = notify.New()
	// notify.UseServices(rocketChatSvc)
	notifier.UseServices(rocketChatSvc)
	err = notifier.Send(context.Background(), "Hi", "I am the bot of the mmw events' dispatcher.")

	return eris.Wrap(err, "failed to send message to rocketchat")
}

// Start boots up one listener per topic and blocks until shutdown (ctx canceled).
func (m *App) Start(ctx context.Context) error {
	// For the demostration to the OVYA team
	if err := initRocketNotifier(); err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, topic := range m.topics {
		messages, err := m.subscriber.Subscribe(gCtx, topic)
		if err != nil {
			return eris.Wrapf(err, "subscription to %q failed", topic)
		}

		m.logger.Info("notification worker listening", "topic", topic)

		g.Go(func() error {
			for msg := range messages {
				payloadStr := string(msg.Payload)
				m.logger.Info("event received",
					"topic", topic,
					"payload", payloadStr,
				)
				if notifier != nil {
					_ = notifier.Send(context.Background(), "event received > "+topic, payloadStr)
				}

				msg.Ack()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return eris.Wrap(err, "notification worker error")
	}

	m.logger.Info("notification worker stopped gracefully")

	return nil
}

func (m *App) Close() error {
	m.logger.Info("shutting down notifications module internal resources")

	return nil
}

func SendToRocket() {

}
