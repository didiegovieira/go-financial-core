package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/settings"
)

type Worker struct {
	App *App
}

func (w *Worker) Start() {
	w.App.Start(settings.Settings.Metrics.Name)
	w.App.Logger.Infof("Notification worker started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	w.App.Logger.Infof("Notification worker stopping")
	w.App.Stop()
}
