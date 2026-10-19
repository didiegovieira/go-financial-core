package account

import (
	"net/http"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/account"

	"github.com/didiegovieira/go-financial-core/pkg/api"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

type FindByPhone struct {
	Presenter          api.Presenter
	FindByPhoneUseCase account.FindByPhoneUseCase
}

// FindByPhone godoc
// @Summary      Find Account by Phone
// @Description  Find an existing account by phone number
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        phone  query     string  true  "Phone"
// @Success      200    {object}  account.FindByPhoneOutput "Account found successfully"
// @Failure      400    {object}  api.HttpError"
// @Router       /v1/financial/accounts [get]
func (c *FindByPhone) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := telemetry.TraceHandler(ctx, "Handler.Account.FindByPhone")
		defer finish()

		var input dto.FindByPhoneInput
		if err := ctx.ShouldBindQuery(&input); err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		output, err := c.FindByPhoneUseCase.Execute(ctx, input)
		if err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		c.Presenter.Present(ctx, output, http.StatusOK)
	}
}
