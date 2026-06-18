package di

import (
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api/handler"
	"github.com/google/wire"
)

var wireHandlerSet = wire.NewSet(
	wire.Struct(new(handler.Health), "*"),
)
