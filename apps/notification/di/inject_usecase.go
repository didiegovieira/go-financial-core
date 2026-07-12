package di

import (
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/usecase"
	"github.com/google/wire"
)

var wireUsecaseSet = wire.NewSet(
	usecase.NewSendNotification,
	wire.Bind(new(usecase.SendNotificationUseCase), new(*usecase.SendNotification)),
)
