package wallet

import (
	"context"
	"net/http"
	"testing"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	appErrors "github.com/didiegovieira/go-financial-core/pkg/errors"
	"github.com/google/uuid"
)

type walletRepositoryStub struct {
	debitCalled bool
	debitWallet entity.Wallet
}

type accountRepositoryStub struct {
	account *entity.Account
}

func (r *accountRepositoryStub) Create(ctx context.Context, account entity.Account) error {
	return nil
}

func (r *accountRepositoryStub) FindByDocument(ctx context.Context, account entity.Account) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepositoryStub) FindByPhone(ctx context.Context, account entity.Account) (*entity.Account, error) {
	return r.account, nil
}

func (r *accountRepositoryStub) DeleteByPhone(ctx context.Context, account entity.Account) error {
	return nil
}

func (r *walletRepositoryStub) Create(ctx context.Context, wallet entity.Wallet) error {
	return nil
}

func (r *walletRepositoryStub) DeleteByAccountID(ctx context.Context, accountID uuid.UUID) error {
	return nil
}

func (r *walletRepositoryStub) Credit(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error) {
	return &wallet, nil
}

func (r *walletRepositoryStub) Debit(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error) {
	r.debitCalled = true
	r.debitWallet = wallet
	return &entity.Wallet{
		AccountID: wallet.AccountID,
		Balance:   900_00,
		Status:    "active",
	}, nil
}

type publisherStub struct {
	events []messaging.Event
}

func (p *publisherStub) Publish(ctx context.Context, event messaging.Event) error {
	p.events = append(p.events, event)
	return nil
}

func TestDebitBlocksTransactionThatRequiresReview(t *testing.T) {
	accountRepository := &accountRepositoryStub{}
	repository := &walletRepositoryStub{}
	publisher := &publisherStub{}
	usecase := NewDebit(accountRepository, repository, publisher)

	_, err := usecase.Execute(context.Background(), dto.DebitInput{
		Phone:  "+5511999999999",
		Amount: 500_00,
	})

	httpErr, ok := err.(*appErrors.Http)
	if !ok {
		t.Fatalf("expected http error, got %T", err)
	}
	if httpErr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, httpErr.Code)
	}
	if repository.debitCalled {
		t.Fatal("expected debit repository not to be called")
	}
	if len(publisher.events) != 0 {
		t.Fatal("expected no event to be published")
	}
}

func TestDebitBlocksRejectedTransaction(t *testing.T) {
	accountRepository := &accountRepositoryStub{}
	repository := &walletRepositoryStub{}
	publisher := &publisherStub{}
	usecase := NewDebit(accountRepository, repository, publisher)

	_, err := usecase.Execute(context.Background(), dto.DebitInput{
		Phone:  "+5511999999999",
		Amount: 5_000_00,
	})

	httpErr, ok := err.(*appErrors.Http)
	if !ok {
		t.Fatalf("expected http error, got %T", err)
	}
	if httpErr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, httpErr.Code)
	}
	if repository.debitCalled {
		t.Fatal("expected debit repository not to be called")
	}
	if len(publisher.events) != 0 {
		t.Fatal("expected no event to be published")
	}
}

func TestDebitAllowsApprovedTransaction(t *testing.T) {
	accountID := uuid.New()
	accountRepository := &accountRepositoryStub{
		account: &entity.Account{
			ID:    accountID,
			Phone: "+5511999999999",
		},
	}
	repository := &walletRepositoryStub{}
	publisher := &publisherStub{}
	usecase := NewDebit(accountRepository, repository, publisher)

	output, err := usecase.Execute(context.Background(), dto.DebitInput{
		Phone:  "+5511999999999",
		Amount: 100_00,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repository.debitCalled {
		t.Fatal("expected debit repository to be called")
	}
	if repository.debitWallet.AccountID != accountID {
		t.Fatalf("expected account ID %q, got %q", accountID, repository.debitWallet.AccountID)
	}
	if output.Balance != 900_00 {
		t.Fatalf("expected balance %d, got %d", 900_00, output.Balance)
	}
	if len(publisher.events) != 1 {
		t.Fatalf("expected 1 event to be published, got %d", len(publisher.events))
	}
	if publisher.events[0].Name != messaging.EventWalletDebited {
		t.Fatalf("expected event %q, got %q", messaging.EventWalletDebited, publisher.events[0].Name)
	}
	if publisher.events[0].Key != accountID.String() {
		t.Fatalf("expected event key %q, got %q", accountID.String(), publisher.events[0].Key)
	}
}
