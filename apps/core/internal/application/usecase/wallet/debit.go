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
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/risk"
	"github.com/didiegovieira/go-financial-core/pkg/base"
	appErrors "github.com/didiegovieira/go-financial-core/pkg/errors"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type DebitUseCase = base.UseCase[dto.DebitInput, *dto.DebitOutput]

type Debit struct {
	AccountRepository repository.AccountRepository
	WalletRepository  repository.WalletRepository
	Publisher         messaging.Publisher
	RiskPolicy        risk.Policy
}

func NewDebit(
	accountRepository repository.AccountRepository,
	walletRepository repository.WalletRepository,
	publisher messaging.Publisher,
) *Debit {
	return &Debit{
		AccountRepository: accountRepository,
		WalletRepository:  walletRepository,
		Publisher:         publisher,
		RiskPolicy:        risk.NewDefaultPolicy(),
	}
}

func (c *Debit) Execute(ctx context.Context, input dto.DebitInput) (out *dto.DebitOutput, err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "UseCase.Wallet.Debit", &err)
	defer finish()

	decision := c.RiskPolicy.Evaluate(input.Amount)
	if decision.Status == risk.DecisionRejected {
		return nil, appErrors.HttpForbidden(decision.Reason)
	}
	if decision.Status == risk.DecisionReview {
		return nil, appErrors.HttpForbidden(decision.Reason)
	}

	phone := normalizer.Phone(input.Phone)
	account, err := c.AccountRepository.FindByPhone(ctx, entity.Account{Phone: phone})
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, appErrors.HttpBadRequest(fmt.Sprintf("account with phone number %s does not exist", phone))
	}

	walletEntity := mapper.MapDebitInputToWallet(input, account.ID)
	updatedWallet, err := c.WalletRepository.Debit(ctx, *walletEntity)
	if errors.Is(err, domainerror.ErrDebitNotAllowed) {
		return nil, appErrors.HttpBadRequest("wallet is inactive or has insufficient balance")
	}
	if err != nil {
		return nil, err
	}

	err = c.Publisher.Publish(ctx, messaging.Event{
		Name:    messaging.EventWalletDebited,
		Key:     account.ID.String(),
		Payload: event.NewWalletDebited(account.ID, input.Amount, updatedWallet.Balance, time.Now()),
	})
	if err != nil {
		return nil, err
	}

	return mapper.MapWalletToDebitOutput(updatedWallet), nil
}
