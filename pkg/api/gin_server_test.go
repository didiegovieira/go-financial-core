package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

var testServerPort = "3008"

func TestGinContextPreservesRequestTraceAndCancellation(t *testing.T) {
	server := NewGinServer(&http.Server{})
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: trace.FlagsSampled,
	})
	ctx, cancel := context.WithCancel(trace.ContextWithSpanContext(context.Background(), parent))
	defer cancel()
	server.GetRouter().GET("/trace", func(c *gin.Context) {
		assert.Equal(t, parent, trace.SpanContextFromContext(c))
		cancel()
		assert.Equal(t, context.Canceled, c.Err())
		assert.Equal(t, ctx.Done(), c.Done())
		c.Status(http.StatusOK)
	})
	recorder := httptest.NewRecorder()
	server.GetRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/trace", nil).WithContext(ctx))
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestNewGinServer(t *testing.T) {
	httpServer := &http.Server{}
	ginServer := NewGinServer(httpServer)

	assert.Equal(t, ginServer.server, httpServer)
	assert.Equal(t, ginServer.router, httpServer.Handler)
	assert.Len(t, ginServer.router.Routes(), 0)
}

func TestGinServerRegisterRoutes(t *testing.T) {
	httpServer := &http.Server{}
	ginServer := NewGinServer(httpServer)

	assert.IsType(t, ginServer.GetRouter(), &gin.Engine{})
}

// func TestGinServerStartError(t *testing.T) {
// 	test.Setup(t, nil)

// 	httpServer := &http.Server{
// 		Addr: ":error",
// 	}
// 	ginServer := NewGinServer(httpServer)

// 	err := ginServer.Start()
// 	assert.Error(t, err)
// }

func TestGinServerStartAndShutdown(t *testing.T) {
	httpServer := &http.Server{Addr: ":" + testServerPort}
	ginServer := NewGinServer(httpServer)

	ginServer.GetRouter().GET("/test", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte("up"))
	})

	go func() {
		err := ginServer.Start()
		if err == nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	timeout := 10 * time.Millisecond
	iterations := 0
	up := false
	for !up && iterations < 200 {
		iterations++
		_, err := net.DialTimeout("tcp", "localhost:"+testServerPort, timeout)
		up = err == nil
	}

	requestURL := fmt.Sprintf("http://localhost:%s/test", testServerPort)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	responseData, _ := io.ReadAll(res.Body)

	assert.Equal(t, []byte("up"), responseData)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = ginServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestGinServerShutdown(t *testing.T) {
	httpServer := &http.Server{Addr: ":" + testServerPort}
	ginServer := NewGinServer(httpServer)

	err := ginServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestMockGin(t *testing.T) {
	context, _, recorder := MockGin()

	assert.NotNil(t, context)

	assert.NotNil(t, recorder)

	assert.Equal(t, gin.TestMode, gin.Mode())

	assert.NotNil(t, context.Request)
	assert.Equal(t, &http.Request{}, context.Request)
}
