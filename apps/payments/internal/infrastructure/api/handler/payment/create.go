package payment

import (
	"net/http"

	dto "github.com/didiegovieira/go-payments-core/apps/payments/internal/application/dto/payment"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application/usecase/payment"
	"github.com/didiegovieira/go-payments-core/pkg/api"
	"github.com/didiegovieira/go-payments-core/pkg/metrics"
	"github.com/gin-gonic/gin"
)

type Create struct {
	Presenter     api.Presenter
	CreateUseCase payment.CreateUseCase
}

// Create godoc
// @Summary      Create Payment
// @Description  Create a new payment
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Success      201  {object}  api.Response "Payment created successfully"
// @Router       /v1/payments [post]
func (c *Create) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := metrics.TraceHandler(ctx, "Handler.Payment.Create")
		defer finish()

		var input dto.CreateInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			c.Presenter.Error(ctx, err)
		}

		output, err := c.CreateUseCase.Execute(ctx, input)
		if err != nil {
			c.Presenter.Error(ctx, err)
		}

		c.Presenter.Present(ctx, output, http.StatusCreated)
	}
}
