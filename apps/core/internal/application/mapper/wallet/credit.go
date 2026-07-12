package wallet

import (
	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
)

func MapCreditInputToWallet(input dto.CreditInput, accountID uuid.UUID) *entity.Wallet {
	return &entity.Wallet{
		AccountID: accountID,
		Balance:   input.Amount,
	}
}

func MapWalletToCreditOutput(wallet *entity.Wallet) *dto.CreditOutput {
	return &dto.CreditOutput{
		Balance: wallet.Balance,
	}
}
