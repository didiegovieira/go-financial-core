package di

import (
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application/usecase/payment"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application/usecase/wallet"
	"github.com/google/wire"
)

var wireUsecaseSet = wire.NewSet(
	// Payment
	payment.NewCreate,
	wire.Bind(new(payment.CreateUseCase), new(*payment.Create)),

	// Wallet
	wallet.NewCreate,
	wire.Bind(new(wallet.CreateUseCase), new(*wallet.Create)),
)
