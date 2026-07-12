package mapper

import (
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/infrastructure/database/postgres/model"
)

func AuditLogToModel(auditLog entity.AuditLog) model.AuditLogModel {
	return model.AuditLogModel{
		ID:          auditLog.ID,
		EventName:   auditLog.EventName,
		AggregateID: auditLog.AggregateID,
		Payload:     auditLog.Payload,
		OccurredAt:  auditLog.OccurredAt,
	}
}
