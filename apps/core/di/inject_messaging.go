package di

import (
	appmessaging "github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/messaging/kafka"
	"github.com/google/wire"
)

func provideKafkaProducer() (*kafka.Producer, func(), error) {
	producer, err := kafka.NewProducerFromSettings()
	if err != nil {
		return nil, nil, err
	}

	return producer, func() {
		_ = producer.Close()
	}, nil
}

var wireMessagingSet = wire.NewSet(
	provideKafkaProducer,
	wire.Bind(new(appmessaging.Publisher), new(*kafka.Producer)),
)
