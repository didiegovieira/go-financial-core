package kafka

import "errors"

var (
	ErrEmptyBrokers     = errors.New("kafka brokers cannot be empty")
	ErrEmptyTopic       = errors.New("kafka topic cannot be empty")
	ErrEmptyPayload     = errors.New("kafka payload cannot be empty")
	ErrEmptyGroupID     = errors.New("kafka consumer group id cannot be empty")
	ErrEmptyTopics      = errors.New("kafka consumer topics cannot be empty")
	ErrNilHandler       = errors.New("kafka handler cannot be nil")
	ErrUnknownEventName = errors.New("unknown event name")
)
