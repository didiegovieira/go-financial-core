package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer = trace.NewNoopTracerProvider().Tracer("noop")

// Core tracing functions
func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer().Start(ctx, spanName, opts...)
}

func StartSpanWithAttributes(ctx context.Context, spanName string, attrs ...trace.SpanStartOption) (context.Context, trace.Span) {
	return StartSpan(ctx, spanName, attrs...)
}

func Trace(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, func()) {
	ctx, span := StartSpan(ctx, spanName, opts...)
	return ctx, func() {
		Finish(span, nil)
	}
}

func TraceWithError(ctx context.Context, spanName string, err *error, opts ...trace.SpanStartOption) (context.Context, func()) {
	ctx, span := StartSpan(ctx, spanName, opts...)
	return ctx, func() {
		if err == nil {
			Finish(span, nil)
			return
		}

		Finish(span, *err)
	}
}

func Finish(span trace.Span, err error) {
	if span == nil {
		return
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
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

func tracer() trace.Tracer {
	if Tracer != nil {
		return Tracer
	}

	return trace.NewNoopTracerProvider().Tracer("noop")
}
