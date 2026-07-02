package di

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler/payment"
	"github.com/google/wire"
)

var wireHandlerSet = wire.NewSet(
	wire.Struct(new(handler.Health), "*"),

	// Payment
	wire.Struct(new(payment.Create), "*"),

	// Account
	wire.Struct(new(account.Create), "*"),
)
