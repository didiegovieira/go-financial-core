package di

import (
	"database/sql"

	apprepository "github.com/didiegovieira/go-financial-core/apps/audit/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/infrastructure/database/postgres"
	postgresrepository "github.com/didiegovieira/go-financial-core/apps/audit/internal/infrastructure/database/postgres/repository"
	"github.com/google/wire"
)

func provideDatabase() *sql.DB {
	db, err := postgres.NewConnection()
	if err != nil {
		panic(err)
	}
	return db
}

var wireRepositorySet = wire.NewSet(
	provideDatabase,
	postgresrepository.NewAuditLogRepository,
	wire.Bind(new(apprepository.AuditLogRepository), new(*postgresrepository.AuditLogRepository)),
)
