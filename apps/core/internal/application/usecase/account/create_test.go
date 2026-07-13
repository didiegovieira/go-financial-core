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

func TestCreateCreatesAccountWalletAndPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	walletRepository := repository.NewMockWalletRepository(ctrl)
	publisher := messaging.NewMockPublisher(ctrl)
	usecase := NewCreate(accountRepository, walletRepository, publisher)
	input := dto.CreateInput{
		Name:            "Diego Vieira",
		Email:           "diego@example.com",
		Password:        "secret123",
		Phone:           "+55 (11) 99999-9999",
		DocumentType:    "cpf",
		DocumentNumber:  "12345678900",
		DocumentCountry: "BR",
	}
	var createdAccountID uuid.UUID

	accountRepository.EXPECT().
		FindByDocument(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, account entity.Account) (*entity.Account, error) {
			if account.Phone != "5511999999999" {
				t.Fatalf("expected normalized phone %q, got %q", "5511999999999", account.Phone)
			}
			if account.DocumentNumber != input.DocumentNumber {
				t.Fatalf("expected document number %q, got %q", input.DocumentNumber, account.DocumentNumber)
			}

			return nil, nil
		})
	accountRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, account entity.Account) error {
			if account.ID == uuid.Nil {
				t.Fatal("expected account ID to be generated")
			}
			if account.Name != input.Name {
				t.Fatalf("expected name %q, got %q", input.Name, account.Name)
			}
			if account.Phone != "5511999999999" {
				t.Fatalf("expected normalized phone %q, got %q", "5511999999999", account.Phone)
			}

			createdAccountID = account.ID
			return nil
		})
	walletRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, wallet entity.Wallet) error {
			if wallet.ID == uuid.Nil {
				t.Fatal("expected wallet ID to be generated")
			}
			if wallet.AccountID != createdAccountID {
				t.Fatalf("expected account ID %q, got %q", createdAccountID, wallet.AccountID)
			}
			if wallet.Balance != 0 {
				t.Fatalf("expected initial balance %d, got %d", 0, wallet.Balance)
			}
			if wallet.Status != "active" {
				t.Fatalf("expected status %q, got %q", "active", wallet.Status)
			}

			return nil
		})
	publisher.EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, event messaging.Event) error {
			if event.Name != messaging.EventAccountCreated {
				t.Fatalf("expected event %q, got %q", messaging.EventAccountCreated, event.Name)
			}
			if event.Key != createdAccountID.String() {
				t.Fatalf("expected event key %q, got %q", createdAccountID.String(), event.Key)
			}
			if event.Payload == nil {
				t.Fatal("expected event payload")
			}

			return nil
		})

	output, err := usecase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ID != createdAccountID {
		t.Fatalf("expected output ID %q, got %q", createdAccountID, output.ID)
	}
	if output.Phone != "5511999999999" {
		t.Fatalf("expected output phone %q, got %q", "5511999999999", output.Phone)
	}
}
