package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"
)

func TestOperationalMetricsAndCorrelatedLogs(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	var logs bytes.Buffer
	o, err := NewObserver(provider.Meter("test"), slog.New(slog.NewJSONHandler(&logs, nil)))
	if err != nil {
		t.Fatal(err)
	}
	parent := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}})
	ctx := trace.ContextWithSpanContext(context.Background(), parent)
	router := gin.New()
	router.Use(o.HTTP())
	router.GET("/accounts/:id", func(c *gin.Context) { c.Status(http.StatusBadRequest) })
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/accounts/private-id?phone=secret", nil).WithContext(ctx))
	o.Kafka(ctx, "core.account.created", "audit-service", time.Second, nil)
	o.Kafka(ctx, "core.account.created", "audit-service", 2*time.Second, errors.New("sensitive payload must not appear"))
	var data metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &data); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, scope := range data.ScopeMetrics {
		for _, m := range scope.Metrics {
			seen[m.Name] = true
			switch m.Name {
			case "financial.http.requests":
				points := m.Data.(metricdata.Sum[int64]).DataPoints
				if len(points) != 1 || points[0].Value != 1 {
					t.Fatalf("HTTP count: %v", points)
				}
				attrs := points[0].Attributes.ToSlice()
				for _, attr := range attrs {
					if attr.Key == "route" && attr.Value.AsString() != "/accounts/:id" {
						t.Fatal("unbounded route label")
					}
				}
			case "financial.kafka.messages":
				points := m.Data.(metricdata.Sum[int64]).DataPoints
				if len(points) != 2 {
					t.Fatalf("expected separate success/error series: %v", points)
				}
				for _, p := range points {
					if p.Value != 1 {
						t.Fatal("incorrect Kafka count")
					}
				}
			case "financial.kafka.duration":
				points := m.Data.(metricdata.Histogram[float64]).DataPoints
				if len(points) != 2 || points[0].Count != 1 || points[1].Count != 1 {
					t.Fatal("missing Kafka duration samples")
				}
			}
		}
	}
	for _, name := range []string{"financial.http.requests", "financial.http.duration", "financial.kafka.messages", "financial.kafka.duration"} {
		if !seen[name] {
			t.Fatalf("missing metric %s", name)
		}
	}
	if bytes.Contains(logs.Bytes(), []byte("private-id")) || bytes.Contains(logs.Bytes(), []byte("secret")) || bytes.Contains(logs.Bytes(), []byte("sensitive")) {
		t.Fatal("sensitive request data logged")
	}
	decoder := json.NewDecoder(&logs)
	for i := 0; i < 3; i++ {
		var row map[string]any
		if err := decoder.Decode(&row); err != nil {
			t.Fatal(err)
		}
		if row["trace_id"] != parent.TraceID().String() {
			t.Fatal("missing trace correlation")
		}
		if i == 2 && (row["level"] != "ERROR" || row["result"] != "error") {
			t.Fatal("failure was not logged as an error")
		}
	}
}
