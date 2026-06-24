package wallet

import (
	"net/http"

	dto "github.com/didiegovieira/go-payments-core/apps/payments/internal/application/dto/wallet"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/application/usecase/wallet"
	"github.com/didiegovieira/go-payments-core/pkg/api"
	"github.com/didiegovieira/go-payments-core/pkg/metrics"
	"github.com/gin-gonic/gin"
)

type Create struct {
	Presenter     api.Presenter
	CreateUseCase wallet.CreateUseCase
}

// Create godoc
// @Summary      Create Wallet
// @Description  Create a new wallet
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Success      201  {object}  api.Response "Wallet created successfully"
// @Router       /v1/wallets [post]
func (c *Create) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := metrics.TraceHandler(ctx, "Handler.Wallet.Create")
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
