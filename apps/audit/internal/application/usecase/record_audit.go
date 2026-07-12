package usecase

import (
	"context"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/dto"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
	"github.com/google/uuid"
)

type RecordAuditUseCase interface {
	Execute(ctx context.Context, input dto.RecordAuditInput) error
}

type RecordAudit struct {
	repository repository.AuditLogRepository
}

func NewRecordAudit(repository repository.AuditLogRepository) *RecordAudit {
	return &RecordAudit{
		repository: repository,
	}
}

func (r *RecordAudit) Execute(ctx context.Context, input dto.RecordAuditInput) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "UseCase.Audit.RecordAudit", &err)
	defer finish()

	return r.repository.Save(ctx, entity.AuditLog{
		ID:          uuid.New(),
		EventName:   input.EventName,
		AggregateID: input.AggregateID,
		Payload:     input.Payload,
		OccurredAt:  time.Now().UTC(),
	})
}
