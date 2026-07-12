package di

import (
	appclient "github.com/didiegovieira/go-financial-core/apps/notification/internal/application/port/client"
	noopclient "github.com/didiegovieira/go-financial-core/apps/notification/internal/infrastructure/client/noop"
	"github.com/google/wire"
)

var wireClientSet = wire.NewSet(
	noopclient.NewNotificationClient,
	wire.Bind(new(appclient.NotificationClient), new(*noopclient.NotificationClient)),
)
