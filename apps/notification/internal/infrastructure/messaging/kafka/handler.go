package kafka

import (
	"context"
	"encoding/json"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/dto"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/usecase"
	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type Handler struct {
	sendNotification usecase.SendNotificationUseCase
}

func NewHandler(sendNotification usecase.SendNotificationUseCase) *Handler {
	return &Handler{
		sendNotification: sendNotification,
	}
}

func (h *Handler) Handle(ctx context.Context, message sharedkafka.Message) (err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Handler.Kafka.Notification.Handle", &err)
	defer finish()

	eventName := eventNameByTopic[message.Topic]
	if eventName == "" {
		return sharedkafka.ErrUnknownEventName
	}

	return h.sendNotification.Execute(ctx, dto.SendNotificationInput{
		EventName:   eventName,
		AggregateID: message.Key,
		Payload:     json.RawMessage(message.Value),
	})
}
