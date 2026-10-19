package telemetry

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Observer records bounded operational dimensions, never account IDs or payloads.
type Observer struct {
	httpRequests  metric.Int64Counter
	httpDuration  metric.Float64Histogram
	kafkaMessages metric.Int64Counter
	kafkaDuration metric.Float64Histogram
	logger        *slog.Logger
}

func NewObserver(meter metric.Meter, logger *slog.Logger) (*Observer, error) {
	o := &Observer{logger: logger}
	var err error
	if o.httpRequests, err = meter.Int64Counter("financial.http.requests"); err != nil {
		return nil, err
	}
	buckets := metric.WithExplicitBucketBoundaries(.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)
	if o.httpDuration, err = meter.Float64Histogram("financial.http.duration", metric.WithUnit("s"), buckets); err != nil {
		return nil, err
	}
	if o.kafkaMessages, err = meter.Int64Counter("financial.kafka.messages"); err != nil {
		return nil, err
	}
	if o.kafkaDuration, err = meter.Float64Histogram("financial.kafka.duration", metric.WithUnit("s"), buckets); err != nil {
		return nil, err
	}
	return o, nil
}

var operational = func() *Observer {
	o, err := NewObserver(otel.Meter("financial-core/operations"), slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err != nil {
		panic(err)
	}
	return o
}()

// ObserveKafka records handler completion, not offset commits or delivery guarantees.
func ObserveKafka(ctx context.Context, topic, group string, started time.Time, err error) {
	operational.Kafka(ctx, topic, group, time.Since(started), err)
}

func (o *Observer) Kafka(ctx context.Context, topic, group string, duration time.Duration, err error) {
	result, level := "success", slog.LevelInfo
	if err != nil {
		result, level = "error", slog.LevelError
	}
	attrs := metric.WithAttributes(attribute.String("topic", topic), attribute.String("consumer_group", group), attribute.String("result", result))
	o.kafkaMessages.Add(ctx, 1, attrs)
	o.kafkaDuration.Record(ctx, duration.Seconds(), attrs)
	o.logKafka(ctx, level, topic, group, result, duration)
}
