package di

import (
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/infrastructure/messaging/kafka"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/runtime"
	"github.com/google/wire"
)

var wireMessagingSet = wire.NewSet(
	kafka.NewConfigFromCurrentSettings,
	kafka.NewHandler,
	kafka.NewConsumer,
	wire.Bind(new(runtime.Consumer), new(*kafka.Consumer)),
)
