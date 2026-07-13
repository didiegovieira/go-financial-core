package wallet

import (
	"context"
	"testing"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/wallet"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestCreditAddsBalanceAndPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	walletRepository := repository.NewMockWalletRepository(ctrl)
	publisher := messaging.NewMockPublisher(ctrl)
	usecase := NewCredit(accountRepository, walletRepository, publisher)
	accountID := uuid.New()

	accountRepository.EXPECT().
		FindByPhone(gomock.Any(), entity.Account{Phone: "5511999999999"}).
		Return(&entity.Account{
			ID:    accountID,
			Phone: "5511999999999",
		}, nil)
	walletRepository.EXPECT().
		Credit(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error) {
			if wallet.AccountID != accountID {
				t.Fatalf("expected account ID %q, got %q", accountID, wallet.AccountID)
			}
			if wallet.Balance != 100_00 {
				t.Fatalf("expected credit amount %d, got %d", 100_00, wallet.Balance)
			}

			return &entity.Wallet{
				AccountID: wallet.AccountID,
				Balance:   1_100_00,
				Status:    "active",
			}, nil
		})
	publisher.EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, event messaging.Event) error {
			if event.Name != messaging.EventWalletCredited {
				t.Fatalf("expected event %q, got %q", messaging.EventWalletCredited, event.Name)
			}
			if event.Key != accountID.String() {
				t.Fatalf("expected event key %q, got %q", accountID.String(), event.Key)
			}
			if event.Payload == nil {
				t.Fatal("expected event payload")
			}

			return nil
		})

	output, err := usecase.Execute(context.Background(), dto.CreditInput{
		Phone:  "+55 (11) 99999-9999",
		Amount: 100_00,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Balance != 1_100_00 {
		t.Fatalf("expected balance %d, got %d", 1_100_00, output.Balance)
	}
}
