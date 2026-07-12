package account

import (
	"context"
	"fmt"
	"time"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	mapper "github.com/didiegovieira/go-financial-core/apps/core/internal/application/mapper/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/event"
	"github.com/didiegovieira/go-financial-core/pkg/base"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
)

type DeleteByPhoneUseCase = base.UseCase[dto.DeleteByPhoneInput, *dto.DeleteByPhoneOutput]

type DeleteByPhone struct {
	AccountRepository repository.AccountRepository
	WalletRepository  repository.WalletRepository
	Publisher         messaging.Publisher
}

func NewDeleteByPhone(
	AccountRepository repository.AccountRepository,
	WalletRepository repository.WalletRepository,
	Publisher messaging.Publisher,
) *DeleteByPhone {
	return &DeleteByPhone{
		AccountRepository: AccountRepository,
		WalletRepository:  WalletRepository,
		Publisher:         Publisher,
	}
}

func (c *DeleteByPhone) Execute(ctx context.Context, input dto.DeleteByPhoneInput) (out *dto.DeleteByPhoneOutput, err error) {
	ctx, finish := metrics.TraceWithError(ctx, "UseCase.Account.DeleteByPhone", &err)
	defer finish()

	accountEntity := mapper.MapDeleteByPhoneInputToAccount(input)

	existingAccount, err := c.AccountRepository.FindByPhone(ctx, *accountEntity)
	if err != nil {
		return nil, err
	}
	if existingAccount == nil {
		return out, fmt.Errorf("account with phone number %s does not exist", accountEntity.Phone)
	}

	err = c.WalletRepository.DeleteByAccountID(ctx, existingAccount.ID)
	if err != nil {
		return nil, err
	}

	err = c.AccountRepository.DeleteByPhone(ctx, *accountEntity)
	if err != nil {
		return nil, err
	}

	err = c.Publisher.Publish(ctx, messaging.Event{
		Name:    messaging.EventAccountDeleted,
		Key:     existingAccount.ID.String(),
		Payload: event.NewAccountDeleted(*existingAccount, time.Now()),
	})
	if err != nil {
		return nil, err
	}

	out = mapper.MapAccountToDeleteByPhoneOutput(existingAccount)
	return out, nil
}
