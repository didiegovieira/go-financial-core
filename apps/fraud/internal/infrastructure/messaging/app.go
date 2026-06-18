package messaging

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/didiegovieira/go-payments-core/apps/fraud/internal/application"
	"github.com/didiegovieira/go-payments-core/apps/fraud/internal/settings"
)

type Application struct {
	BaseApp *application.App
}

func (a *Application) Start() {
	a.BaseApp.Start(settings.Settings.Metrics.Name)
	a.BaseApp.Logger.Infof("Notifications worker started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	a.BaseApp.Logger.Infof("Notifications worker stopping")
	a.BaseApp.Stop()
}
