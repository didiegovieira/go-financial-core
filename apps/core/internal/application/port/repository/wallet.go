package repository

import (
	"context"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet entity.Wallet) error
}
