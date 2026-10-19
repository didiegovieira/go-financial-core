package postgres

import (
	"context"
	"database/sql"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/infrastructure/database/postgres/mapper"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/infrastructure/database/postgres/query"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Save(ctx context.Context, auditLog entity.AuditLog) (err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Repository.Audit.AuditLog.Save", &err)
	defer finish()

	auditLogModel := mapper.AuditLogToModel(auditLog)

	_, err = r.db.ExecContext(
		ctx,
		query.AuditLogCreate,
		auditLogModel.ID,
		auditLogModel.EventName,
		auditLogModel.AggregateID,
		auditLogModel.Payload,
		auditLogModel.OccurredAt,
	)
	return err
}
