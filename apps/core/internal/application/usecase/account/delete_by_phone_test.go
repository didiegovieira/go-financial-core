package account

import (
	"context"
	"testing"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestDeleteByPhoneDeletesWalletAccountAndPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	walletRepository := repository.NewMockWalletRepository(ctrl)
	publisher := messaging.NewMockPublisher(ctrl)
	usecase := NewDeleteByPhone(accountRepository, walletRepository, publisher)
	accountID := uuid.New()
	account := &entity.Account{
		ID:    accountID,
		Phone: "5511999999999",
	}

	accountRepository.EXPECT().
		FindByPhone(gomock.Any(), entity.Account{Phone: "5511999999999"}).
		Return(account, nil)
	walletRepository.EXPECT().
		DeleteByAccountID(gomock.Any(), accountID).
		Return(nil)
	accountRepository.EXPECT().
		DeleteByPhone(gomock.Any(), entity.Account{Phone: "5511999999999"}).
		Return(nil)
	publisher.EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, event messaging.Event) error {
			if event.Name != messaging.EventAccountDeleted {
				t.Fatalf("expected event %q, got %q", messaging.EventAccountDeleted, event.Name)
			}
			if event.Key != accountID.String() {
				t.Fatalf("expected event key %q, got %q", accountID.String(), event.Key)
			}
			if event.Payload == nil {
				t.Fatal("expected event payload")
			}

			return nil
		})

	output, err := usecase.Execute(context.Background(), dto.DeleteByPhoneInput{
		Phone: "+55 (11) 99999-9999",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !output.Success {
		t.Fatal("expected success output")
	}
}
