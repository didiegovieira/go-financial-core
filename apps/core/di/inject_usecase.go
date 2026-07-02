package di

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/payment"
	"github.com/google/wire"
)

var wireUsecaseSet = wire.NewSet(
	// Payment
	payment.NewCreate,
	wire.Bind(new(payment.CreateUseCase), new(*payment.Create)),

	// Account
	account.NewCreate,
	wire.Bind(new(account.CreateUseCase), new(*account.Create)),
)
