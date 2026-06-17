package api

import "github.com/gin-gonic/gin"

type Presenter interface {
	Error(c *gin.Context, err error)
	Present(c *gin.Context, body interface{}, code int)
}

type HttpError struct {
	Error string `json:"error"`
}

type HttpJSON struct {
	Key string `json:"key"`
}
