package client

import (
	"context"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/domain/entity"
)

type NotificationClient interface {
	Send(ctx context.Context, notification entity.Notification) error
}
