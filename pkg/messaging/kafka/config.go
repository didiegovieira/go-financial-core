package kafka

import (
	"strings"
	"time"
)

type ProducerConfig struct {
	Brokers                []string
	BatchTimeout           time.Duration
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
	AllowAutoTopicCreation bool
}

type ConsumerConfig struct {
	Brokers        []string
	GroupID        string
	Topics         []string
	MinBytes       int
	MaxBytes       int
	MaxWait        time.Duration
	CommitInterval time.Duration
}

func (c ProducerConfig) WithDefaults() ProducerConfig {
	c.Brokers = NormalizeStrings(c.Brokers)
	if c.BatchTimeout == 0 {
		c.BatchTimeout = 10 * time.Millisecond
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 10 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 10 * time.Second
	}

	return c
}

func (c ProducerConfig) Validate() error {
	if len(c.Brokers) == 0 {
		return ErrEmptyBrokers
	}

	return nil
}

func (c ConsumerConfig) WithDefaults() ConsumerConfig {
	c.Brokers = NormalizeStrings(c.Brokers)
	c.Topics = NormalizeStrings(c.Topics)
	if c.MinBytes == 0 {
		c.MinBytes = 1
	}
	if c.MaxBytes == 0 {
		c.MaxBytes = 10e6
	}
	if c.MaxWait == 0 {
		c.MaxWait = 10 * time.Second
	}

	return c
}

func (c ConsumerConfig) Validate() error {
	if len(c.Brokers) == 0 {
		return ErrEmptyBrokers
	}
	if c.GroupID == "" {
		return ErrEmptyGroupID
	}
	if len(c.Topics) == 0 {
		return ErrEmptyTopics
	}

	return nil
}

func NormalizeStrings(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			normalized = append(normalized, value)
		}
	}

	return normalized
}
