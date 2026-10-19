// Package telemetry centralizes SDK setup, tracing, metrics, correlated logs,
// and context propagation shared by all applications.
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

type SDK struct {
	traces  *sdktrace.TracerProvider
	metrics *sdkmetric.MeterProvider
}

// Start reads OTEL_* configuration. Logs remain collected from Docker by Alloy.
// Call once after loading application environment, before composing dependencies.
func Start(ctx context.Context) (_ *SDK, err error) {
	for _, signal := range []string{"TRACES", "METRICS"} {
		if value := env("OTEL_"+signal+"_EXPORTER", "otlp"); value != "otlp" && value != "none" {
			return nil, fmt.Errorf("unsupported OTEL_%s_EXPORTER: %q (use otlp or none)", signal, value)
		}
		if env("OTEL_"+signal+"_EXPORTER", "otlp") == "otlp" {
			if protocol := protocol(signal); protocol != "grpc" && protocol != "http/protobuf" {
				return nil, fmt.Errorf("unsupported OTLP %s protocol: %q", signal, protocol)
			}
		}
	}
	if env("OTEL_LOGS_EXPORTER", "none") != "none" {
		return nil, errors.New("OTEL_LOGS_EXPORTER must be none: Alloy collects stdout/stderr")
	}
	var propagators []propagation.TextMapPropagator
	for _, name := range strings.Split(env("OTEL_PROPAGATORS", "tracecontext,baggage"), ",") {
		switch strings.TrimSpace(name) {
		case "tracecontext":
			propagators = append(propagators, propagation.TraceContext{})
		case "baggage":
			propagators = append(propagators, propagation.Baggage{})
		case "none":
		default:
			return nil, fmt.Errorf("unsupported OTEL_PROPAGATORS entry: %q", name)
		}
	}
	res, err := resource.New(ctx, resource.WithAttributes(attribute.String("service.name", "go-financial-core")), resource.WithFromEnv(), resource.WithTelemetrySDK(), resource.WithProcess(), resource.WithHost())
	if err != nil {
		return nil, fmt.Errorf("telemetry resource: %w", err)
	}
	sdk := &SDK{}
	defer func() {
		if err != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = errors.Join(err, sdk.Shutdown(cleanupCtx))
		}
	}()
	if env("OTEL_TRACES_EXPORTER", "otlp") != "none" {
		var exporter sdktrace.SpanExporter
		if protocol("TRACES") == "http/protobuf" {
			exporter, err = otlptracehttp.New(ctx)
		} else {
			exporter, err = otlptracegrpc.New(ctx)
		}
		if err != nil {
			return nil, fmt.Errorf("trace exporter: %w", err)
		}
		opts := []sdktrace.TracerProviderOption{sdktrace.WithResource(res), sdktrace.WithBatcher(exporter)}
		// Preserve the previous distribution's default; explicit OTEL sampler wins.
		if os.Getenv("OTEL_TRACES_SAMPLER") == "" {
			opts = append(opts, sdktrace.WithSampler(sdktrace.AlwaysSample()))
		}
		sdk.traces = sdktrace.NewTracerProvider(opts...)
	}
	if env("OTEL_METRICS_EXPORTER", "otlp") != "none" {
		var exporter sdkmetric.Exporter
		if protocol("METRICS") == "http/protobuf" {
			exporter, err = otlpmetrichttp.New(ctx)
		} else {
			exporter, err = otlpmetricgrpc.New(ctx)
		}
		if err != nil {
			return nil, fmt.Errorf("metric exporter: %w", err)
		}
		sdk.metrics = sdkmetric.NewMeterProvider(sdkmetric.WithResource(res), sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
		if err = runtime.Start(runtime.WithMeterProvider(sdk.metrics)); err != nil {
			return nil, fmt.Errorf("runtime instrumentation: %w", err)
		}
	}
	// Publish global providers only after initialization has succeeded.
	if sdk.traces != nil {
		otel.SetTracerProvider(sdk.traces)
	} else {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
	}
	if sdk.metrics != nil {
		otel.SetMeterProvider(sdk.metrics)
	} else {
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagators...))
	return sdk, nil
}

func (s *SDK) Shutdown(ctx context.Context) error {
	var errs []error
	if s.traces != nil {
		errs = append(errs, s.traces.Shutdown(ctx))
	}
	if s.metrics != nil {
		errs = append(errs, s.metrics.Shutdown(ctx))
	}
	return errors.Join(errs...)
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return strings.ToLower(value)
	}
	return fallback
}

func protocol(signal string) string {
	return env("OTEL_EXPORTER_OTLP_"+signal+"_PROTOCOL", env("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc"))
}
