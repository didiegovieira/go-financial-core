package mapper

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/model"
)

func WalletToEntity(w model.WalletModel) entity.Wallet {
	return entity.Wallet{
		ID:        w.ID,
		AccountID: w.AccountID,
		Balance:   w.Balance,
		Status:    w.Status,
	}
}

func WalletToModel(w entity.Wallet) model.WalletModel {
	return model.WalletModel{
		ID:        w.ID,
		AccountID: w.AccountID,
		Balance:   w.Balance,
		Status:    w.Status,
	}
}
