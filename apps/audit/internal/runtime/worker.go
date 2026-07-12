package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/settings"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type Worker struct {
	App      *App
	Consumer Consumer
}

func (w *Worker) Start() {
	w.App.Start(settings.Settings.Metrics.Name)
	w.App.Logger.Infof("Audit worker started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Consumer.Start(ctx)
	}()

	select {
	case <-ctx.Done():
		w.App.Logger.Infof("Audit worker stopping")
	case err := <-errCh:
		if err != nil {
			w.App.Logger.Errorf("Audit worker stopped with error: %v", err)
		}
	}

	if err := w.Consumer.Stop(); err != nil {
		w.App.Logger.Errorf("Failed to stop audit consumer: %v", err)
	}
	w.App.Stop()
}
