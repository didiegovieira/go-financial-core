package di

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/wallet"
	"github.com/google/wire"
)

var wireUsecaseSet = wire.NewSet(
	// Account
	account.NewCreate,
	wire.Bind(new(account.CreateUseCase), new(*account.Create)),
	account.NewFindByPhone,
	wire.Bind(new(account.FindByPhoneUseCase), new(*account.FindByPhone)),
	account.NewDeleteByPhone,
	wire.Bind(new(account.DeleteByPhoneUseCase), new(*account.DeleteByPhone)),

	// Wallet
	wallet.NewCredit,
	wire.Bind(new(wallet.CreditUseCase), new(*wallet.Credit)),
	wallet.NewDebit,
	wire.Bind(new(wallet.DebitUseCase), new(*wallet.Debit)),
)
