package di

import (
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/settings"
	"github.com/didiegovieira/go-payments-core/pkg/log/implement"
	"github.com/google/wire"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func provideTracer() trace.Tracer {
	return otel.Tracer(settings.Settings.Metrics.Name)
}

func provideLogger() implement.Logger {
	return implement.NewLogrus()
}

var wireCommonSet = wire.NewSet(
	provideLogger,
	provideTracer,
	wire.Struct(new(application.App), "*"),
)
