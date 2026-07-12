package kafka

import (
	"context"

	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
)

type Consumer struct {
	consumer *sharedkafka.Consumer
}

func NewConsumer(config Config, handler *Handler) (*Consumer, error) {
	consumer, err := sharedkafka.NewConsumer(config.ConsumerConfig(), handler)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Consumer.Kafka.Notification.Start", &err)
	defer finish()

	return c.consumer.Start(ctx)
}

func (c *Consumer) Stop() (err error) {
	_, finish := metrics.TraceWithError(context.Background(), "Consumer.Kafka.Notification.Stop", &err)
	defer finish()

	return c.consumer.Stop()
}
