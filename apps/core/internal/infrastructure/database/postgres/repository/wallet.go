package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	domainerror "github.com/didiegovieira/go-financial-core/apps/core/internal/domain/error"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/mapper"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/query"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
	"github.com/google/uuid"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet entity.Wallet) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Repository.Wallet.Create", &err)
	defer finish()

	walletModel := mapper.WalletToModel(wallet)

	_, err = r.db.ExecContext(
		ctx,
		query.WalletCreate,
		walletModel.ID,
		walletModel.AccountID,
		walletModel.Balance,
		walletModel.Status,
	)
	return err
}

func (r *WalletRepository) DeleteByAccountID(ctx context.Context, accountID uuid.UUID) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Repository.Wallet.DeleteByAccountID", &err)
	defer finish()

	_, err = r.db.ExecContext(
		ctx,
		query.WalletDeleteByAccountID,
		accountID,
	)
	return err
}

func (r *WalletRepository) Credit(ctx context.Context, wallet entity.Wallet) (updatedWallet *entity.Wallet, err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Repository.Wallet.Credit", &err)
	defer finish()

	walletModel := mapper.WalletToModel(wallet)

	err = r.db.QueryRowContext(
		ctx,
		query.WalletCredit,
		walletModel.Balance,
		walletModel.AccountID,
	).Scan(&walletModel.Balance)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domainerror.ErrWalletUnavailable
	}
	if err != nil {
		return nil, err
	}

	updated := mapper.WalletToEntity(walletModel)
	return &updated, nil
}

func (r *WalletRepository) Debit(ctx context.Context, wallet entity.Wallet) (updatedWallet *entity.Wallet, err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Repository.Wallet.Debit", &err)
	defer finish()

	walletModel := mapper.WalletToModel(wallet)

	err = r.db.QueryRowContext(
		ctx,
		query.WalletDebit,
		walletModel.Balance,
		walletModel.AccountID,
	).Scan(&walletModel.Balance)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domainerror.ErrDebitNotAllowed
	}
	if err != nil {
		return nil, err
	}

	updated := mapper.WalletToEntity(walletModel)
	return &updated, nil
}
