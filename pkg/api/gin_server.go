package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type GinServer[T any] struct {
	router T
	server *http.Server
}

func NewGinServer[T *gin.Engine](httpServer *http.Server) *GinServer[T] {
	router := gin.Default()
	httpServer.Handler = router

	return &GinServer[T]{
		router: router,
		server: httpServer,
	}
}

func (s *GinServer[T]) GetRouter() T {
	return s.router
}

func (s *GinServer[T]) Start() error {
	fmt.Printf("Starting web server on port %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *GinServer[T]) Shutdown(ctx context.Context) error {
	fmt.Printf("Stopping web server")
	return s.server.Shutdown(ctx)
}

func MockGin() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	responseRecorder := httptest.NewRecorder()
	c, e := gin.CreateTestContext(responseRecorder)
	c.Request = &http.Request{}
	return c, e, responseRecorder
}
