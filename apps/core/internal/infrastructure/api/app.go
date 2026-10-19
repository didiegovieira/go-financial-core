package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/application"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/handler/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/api/middleware"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	"github.com/didiegovieira/go-financial-core/pkg/api"
	"github.com/gin-gonic/gin"
)

type Application struct {
	BaseApp *application.App
	Server  api.Server[*gin.Engine]

	// Middleware
	CorsMiddleware *middleware.Cors

	// Health
	HealthHandler *handler.Health

	// Account
	AccountCreateHandler *account.Create
	FindByPhoneHandler   *account.FindByPhone
	DeleteByPhoneHandler *account.DeleteByPhone

	// Wallet
	WalletCreditHandler *wallet.Credit
	WalletDebitHandler  *wallet.Debit
}

func (a *Application) Start() {
	a.BaseApp.Start(settings.Settings.Telemetry.Name)

	a.SetupRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.Server.Shutdown(shutdownCtx); err != nil {
			a.BaseApp.Logger.Infof("Server forced to shutdown: %v", err)
		}
	}()

	if err := a.Server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.BaseApp.Logger.Infof("Failed to start server: %v", err)
	}

	a.BaseApp.Logger.Infof("Server exited properly")
	a.BaseApp.Stop()
}
