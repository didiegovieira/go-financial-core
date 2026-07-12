package repository

import (
	"context"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet entity.Wallet) error
	DeleteByAccountID(ctx context.Context, accountID uuid.UUID) error
	Credit(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error)
	Debit(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error)
}
