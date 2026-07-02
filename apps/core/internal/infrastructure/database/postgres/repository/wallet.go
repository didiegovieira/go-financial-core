package postgres

import (
	"context"
	"database/sql"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/mapper"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/query"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet entity.Wallet) error {
	walletModel := mapper.WalletToModel(wallet)

	_, err := r.db.ExecContext(
		ctx,
		query.WalletCreate,
		walletModel.ID,
		walletModel.AccountID,
		walletModel.Balance,
		walletModel.Status,
	)
	return err
}
