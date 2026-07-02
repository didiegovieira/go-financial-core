package di

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/middleware"
	"github.com/google/wire"
)

var wireMiddlewareSet = wire.NewSet(
	wire.Struct(new(middleware.Cors), "*"),
)
