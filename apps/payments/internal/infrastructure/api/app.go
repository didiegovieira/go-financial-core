package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api/handler"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/settings"
	"github.com/didiegovieira/go-payments-core/pkg/api"
	"github.com/gin-gonic/gin"
)

type Application struct {
	BaseApp *application.App
	Server  api.Server[*gin.Engine]

	// Health
	HealthHandler *handler.Health
}

func (a *Application) Start() {
	a.BaseApp.Start(settings.Settings.Metrics.Name)

	a.SetupRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.Server.Shutdown(shutdownCtx); err != nil {
			a.BaseApp.Logger.Printf("Server forced to shutdown: %v", err)
		}
	}()

	if err := a.Server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.BaseApp.Logger.Printf("Failed to start server: %v", err)
	}

	a.BaseApp.Logger.Println("Server exited properly")
	a.BaseApp.Stop()
}
