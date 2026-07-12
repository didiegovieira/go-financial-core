package usecase

import (
	"context"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/dto"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/port/client"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
	"github.com/google/uuid"
)

type SendNotificationUseCase interface {
	Execute(ctx context.Context, input dto.SendNotificationInput) error
}

type SendNotification struct {
	client client.NotificationClient
}

func NewSendNotification(client client.NotificationClient) *SendNotification {
	return &SendNotification{
		client: client,
	}
}

func (s *SendNotification) Execute(ctx context.Context, input dto.SendNotificationInput) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "UseCase.Notification.SendNotification", &err)
	defer finish()

	notification := entity.Notification{
		ID:          uuid.NewString(),
		EventName:   input.EventName,
		AggregateID: input.AggregateID,
		Channel:     "external",
		Recipient:   input.AggregateID,
		Subject:     input.EventName,
		Body:        string(input.Payload),
		Status:      "pending",
		CreatedAt:   time.Now().UTC(),
	}

	return s.client.Send(ctx, notification)
}
