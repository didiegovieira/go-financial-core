package messaging

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/didiegovieira/go-payments-core/apps/audit/internal/application"
	"github.com/didiegovieira/go-payments-core/apps/audit/internal/settings"
)

type Application struct {
	BaseApp  *application.App
	Consumer *Consumer // seu Kafka consumer
}

func (a *Application) Start() {
	a.BaseApp.Start(settings.Settings.Metrics.Name)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		a.Consumer.Stop()
	}()

	a.Consumer.Start()
	a.BaseApp.Logger.Println("Worker exited properly")
}

func (a *Application) Stop() {
	a.Consumer.Stop()
	a.BaseApp.Stop()
}
