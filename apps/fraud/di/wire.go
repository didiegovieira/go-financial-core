//go:build wireinject
// +build wireinject

package di

import (
	"github.com/didiegovieira/go-payments-core/apps/fraud/internal/infrastructure/messaging"
	"github.com/google/wire"
)

var wireWorkerSet = wire.NewSet(
	wireCommonSet,
	wire.Struct(new(messaging.Application), "*"),
)

func InitializeWorker() (*messaging.Application, func(), error) {
	wire.Build(wireWorkerSet)
	return &messaging.Application{}, func() {}, nil
}
