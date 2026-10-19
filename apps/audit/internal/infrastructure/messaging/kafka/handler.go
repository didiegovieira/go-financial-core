package kafka

import (
	"context"
	"encoding/json"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/dto"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/usecase"
	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type Handler struct {
	recordAudit usecase.RecordAuditUseCase
}

func NewHandler(recordAudit usecase.RecordAuditUseCase) *Handler {
	return &Handler{
		recordAudit: recordAudit,
	}
}

func (h *Handler) Handle(ctx context.Context, message sharedkafka.Message) (err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Handler.Kafka.Audit.Handle", &err)
	defer finish()

	eventName := eventNameByTopic[message.Topic]
	if eventName == "" {
		return sharedkafka.ErrUnknownEventName
	}

	return h.recordAudit.Execute(ctx, dto.RecordAuditInput{
		EventName:   eventName,
		AggregateID: message.Key,
		Payload:     json.RawMessage(message.Value),
	})
}
