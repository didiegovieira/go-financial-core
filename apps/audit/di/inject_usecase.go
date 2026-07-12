package di

import (
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/usecase"
	"github.com/google/wire"
)

var wireUsecaseSet = wire.NewSet(
	usecase.NewRecordAudit,
	wire.Bind(new(usecase.RecordAuditUseCase), new(*usecase.RecordAudit)),
)
