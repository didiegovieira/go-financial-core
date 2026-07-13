package wallet

import (
	"context"
	"net/http"
	"testing"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	appErrors "github.com/didiegovieira/go-financial-core/pkg/errors"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestDebitBlocksTransactionThatRequiresReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	walletRepository := repository.NewMockWalletRepository(ctrl)
	publisher := messaging.NewMockPublisher(ctrl)
	usecase := NewDebit(accountRepository, walletRepository, publisher)

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
}

func TestDebitBlocksRejectedTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	walletRepository := repository.NewMockWalletRepository(ctrl)
	publisher := messaging.NewMockPublisher(ctrl)
	usecase := NewDebit(accountRepository, walletRepository, publisher)

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
}

func TestDebitAllowsApprovedTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	walletRepository := repository.NewMockWalletRepository(ctrl)
	publisher := messaging.NewMockPublisher(ctrl)
	usecase := NewDebit(accountRepository, walletRepository, publisher)
	accountID := uuid.New()

	accountRepository.EXPECT().
		FindByPhone(gomock.Any(), entity.Account{Phone: "5511999999999"}).
		Return(&entity.Account{
			ID:    accountID,
			Phone: "5511999999999",
		}, nil)
	walletRepository.EXPECT().
		Debit(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error) {
			if wallet.AccountID != accountID {
				t.Fatalf("expected account ID %q, got %q", accountID, wallet.AccountID)
			}
			if wallet.Balance != 100_00 {
				t.Fatalf("expected debit amount %d, got %d", 100_00, wallet.Balance)
			}

			return &entity.Wallet{
				AccountID: wallet.AccountID,
				Balance:   900_00,
				Status:    "active",
			}, nil
		})
	publisher.EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, event messaging.Event) error {
			if event.Name != messaging.EventWalletDebited {
				t.Fatalf("expected event %q, got %q", messaging.EventWalletDebited, event.Name)
			}
			if event.Key != accountID.String() {
				t.Fatalf("expected event key %q, got %q", accountID.String(), event.Key)
			}
			if event.Payload == nil {
				t.Fatal("expected event payload")
			}

			return nil
		})

	output, err := usecase.Execute(context.Background(), dto.DebitInput{
		Phone:  "+55 (11) 99999-9999",
		Amount: 100_00,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Balance != 900_00 {
		t.Fatalf("expected balance %d, got %d", 900_00, output.Balance)
	}
}
