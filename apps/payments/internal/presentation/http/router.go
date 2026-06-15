package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AppName     string
	Environment string
}

func NewRouter(config RouterConfig) http.Handler {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	registerHealthRoutes(router, config)

	return router
}

func registerHealthRoutes(router *gin.Engine, config RouterConfig) {
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"service":     config.AppName,
			"environment": config.Environment,
		})
	})

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
