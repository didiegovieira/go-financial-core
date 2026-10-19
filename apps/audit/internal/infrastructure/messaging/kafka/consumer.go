package kafka

import (
	"context"

	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
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
	ctx, finish := telemetry.TraceWithError(ctx, "Consumer.Kafka.Audit.Start", &err)
	defer finish()

	return c.consumer.Start(ctx)
}

func (c *Consumer) Stop() (err error) {
	_, finish := telemetry.TraceWithError(context.Background(), "Consumer.Kafka.Audit.Stop", &err)
	defer finish()

	return c.consumer.Stop()
}
