package wallet

import (
	"net/http"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	walletUseCase "github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/wallet"
	"github.com/didiegovieira/go-financial-core/pkg/api"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

type Debit struct {
	Presenter    api.Presenter
	DebitUseCase walletUseCase.DebitUseCase
}

// Debit godoc
// @Summary      Debit Wallet
// @Description  Remove an amount in cents from an active wallet
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        input body dto.DebitInput true "Debit data"
// @Success      200 {object} dto.DebitOutput
// @Failure      400 {object} api.HttpError
// @Router       /v1/financial/wallets/debits [post]
func (c *Debit) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := telemetry.TraceHandler(ctx, "Handler.Wallet.Debit")
		defer finish()

		var input dto.DebitInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		output, err := c.DebitUseCase.Execute(ctx, input)
		if err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		c.Presenter.Present(ctx, output, http.StatusOK)
	}
}
