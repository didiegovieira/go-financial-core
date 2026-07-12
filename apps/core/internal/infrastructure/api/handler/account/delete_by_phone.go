package account

import (
	"net/http"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/account"

	"github.com/didiegovieira/go-financial-core/pkg/api"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
	"github.com/gin-gonic/gin"
)

type DeleteByPhone struct {
	Presenter            api.Presenter
	DeleteByPhoneUseCase account.DeleteByPhoneUseCase
}

// DeleteByPhone godoc
// @Summary      Delete Account by Phone
// @Description  Delete an existing account by phone number
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.DeleteByPhoneOutput "Account found successfully"
// @Router       /v1/financial/accounts [get]
func (c *DeleteByPhone) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := metrics.TraceHandler(ctx, "Handler.Account.DeleteByPhone")
		defer finish()

		var input dto.DeleteByPhoneInput
		if err := ctx.ShouldBindQuery(&input); err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		output, err := c.DeleteByPhoneUseCase.Execute(ctx, input)
		if err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		c.Presenter.Present(ctx, output, http.StatusOK)
	}
}
