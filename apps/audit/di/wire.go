//go:build wireinject
// +build wireinject

package di

import (
	"github.com/didiegovieira/go-payments-core/apps/audit/internal/infrastructure/messaging"
	"github.com/google/wire"
)

var wireApiSet = wire.NewSet(

// wire.Struct(new(messaging.Application), "*"),
)

var wireTestSet = wire.NewSet(

// wire.Struct(new(api.Application), "*"),
// wire.Struct(new(test.Application), "*"),
)

func InitializeWorker() (*messaging.Application, func(), error) {
	wire.Build(wireApiSet)
	return &messaging.Application{}, func() {}, nil
}

// func InitilizeTests(mockCtrl *gomock.Controller) (*test.Application, func(), error) {
// 	wire.Build(wireTestSet)

// 	return &test.Application{}, func() {}, nil
// }
