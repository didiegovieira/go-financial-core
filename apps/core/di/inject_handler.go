package di

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler/wallet"
	"github.com/google/wire"
)

var wireHandlerSet = wire.NewSet(
	wire.Struct(new(handler.Health), "*"),

	// Account
	wire.Struct(new(account.Create), "*"),
	wire.Struct(new(account.FindByPhone), "*"),
	wire.Struct(new(account.DeleteByPhone), "*"),

	// Wallet
	wire.Struct(new(wallet.Credit), "*"),
	wire.Struct(new(wallet.Debit), "*"),
)
