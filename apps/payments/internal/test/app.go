package test

import (
	"net/http/httptest"

	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/infrastructure/api"
	"go.uber.org/mock/gomock"
)

type Application struct {
	BaseApp *application.App
	Api     *api.Application

	MockCtrl *gomock.Controller

	ApiUrl    string           `wire:"-"`
	ApiServer *httptest.Server `wire:"-"`
}

func (a *Application) RunApiServer() *httptest.Server {
	a.BaseApp.Start("test-api-server")

	a.Api.SetupRoutes()

	a.ApiServer = httptest.NewServer(a.Api.Server.GetRouter())
	a.ApiUrl = a.ApiServer.URL + "/v1/payments"

	return a.ApiServer
}

func (a *Application) ApiCleanup() {
	a.ApiServer.Close()
}
