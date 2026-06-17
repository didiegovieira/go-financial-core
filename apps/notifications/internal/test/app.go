package test

import (
	"net/http/httptest"

	"github.com/didiegovieira/go-payments-core/apps/notifications/internal/application"
	"go.uber.org/mock/gomock"
)

type Application struct {
	BaseApp *application.App

	MockCtrl *gomock.Controller

	ApiUrl    string           `wire:"-"`
	ApiServer *httptest.Server `wire:"-"`
}

func (a *Application) RunApiServer() *httptest.Server {
	a.BaseApp.Start("test-api-server")

	return a.ApiServer
}

func (a *Application) ApiCleanup() {
	a.ApiServer.Close()
}
