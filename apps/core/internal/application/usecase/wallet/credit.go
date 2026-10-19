package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	mapper "github.com/didiegovieira/go-financial-core/apps/core/internal/application/mapper/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/normalizer"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	domainerror "github.com/didiegovieira/go-financial-core/apps/core/internal/domain/error"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/event"
	"github.com/didiegovieira/go-financial-core/pkg/base"
	appErrors "github.com/didiegovieira/go-financial-core/pkg/errors"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type CreditUseCase = base.UseCase[dto.CreditInput, *dto.CreditOutput]

type Credit struct {
	AccountRepository repository.AccountRepository
	WalletRepository  repository.WalletRepository
	Publisher         messaging.Publisher
}

func NewCredit(
	accountRepository repository.AccountRepository,
	walletRepository repository.WalletRepository,
	publisher messaging.Publisher,
) *Credit {
	return &Credit{
		AccountRepository: accountRepository,
		WalletRepository:  walletRepository,
		Publisher:         publisher,
	}
}

func (c *Credit) Execute(ctx context.Context, input dto.CreditInput) (out *dto.CreditOutput, err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "UseCase.Wallet.Credit", &err)
	defer finish()

	phone := normalizer.Phone(input.Phone)
	account, err := c.AccountRepository.FindByPhone(ctx, entity.Account{Phone: phone})
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, appErrors.HttpBadRequest(fmt.Sprintf("account with phone number %s does not exist", phone))
	}

	walletEntity := mapper.MapCreditInputToWallet(input, account.ID)
	updatedWallet, err := c.WalletRepository.Credit(ctx, *walletEntity)
	if errors.Is(err, domainerror.ErrWalletUnavailable) {
		return nil, appErrors.HttpBadRequest("wallet does not exist or is inactive")
	}
	if err != nil {
		return nil, err
	}

	err = c.Publisher.Publish(ctx, messaging.Event{
		Name:    messaging.EventWalletCredited,
		Key:     account.ID.String(),
		Payload: event.NewWalletCredited(account.ID, input.Amount, updatedWallet.Balance, time.Now()),
	})
	if err != nil {
		return nil, err
	}

	return mapper.MapWalletToCreditOutput(updatedWallet), nil
}
