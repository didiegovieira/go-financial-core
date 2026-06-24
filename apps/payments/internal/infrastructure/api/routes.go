package api

import (
	"fmt"

	"github.com/didiegovieira/go-payments-core/apps/payments/internal/settings"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	prefix = "v1/payments"
)

func (a *Application) SetupRoutes() {
	router := a.Server.GetRouter()
	router.Use(otelgin.Middleware(settings.Settings.Metrics.Name))

	// Swagger Docs
	docs.SwaggerInfo.Title = "Go Payments API"
	docs.SwaggerInfo.BasePath = prefix

	// Redirect to swagger docs
	router.GET("/docs/payments", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Swagger Docs
	router.GET("/docs/payments/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Base Routes
	base := router.Group(prefix, a.CorsMiddleware.Handle())
	{
		// Health Check
		base.GET("/health", a.HealthHandler.Handle())
	}

	// Log Registered Routes for Debugging
	for _, route := range router.Routes() {
		if a.BaseApp != nil {
			a.BaseApp.Logger.Infof("Registered route: %s %s", route.Method, route.Path)
		} else {
			fmt.Printf("Registered route: %s %s\n", route.Method, route.Path)
		}
	}
}
