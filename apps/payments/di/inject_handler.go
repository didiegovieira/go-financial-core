package di

import (
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api/handler"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api/handler/payment"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api/handler/wallet"
	"github.com/google/wire"
)

var wireHandlerSet = wire.NewSet(
	wire.Struct(new(handler.Health), "*"),

	// Payment
	wire.Struct(new(payment.Create), "*"),

	// Wallet
	wire.Struct(new(wallet.Create), "*"),
)
