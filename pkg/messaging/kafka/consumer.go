package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/didiegovieira/go-financial-core/pkg/telemetry"

	kafkago "github.com/segmentio/kafka-go"
)

type Message struct {
	Topic string
	Key   string
	Value []byte
}

type Handler interface {
	Handle(ctx context.Context, message Message) error
}

type Consumer struct {
	reader  *kafkago.Reader
	handler Handler
	group   string
}

func NewConsumer(config ConsumerConfig, handler Handler) (*Consumer, error) {
	config = config.WithDefaults()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, ErrNilHandler
	}

	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:        config.Brokers,
			GroupID:        config.GroupID,
			GroupTopics:    config.Topics,
			MinBytes:       config.MinBytes,
			MaxBytes:       config.MaxBytes,
			MaxWait:        config.MaxWait,
			CommitInterval: config.CommitInterval,
		}),
		handler: handler,
		group:   config.GroupID,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}

		messageCtx := extractContext(ctx, message.Headers)
		err = c.handle(messageCtx, Message{
			Topic: message.Topic,
			Key:   string(message.Key),
			Value: message.Value,
		})
		if err != nil {
			return err
		}

		if err = c.reader.CommitMessages(ctx, message); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}
	}
}

func (c *Consumer) handle(ctx context.Context, message Message) (err error) {
	started := time.Now()
	defer func() { telemetry.ObserveKafka(ctx, message.Topic, c.group, started, err) }()
	return c.handler.Handle(ctx, message)
}

func (c *Consumer) Stop() error {
	return c.reader.Close()
}
