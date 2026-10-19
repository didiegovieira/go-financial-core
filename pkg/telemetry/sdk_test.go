package telemetry

import (
	"context"
	"net"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	metricspb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	tracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

type traceReceiver struct {
	tracepb.UnimplementedTraceServiceServer
	requests chan *tracepb.ExportTraceServiceRequest
}

func (r *traceReceiver) Export(_ context.Context, req *tracepb.ExportTraceServiceRequest) (*tracepb.ExportTraceServiceResponse, error) {
	r.requests <- req
	return &tracepb.ExportTraceServiceResponse{}, nil
}

type metricReceiver struct {
	metricspb.UnimplementedMetricsServiceServer
	requests chan *metricspb.ExportMetricsServiceRequest
}

func (r *metricReceiver) Export(_ context.Context, req *metricspb.ExportMetricsServiceRequest) (*metricspb.ExportMetricsServiceResponse, error) {
	r.requests <- req
	return &metricspb.ExportMetricsServiceResponse{}, nil
}

func preserveGlobals(t *testing.T) {
	tp, mp, propagator := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetTextMapPropagator()
	t.Cleanup(func() { otel.SetTracerProvider(tp); otel.SetMeterProvider(mp); otel.SetTextMapPropagator(propagator) })
}

func TestOfficialSDKExportsAndFlushesOnShutdown(t *testing.T) {
	preserveGlobals(t)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	traces := &traceReceiver{requests: make(chan *tracepb.ExportTraceServiceRequest, 10)}
	metrics := &metricReceiver{requests: make(chan *metricspb.ExportMetricsServiceRequest, 10)}
	tracepb.RegisterTraceServiceServer(server, traces)
	metricspb.RegisterMetricsServiceServer(server, metrics)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	for k, v := range map[string]string{
		"OTEL_SERVICE_NAME": "official-sdk-test", "OTEL_RESOURCE_ATTRIBUTES": "deployment.environment.name=test",
		"OTEL_EXPORTER_OTLP_ENDPOINT": "http://" + listener.Addr().String(),
		"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc", "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "grpc", "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL": "grpc",
		"OTEL_TRACES_EXPORTER": "otlp", "OTEL_METRICS_EXPORTER": "otlp", "OTEL_LOGS_EXPORTER": "none",
		"OTEL_PROPAGATORS": "tracecontext,baggage", "OTEL_METRIC_EXPORT_INTERVAL": "3600000", "OTEL_BSP_SCHEDULE_DELAY": "3600000", "OTEL_TRACES_SAMPLER": "always_on",
	} {
		t.Setenv(k, v)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sdk, err := Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sdk.Shutdown(context.Background()) })
	spanCtx, span := otel.Tracer("test").Start(ctx, "official-sdk-span")
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(spanCtx, carrier)
	if carrier.Get("traceparent") == "" {
		t.Fatal("trace propagation missing")
	}
	span.End()
	counter, err := otel.Meter("test").Int64Counter("sdk.test.requests")
	if err != nil {
		t.Fatal(err)
	}
	counter.Add(spanCtx, 3)
	if err := sdk.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case req := <-traces.requests:
		foundService, foundSpan := false, false
		for _, batch := range req.ResourceSpans {
			for _, a := range batch.Resource.Attributes {
				if a.Key == "service.name" && a.Value.GetStringValue() == "official-sdk-test" {
					foundService = true
				}
			}
			for _, scope := range batch.ScopeSpans {
				for _, s := range scope.Spans {
					if s.Name == "official-sdk-span" {
						foundSpan = true
					}
				}
			}
		}
		if !foundService || !foundSpan {
			t.Fatal("trace resource or span lost")
		}
	case <-ctx.Done():
		t.Fatal("trace not flushed")
	}
	select {
	case req := <-metrics.requests:
		found := false
		for _, batch := range req.ResourceMetrics {
			for _, scope := range batch.ScopeMetrics {
				for _, m := range scope.Metrics {
					if m.Name == "sdk.test.requests" && m.GetSum().DataPoints[0].GetAsInt() == 3 {
						found = true
					}
				}
			}
		}
		if !found {
			t.Fatal("metric value not flushed")
		}
	case <-ctx.Done():
		t.Fatal("metrics not flushed")
	}
}

func TestDisabledExportersAndInvalidConfiguration(t *testing.T) {
	preserveGlobals(t)
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	t.Setenv("OTEL_LOGS_EXPORTER", "none")
	t.Setenv("OTEL_PROPAGATORS", "tracecontext,baggage")
	sdk, err := Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if sdk.traces != nil || sdk.metrics != nil {
		t.Fatal("disabled exporters initialized")
	}
	if err := sdk.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ key, value string }{
		{"OTEL_TRACES_EXPORTER", "jaeger"}, {"OTEL_METRICS_EXPORTER", "prometheus"},
		{"OTEL_LOGS_EXPORTER", "otlp"}, {"OTEL_PROPAGATORS", "unknown"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			t.Setenv(tc.key, tc.value)
			if _, err := Start(context.Background()); err == nil {
				t.Fatal("unsupported configuration accepted")
			}
		})
	}
}
