package application

import (
	"github.com/didiegovieira/go-payments-core/apps/notifications/internal/settings"
	"github.com/didiegovieira/go-payments-core/pkg/log"
	"github.com/didiegovieira/go-payments-core/pkg/log/implement"
	"github.com/didiegovieira/go-payments-core/pkg/metrics"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	Logger implement.Logger
	Tracer trace.Tracer
}

func (a *App) Start(serviceName string) {
	log.Logger = a.Logger

	if settings.Settings.IsLocal() {
		metrics.Tracer = trace.NewNoopTracerProvider().Tracer("local")
	} else {
		metrics.Tracer = a.Tracer
	}

	a.Logger.Tag("service", serviceName)
}

func (a *App) Stop() {
	a.Logger.Infof("Application stopping")
}
