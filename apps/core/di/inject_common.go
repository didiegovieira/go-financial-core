package di

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	"github.com/didiegovieira/go-financial-core/pkg/log/implement"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/trace"
)

func provideTracer() trace.Tracer {
	return telemetry.NewTracer(settings.Settings.Telemetry.Name)
}

func provideLogger() implement.Logger {
	return implement.NewLogrus()
}

var wireCommonSet = wire.NewSet(
	provideLogger,
	provideTracer,
	wire.Struct(new(application.App), "*"),
)
