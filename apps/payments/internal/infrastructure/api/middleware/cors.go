package middleware

import (
	"github.com/didiegovieira/go-payments-core/pkg/api"
	"github.com/gin-gonic/gin"
)

type Cors struct {
	Presenter api.Presenter
}

func (c Cors) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Next()
	}
}
