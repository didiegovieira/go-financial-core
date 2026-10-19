package account

import (
	"context"
	"fmt"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	mapper "github.com/didiegovieira/go-financial-core/apps/core/internal/application/mapper/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/pkg/base"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type FindByPhoneUseCase = base.UseCase[dto.FindByPhoneInput, *dto.FindByPhoneOutput]

type FindByPhone struct {
	AccountRepository repository.AccountRepository
}

func NewFindByPhone(
	AccountRepository repository.AccountRepository,
) *FindByPhone {
	return &FindByPhone{
		AccountRepository: AccountRepository,
	}
}

func (c *FindByPhone) Execute(ctx context.Context, input dto.FindByPhoneInput) (out *dto.FindByPhoneOutput, err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "UseCase.Account.FindByPhone", &err)
	defer finish()

	accountEntity := mapper.MapFindByPhoneInputToAccount(input)
	existingAccount, err := c.AccountRepository.FindByPhone(ctx, *accountEntity)
	if err != nil {
		return nil, err
	}
	if existingAccount == nil {
		return out, fmt.Errorf("account with phone number %s does not exist", accountEntity.Phone)
	}

	out = mapper.MapAccountToFindByPhoneOutput(existingAccount)

	return out, nil
}
