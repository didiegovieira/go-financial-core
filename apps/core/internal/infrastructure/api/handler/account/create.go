package account

import (
	"net/http"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/usecase/account"

	"github.com/didiegovieira/go-financial-core/pkg/api"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
	"github.com/gin-gonic/gin"
)

type Create struct {
	Presenter     api.Presenter
	CreateUseCase account.CreateUseCase
}

// Create godoc
// @Summary      Create Account
// @Description  Create a new account
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Success      201  {object}  dto.CreateOutput "Account created successfully"
// @Router       /v1/financial/accounts [post]
func (c *Create) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		finish := metrics.TraceHandler(ctx, "Handler.Account.Create")
		defer finish()

		var input dto.CreateInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		output, err := c.CreateUseCase.Execute(ctx, input)
		if err != nil {
			c.Presenter.Error(ctx, err)
			return
		}

		c.Presenter.Present(ctx, output, http.StatusCreated)
	}
}
