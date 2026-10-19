package telemetry

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// Operational records contain correlation and bounded dimensions, not payloads.
func (o *Observer) logHTTP(ctx context.Context, level slog.Level, method, route string, status int, duration time.Duration) {
	o.logger.Log(ctx, level, "http_request_completed", "service", os.Getenv("OTEL_SERVICE_NAME"), "trace_id", GetTraceID(ctx), "span_id", GetSpanID(ctx), "method", method, "route", route, "status", status, "duration_seconds", duration.Seconds())
}

func (o *Observer) logKafka(ctx context.Context, level slog.Level, topic, group, result string, duration time.Duration) {
	o.logger.Log(ctx, level, "kafka_message_processed", "service", os.Getenv("OTEL_SERVICE_NAME"), "trace_id", GetTraceID(ctx), "span_id", GetSpanID(ctx), "topic", topic, "consumer_group", group, "result", result, "duration_seconds", duration.Seconds())
}
