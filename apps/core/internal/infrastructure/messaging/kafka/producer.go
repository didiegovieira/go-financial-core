package kafka

import (
	"context"

	appmessaging "github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
)

type Producer struct {
	producer *sharedkafka.Producer
}

func NewProducer(config Config) (*Producer, error) {
	producer, err := sharedkafka.NewProducer(config.ProducerConfig())
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
	}, nil
}

func NewProducerFromSettings() (*Producer, error) {
	return NewProducer(NewConfigFromSettings(settings.Settings.Kafka))
}

func (p *Producer) Publish(ctx context.Context, event appmessaging.Event) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Messaging.Kafka.Publish", &err)
	defer finish()

	topic, ok := topicByEventName[event.Name]
	if !ok {
		return sharedkafka.ErrUnknownEventName
	}

	payload, err := sharedkafka.MarshalPayload(event.Payload)
	if err != nil {
		return err
	}

	return p.PublishRaw(ctx, topic, event.Key, payload)
}

func (p *Producer) PublishRaw(ctx context.Context, topic string, key string, payload []byte) (err error) {
	ctx, finish := metrics.TraceWithError(ctx, "Messaging.Kafka.PublishRaw", &err)
	defer finish()

	return p.producer.PublishRaw(ctx, topic, key, payload)
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

var _ appmessaging.Publisher = (*Producer)(nil)
