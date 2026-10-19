package kafka

import (
	"context"
	"testing"

	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestPropagationAcrossKafkaWorkers(t *testing.T) {
	previous := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	t.Cleanup(func() { otel.SetTextMapPropagator(previous) })

	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3},
		SpanID:     trace.SpanID{4, 5, 6},
		TraceFlags: trace.FlagsSampled,
	})
	member, err := baggage.NewMember("test-source", "poc")
	if err != nil {
		t.Fatal(err)
	}
	bag, err := baggage.New(member)
	if err != nil {
		t.Fatal(err)
	}
	ctx := baggage.ContextWithBaggage(trace.ContextWithSpanContext(context.Background(), parent), bag)
	headers := injectContext(ctx)
	for _, worker := range []string{"audit", "notifications"} {
		t.Run(worker, func(t *testing.T) {
			ctx := extractContext(context.Background(), headers)
			got := trace.SpanContextFromContext(ctx)
			if got.TraceID() != parent.TraceID() || got.SpanID() != parent.SpanID() || !got.IsRemote() || !got.IsSampled() {
				t.Fatalf("worker lost parent trace context: %v", got)
			}
			if baggage.FromContext(ctx).Member("test-source").Value() != "poc" {
				t.Fatal("worker lost baggage")
			}
		})
	}
	// Older events and malformed headers must remain consumable without a parent.
	for _, headers := range [][]kafkago.Header{nil, {{Key: "traceparent", Value: []byte("invalid")}}} {
		if trace.SpanContextFromContext(extractContext(context.Background(), headers)).IsValid() {
			t.Fatal("unexpected parent for an event without valid trace context")
		}
	}
}
