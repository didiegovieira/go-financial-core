package kafka

import (
	"time"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
)

type Config struct {
	Brokers                []string
	BatchTimeout           time.Duration
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
	AllowAutoTopicCreation bool
}

func NewConfigFromSettings(config settings.KafkaSpecification) Config {
	return Config{
		Brokers: config.Brokers,
	}
}

func (c Config) WithDefaults() Config {
	producerConfig := c.ProducerConfig().WithDefaults()
	c.Brokers = producerConfig.Brokers
	c.BatchTimeout = producerConfig.BatchTimeout
	c.ReadTimeout = producerConfig.ReadTimeout
	c.WriteTimeout = producerConfig.WriteTimeout
	c.AllowAutoTopicCreation = producerConfig.AllowAutoTopicCreation

	return c
}

func (c Config) Validate() error {
	return c.ProducerConfig().Validate()
}

func (c Config) ProducerConfig() sharedkafka.ProducerConfig {
	return sharedkafka.ProducerConfig{
		Brokers:                c.Brokers,
		BatchTimeout:           c.BatchTimeout,
		ReadTimeout:            c.ReadTimeout,
		WriteTimeout:           c.WriteTimeout,
		AllowAutoTopicCreation: c.AllowAutoTopicCreation,
	}
}
