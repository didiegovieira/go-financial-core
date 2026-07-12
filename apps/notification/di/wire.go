//go:build wireinject
// +build wireinject

package di

import (
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/runtime"
	"github.com/google/wire"
)

var wireWorkerSet = wire.NewSet(
	wireCommonSet,
	wireClientSet,
	wireUsecaseSet,
	wireMessagingSet,
	wire.Struct(new(runtime.Worker), "*"),
)

func InitializeWorker() (*runtime.Worker, func(), error) {
	wire.Build(wireWorkerSet)
	return &runtime.Worker{}, func() {}, nil
}
