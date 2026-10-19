package http

import (
	"context"
	"encoding/json"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/domain/entity"
	sharedhttp "github.com/didiegovieira/go-financial-core/pkg/http"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type NotificationClient struct {
	wrapper sharedhttp.Wrapper
	url     string
}

func NewNotificationClient(wrapper sharedhttp.Wrapper, url string) *NotificationClient {
	return &NotificationClient{
		wrapper: wrapper,
		url:     url,
	}
}

func (c *NotificationClient) Send(ctx context.Context, notification entity.Notification) (err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Client.Notification.HTTP.Send", &err)
	defer finish()

	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	response, err := c.wrapper.Post(ctx, c.url, sharedhttp.Request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return response.Error()
}
