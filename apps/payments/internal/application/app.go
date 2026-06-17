package application

import (
	"log"

	"go.opentelemetry.io/otel/trace"
)

type App struct {
	Logger log.Logger
	Tracer trace.Tracer
}

func (a *App) Start(serviceName string) {
	// log2.Logger = a.Logger

	// if settings.Settings.IsLocal() {
	// 	metrics2.Tracer = trace.NewNoopTracerProvider().Tracer("local")
	// } else {
	// 	metrics2.Tracer = a.Tracer
	// }

	// a.Logger.Tag("service", serviceName)
}

func (a *App) Stop() {
	a.Logger.Println("Application stopping")
}
