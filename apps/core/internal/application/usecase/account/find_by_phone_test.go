package account

import (
	"context"
	"testing"

	dto "github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestFindByPhoneReturnsMappedAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountRepository := repository.NewMockAccountRepository(ctrl)
	usecase := NewFindByPhone(accountRepository)
	account := &entity.Account{
		ID:              uuid.New(),
		Name:            "Diego Vieira",
		Email:           "diego@example.com",
		Phone:           "5511999999999",
		DocumentType:    "cpf",
		DocumentNumber:  "12345678900",
		DocumentCountry: "BR",
	}

	accountRepository.EXPECT().
		FindByPhone(gomock.Any(), entity.Account{Phone: "5511999999999"}).
		Return(account, nil)

	output, err := usecase.Execute(context.Background(), dto.FindByPhoneInput{
		Phone: "+55 (11) 99999-9999",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != account.Name {
		t.Fatalf("expected name %q, got %q", account.Name, output.Name)
	}
	if output.Phone != account.Phone {
		t.Fatalf("expected phone %q, got %q", account.Phone, output.Phone)
	}
	if output.DocumentNumber != account.DocumentNumber {
		t.Fatalf("expected document number %q, got %q", account.DocumentNumber, output.DocumentNumber)
	}
}
