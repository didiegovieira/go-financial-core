package di

import (
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/runtime"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/settings"
	"github.com/didiegovieira/go-financial-core/pkg/log/implement"
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
	wire.Struct(new(runtime.App), "*"),
)
