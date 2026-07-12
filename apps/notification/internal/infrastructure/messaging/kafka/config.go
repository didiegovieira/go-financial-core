package kafka

import (
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/settings"
	sharedkafka "github.com/didiegovieira/go-financial-core/pkg/messaging/kafka"
)

type Config struct {
	Brokers []string
	GroupID string
	Topics  []string
}

func NewConfigFromSettings(config settings.KafkaSpecification) Config {
	return Config{
		Brokers: config.Brokers,
		GroupID: config.ConsumerGroupID,
		Topics:  Topics,
	}
}

func NewConfigFromCurrentSettings() Config {
	return NewConfigFromSettings(settings.Settings.Kafka)
}

func (c Config) ConsumerConfig() sharedkafka.ConsumerConfig {
	return sharedkafka.ConsumerConfig{
		Brokers: c.Brokers,
		GroupID: c.GroupID,
		Topics:  c.Topics,
	}
}
