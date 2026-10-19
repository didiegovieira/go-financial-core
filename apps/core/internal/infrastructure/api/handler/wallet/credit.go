package wallet

import (
	"net/http"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	walletUseCase "github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/wallet"
	"github.com/didiegovieira/go-financial-core/pkg/api"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

type Credit struct {
	Presenter     api.Presenter
	CreditUseCase walletUseCase.CreditUseCase
}

// Credit godoc
// @Summary      Credit Wallet
// @Description  Add an amount in cents to an active wallet
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        input body dto.CreditInput true "Credit data"
// @Success      200 {object} dto.CreditOutput
// @Router       /v1/financial/wallets/credits [post]
func (c *Credit) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := telemetry.TraceHandler(ctx, "Handler.Wallet.Credit")
		defer finish()

		var input dto.CreditInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		output, err := c.CreditUseCase.Execute(ctx, input)
		if err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		c.Presenter.Present(ctx, output, http.StatusOK)
	}
}
