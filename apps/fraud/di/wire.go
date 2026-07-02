//go:build wireinject
// +build wireinject

package di

import (
	"github.com/didiegovieira/go-financial-core/apps/fraud/internal/service"
	"github.com/google/wire"
)

var wireWorkerSet = wire.NewSet(
	wireCommonSet,
	wire.Struct(new(service.Worker), "*"),
)

func InitializeWorker() (*service.Worker, func(), error) {
	wire.Build(wireWorkerSet)
	return &service.Worker{}, func() {}, nil
}
