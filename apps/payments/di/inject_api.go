package di

import (
	"net/http"

	"github.com/didiegovieira/go-payments-core/apps/payments/internal/settings"
	"github.com/didiegovieira/go-payments-core/pkg/api"
	"github.com/didiegovieira/go-payments-core/pkg/api/presenter"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func provideHttpServer() *http.Server {
	return &http.Server{
		Addr:         settings.Settings.HttpServer.Port,
		ReadTimeout:  settings.Settings.HttpServer.ReadTimeout,
		WriteTimeout: settings.Settings.HttpServer.WriteTimeout,
	}
}

func provideGinServer(httpServer *http.Server) api.Server[*gin.Engine] {
	return api.NewGinServer(httpServer)
}

func providePresenter() api.Presenter {
	return presenter.NewJson()
}

var wireServerSet = wire.NewSet(
	provideHttpServer,
	provideGinServer,
	providePresenter,
)
