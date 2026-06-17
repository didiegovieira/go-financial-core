package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

const (
	prefix = "v1/payments"
)

func (a *Application) SetupRoutes() {
	router := a.Server.GetRouter()

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
	base := router.Group(prefix)
	{
		// Health Check
		base.GET("/health", a.HealthHandler.Handle())
	}

	// Log Registered Routes for Debugging
	for _, route := range router.Routes() {
		if a.BaseApp != nil {
			a.BaseApp.Logger.Printf("Registered route: %s %s", route.Method, route.Path)
		} else {
			fmt.Printf("Registered route: %s %s\n", route.Method, route.Path)
		}
	}
}
