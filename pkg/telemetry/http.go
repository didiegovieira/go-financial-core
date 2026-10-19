package telemetry

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// TraceHTTP instruments incoming requests using the configured SDK provider.
func TraceHTTP(serviceName string) gin.HandlerFunc { return otelgin.Middleware(serviceName) }

func TraceHandler(ctx *gin.Context, spanName string, opts ...trace.SpanStartOption) func() {
	reqCtx, finish := Trace(ctx.Request.Context(), spanName, opts...)
	ctx.Request = ctx.Request.WithContext(reqCtx)

	return finish
}

func TraceHandlerWithError(ctx *gin.Context, spanName string, err *error, opts ...trace.SpanStartOption) func() {
	reqCtx, finish := TraceWithError(ctx.Request.Context(), spanName, err, opts...)
	ctx.Request = ctx.Request.WithContext(reqCtx)

	return finish
}

func ObserveHTTP() gin.HandlerFunc { return operational.HTTP() }

func (o *Observer) HTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := c.Writer.Status()
		attrs := metric.WithAttributes(attribute.String("method", c.Request.Method), attribute.String("route", route), attribute.String("status", strconv.Itoa(status)))
		ctx := c.Request.Context()
		o.httpRequests.Add(ctx, 1, attrs)
		o.httpDuration.Record(ctx, time.Since(started).Seconds(), attrs)
		level := slog.LevelInfo
		if status >= 400 {
			level = slog.LevelWarn
		}
		if status >= 500 {
			level = slog.LevelError
		}
		o.logHTTP(ctx, level, c.Request.Method, route, status, time.Since(started))
	}
}
