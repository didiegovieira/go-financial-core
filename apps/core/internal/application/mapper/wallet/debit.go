package wallet

import (
	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
)

func MapDebitInputToWallet(input dto.DebitInput, accountID uuid.UUID) *entity.Wallet {
	return &entity.Wallet{
		AccountID: accountID,
		Balance:   input.Amount,
	}
}

func MapWalletToDebitOutput(wallet *entity.Wallet) *dto.DebitOutput {
	return &dto.DebitOutput{
		Balance: wallet.Balance,
	}
}
