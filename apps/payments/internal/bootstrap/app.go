package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didiegovieira/go-payments-core/apps/payments/internal/config"
	presentationhttp "github.com/didiegovieira/go-payments-core/apps/payments/internal/presentation/http"
)

type App struct {
	settings config.Settings
	logger   *slog.Logger
	server   *http.Server
}

func NewApp(settings config.Settings, logger *slog.Logger, server *http.Server) *App {
	return &App{
		settings: settings,
		logger:   logger,
		server:   server,
	}
}

func NewLogger(settings config.Settings) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel(settings.App.Environment),
	}))
}

func NewRouterConfig(settings config.Settings) presentationhttp.RouterConfig {
	return presentationhttp.RouterConfig{
		AppName:     settings.App.Name,
		Environment: settings.App.Environment,
	}
}

func NewHTTPServer(settings config.Settings, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              settings.HTTP.Address(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		a.logger.Info(
			"starting payments api",
			"address", a.server.Addr,
			"environment", a.settings.App.Environment,
		)

		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return a.shutdown()
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}
}

func (a *App) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.settings.Shutdown.Timeout)
	defer cancel()

	a.logger.Info("shutting down payments api")

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.logger.Info("payments api stopped")
	return nil
}

func logLevel(environment string) slog.Level {
	if environment == "production" {
		return slog.LevelInfo
	}

	return slog.LevelDebug
}
