package api

import (
	"fmt"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	prefix = "v1/financial"
)

func (a *Application) SetupRoutes() {
	router := a.Server.GetRouter()
	router.Use(otelgin.Middleware(settings.Settings.Metrics.Name))

	// Swagger Docs
	docs.SwaggerInfo.Title = "Go Financial Core API"
	docs.SwaggerInfo.BasePath = prefix

	// Redirect to swagger docs
	router.GET("/docs/financial", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Swagger Docs
	router.GET("/docs/financial/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Base Routes
	base := router.Group(prefix, a.CorsMiddleware.Handle())
	{
		// Health Check
		base.GET("/health", a.HealthHandler.Handle())

		accounts := base.Group("/accounts")
		{
			accounts.POST("", a.AccountCreateHandler.Handle())
		}

		payments := base.Group("/payments")
		{
			payments.POST("", a.PaymentCreateHandler.Handle())
		}
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
