package di

import (
	"database/sql"

	appRepository "github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres"
	postgresRepository "github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/repository"
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

	// Account
	postgresRepository.NewAccountRepository,
	wire.Bind(new(appRepository.AccountRepository), new(*postgresRepository.AccountRepository)),

	// Wallet
	postgresRepository.NewWalletRepository,
	wire.Bind(new(appRepository.WalletRepository), new(*postgresRepository.WalletRepository)),
)
