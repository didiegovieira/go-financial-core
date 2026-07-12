package account

import (
	"context"
	"fmt"
	"time"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	mapper "github.com/didiegovieira/go-financial-core/apps/core/internal/application/mapper/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/event"
	"github.com/didiegovieira/go-financial-core/pkg/base"
	"github.com/didiegovieira/go-financial-core/pkg/metrics"
)

type CreateUseCase = base.UseCase[dto.CreateInput, *dto.CreateOutput]

type Create struct {
	AccountRepository repository.AccountRepository
	WalletRepository  repository.WalletRepository
	Publisher         messaging.Publisher
}

func NewCreate(
	AccountRepository repository.AccountRepository,
	WalletRepository repository.WalletRepository,
	Publisher messaging.Publisher,
) *Create {
	return &Create{
		AccountRepository: AccountRepository,
		WalletRepository:  WalletRepository,
		Publisher:         Publisher,
	}
}

func (c *Create) Execute(ctx context.Context, input dto.CreateInput) (out *dto.CreateOutput, err error) {
	ctx, finish := metrics.TraceWithError(ctx, "UseCase.Account.Create", &err)
	defer finish()

	accountEntity := mapper.MapCreateInputToAccount(input)

	existingAccount, err := c.AccountRepository.FindByDocument(ctx, *accountEntity)
	if err != nil {
		return nil, err
	}

	if existingAccount != nil {
		return out, fmt.Errorf("account with document number %s already exists", accountEntity.DocumentNumber)
	}

	accountEntity.NewAccount()
	err = c.AccountRepository.Create(ctx, *accountEntity)
	if err != nil {
		return out, err
	}

	walletEntity := entity.NewWallet(accountEntity.ID)
	err = c.WalletRepository.Create(ctx, *walletEntity)
	if err != nil {
		return out, err
	}

	err = c.Publisher.Publish(ctx, messaging.Event{
		Name:    messaging.EventAccountCreated,
		Key:     accountEntity.ID.String(),
		Payload: event.NewAccountCreated(*accountEntity, time.Now()),
	})
	if err != nil {
		return out, err
	}

	out = mapper.MapAccountToCreateOutput(accountEntity)

	return out, nil
}
