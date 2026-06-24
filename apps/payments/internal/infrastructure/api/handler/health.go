package handler

import (
	"net/http"

	"github.com/didiegovieira/go-payments-core/pkg/api"
	"github.com/didiegovieira/go-payments-core/pkg/metrics"
	"github.com/gin-gonic/gin"
)

type Health struct {
	Presenter api.Presenter
}

// Health godoc
// @Summary      Health Check
// @Description  Check if the service is alive
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {boolean}  true "Service is healthy"
// @Router       /v1/payments/health [get]
func (h *Health) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := metrics.TraceHandler(ctx, "Handler.Health")
		defer finish()

		h.Presenter.Present(ctx, gin.H{
			"ok": true,
		}, http.StatusOK)
	}
}
