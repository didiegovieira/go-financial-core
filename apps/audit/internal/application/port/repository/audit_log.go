package repository

import (
	"context"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/domain/entity"
)

type AuditLogRepository interface {
	Save(ctx context.Context, auditLog entity.AuditLog) error
}
