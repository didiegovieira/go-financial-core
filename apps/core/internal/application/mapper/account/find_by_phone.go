package account

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/normalizer"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
)

func MapFindByPhoneInputToAccount(dto account.FindByPhoneInput) *entity.Account {
	return &entity.Account{
		Phone: normalizer.Phone(dto.Phone),
	}
}

func MapAccountToFindByPhoneOutput(accountEntity *entity.Account) *account.FindByPhoneOutput {
	return &account.FindByPhoneOutput{
		Name:            accountEntity.Name,
		Email:           accountEntity.Email,
		Phone:           accountEntity.Phone,
		DocumentType:    accountEntity.DocumentType,
		DocumentNumber:  accountEntity.DocumentNumber,
		DocumentCountry: accountEntity.DocumentCountry,
	}
}
