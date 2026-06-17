package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

// Core tracing functions
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return Tracer.Start(ctx, spanName)
}

func StartSpanWithAttributes(ctx context.Context, spanName string, attrs ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer.Start(ctx, spanName, attrs...)
}

func Finish(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

// Trace context utilities - HIGHLY RECOMMENDED
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// For advanced use cases
func GetTraceInfo(ctx context.Context) (traceID, spanID string, isValid bool) {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	if spanCtx.IsValid() {
		return spanCtx.TraceID().String(), spanCtx.SpanID().String(), true
	}
	return "", "", false
}

// Helper to easily add attributes
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.SetAttributes(attrs...)
	}
}

// Helper to add events
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}
