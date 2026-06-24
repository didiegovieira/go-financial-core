package di

import (
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api/middleware"
	"github.com/google/wire"
)

var wireMiddlewareSet = wire.NewSet(
	wire.Struct(new(middleware.Cors), "*"),
)
