package kafka

import (
	"context"
	"encoding/json"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(config ProducerConfig) (*Producer, error) {
	config = config.WithDefaults()
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Producer{
		writer: &kafkago.Writer{
			Addr:                   kafkago.TCP(config.Brokers...),
			Balancer:               &kafkago.Hash{},
			BatchTimeout:           config.BatchTimeout,
			ReadTimeout:            config.ReadTimeout,
			WriteTimeout:           config.WriteTimeout,
			RequiredAcks:           kafkago.RequireAll,
			AllowAutoTopicCreation: config.AllowAutoTopicCreation,
		},
	}, nil
}

func (p *Producer) PublishRaw(ctx context.Context, topic string, key string, payload []byte) error {
	if topic == "" {
		return ErrEmptyTopic
	}
	if len(payload) == 0 {
		return ErrEmptyPayload
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   payload,
		Headers: injectContext(ctx),
		Time:    time.Now().UTC(),
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func MarshalPayload(payload any) ([]byte, error) {
	switch value := payload.(type) {
	case nil:
		return nil, ErrEmptyPayload
	case []byte:
		return value, nil
	case json.RawMessage:
		return value, nil
	case string:
		return []byte(value), nil
	default:
		return json.Marshal(value)
	}
}
