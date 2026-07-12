package account

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/normalizer"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
)

func MapCreateInputToAccount(dto account.CreateInput) *entity.Account {
	return &entity.Account{
		Name:            dto.Name,
		Email:           dto.Email,
		Password:        dto.Password,
		Phone:           normalizer.Phone(dto.Phone),
		DocumentType:    dto.DocumentType,
		DocumentNumber:  dto.DocumentNumber,
		DocumentCountry: dto.DocumentCountry,
	}
}

func MapAccountToCreateOutput(accountEntity *entity.Account) *account.CreateOutput {
	return &account.CreateOutput{
		ID:              accountEntity.ID,
		Name:            accountEntity.Name,
		Email:           accountEntity.Email,
		Phone:           accountEntity.Phone,
		DocumentType:    accountEntity.DocumentType,
		DocumentNumber:  accountEntity.DocumentNumber,
		DocumentCountry: accountEntity.DocumentCountry,
	}
}
