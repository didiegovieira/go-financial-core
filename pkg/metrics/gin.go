package metrics

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

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
