package runtime

import (
	"github.com/didiegovieira/go-financial-core/pkg/log"
	"github.com/didiegovieira/go-financial-core/pkg/log/implement"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	Logger implement.Logger
	Tracer trace.Tracer
}

func (a *App) Start(serviceName string) {
	log.Logger = a.Logger
	telemetry.Tracer = a.Tracer

	a.Logger.Tag("service", serviceName)
}

func (a *App) Stop() {
	a.Logger.Infof("Application stopping")
}
