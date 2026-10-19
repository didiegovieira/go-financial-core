package kafka

import (
	"context"
	"strings"

	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
	kafkago "github.com/segmentio/kafka-go"
)

// headerCarrier transports OpenTelemetry context without changing event payloads.
type headerCarrier []kafkago.Header

func (c headerCarrier) Get(key string) string {
	for _, header := range c {
		if strings.EqualFold(header.Key, key) {
			return string(header.Value)
		}
	}
	return ""
}

func (c *headerCarrier) Set(key, value string) {
	for i := range *c {
		if strings.EqualFold((*c)[i].Key, key) {
			(*c)[i].Value = []byte(value)
			return
		}
	}
	*c = append(*c, kafkago.Header{Key: key, Value: []byte(value)})
}

func (c headerCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for _, header := range c {
		keys = append(keys, header.Key)
	}
	return keys
}

func injectContext(ctx context.Context) []kafkago.Header {
	var carrier headerCarrier
	telemetry.InjectContext(ctx, &carrier)
	return []kafkago.Header(carrier)
}

func extractContext(ctx context.Context, headers []kafkago.Header) context.Context {
	carrier := headerCarrier(headers)
	return telemetry.ExtractContext(ctx, &carrier)
}
