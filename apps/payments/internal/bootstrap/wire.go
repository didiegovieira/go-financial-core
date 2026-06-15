//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/config"
	presentationhttp "github.com/didiegovieira/go-payments-core/apps/payments/internal/presentation/http"
	"github.com/google/wire"
)

// InitializeApp composes the payments API dependency graph.
func InitializeApp(settings config.Settings) *App {
	wire.Build(
		NewLogger,
		NewRouterConfig,
		presentationhttp.NewRouter,
		NewHTTPServer,
		NewApp,
	)

	return nil
}
