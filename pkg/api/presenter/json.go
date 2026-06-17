package presenter

import (
	"net/http"

	appErr "github.com/didiegovieira/go-payments-core/pkg/errors"
	"github.com/didiegovieira/go-payments-core/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type JsonConfig struct {
	IncludeTraceIDInHeader bool
	IncludeTraceIDInBody   bool
	TraceIDHeaderName      string
}

type Json struct {
	config JsonConfig
}

func NewJson() Json {
	return Json{
		config: JsonConfig{
			IncludeTraceIDInHeader: true,
			IncludeTraceIDInBody:   false,
			TraceIDHeaderName:      "X-Trace-ID",
		},
	}
}

func NewJsonWithConfig(config JsonConfig) Json {
	return Json{config: config}
}

func (j Json) setTraceID(c *gin.Context, response gin.H) {
	traceID := metrics.GetTraceID(c.Request.Context())
	if traceID == "" {
		return
	}

	if j.config.IncludeTraceIDInHeader {
		c.Header(j.config.TraceIDHeaderName, traceID)
	}

	if j.config.IncludeTraceIDInBody && response != nil {
		response["trace_id"] = traceID
	}
}

func (j Json) Error(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	response := gin.H{"error": err.Error()}

	switch e := err.(type) {
	case *appErr.Validation:
		code = http.StatusBadRequest
		response["messages"] = e.Errors

	case *appErr.Http:
		code = e.Code
	}

	j.setTraceID(c, response)

	c.AbortWithStatusJSON(code, response)
}

func (j Json) Present(c *gin.Context, body interface{}, code int) {
	j.setTraceID(c, nil)

	c.JSON(code, body)
}

func (j Json) PresentWithHeaders(c *gin.Context, body interface{}, code int, headers map[string]string) {
	j.setTraceID(c, nil)

	for key, value := range headers {
		c.Header(key, value)
	}

	c.JSON(code, body)
}
