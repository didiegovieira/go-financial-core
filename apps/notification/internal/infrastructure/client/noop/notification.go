package noop

import (
	"context"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/pkg/log/implement"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type NotificationClient struct {
	logger implement.Logger
}

func NewNotificationClient(logger implement.Logger) *NotificationClient {
	return &NotificationClient{
		logger: logger,
	}
}

func (c *NotificationClient) Send(ctx context.Context, notification entity.Notification) (err error) {
	_, finish := telemetry.TraceWithError(ctx, "Client.Notification.Noop.Send", &err)
	defer finish()

	c.logger.Infof(
		"Notification sent event=%s aggregate_id=%s id=%s channel=%s",
		notification.EventName,
		notification.AggregateID,
		notification.ID,
		notification.Channel,
	)
	return nil
}
