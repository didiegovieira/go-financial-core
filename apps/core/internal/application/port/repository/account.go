package repository

import (
	"context"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
)

type AccountRepository interface {
	Create(ctx context.Context, account entity.Account) error
	FindByDocument(ctx context.Context, account entity.Account) (*entity.Account, error)
	FindByPhone(ctx context.Context, account entity.Account) (*entity.Account, error)
	DeleteByPhone(ctx context.Context, account entity.Account) error
}
