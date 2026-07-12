package account

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/dto/account"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/application/normalizer"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
)

func MapDeleteByPhoneInputToAccount(dto account.DeleteByPhoneInput) *entity.Account {
	return &entity.Account{
		Phone: normalizer.Phone(dto.Phone),
	}
}

func MapAccountToDeleteByPhoneOutput(accountEntity *entity.Account) *account.DeleteByPhoneOutput {
	return &account.DeleteByPhoneOutput{
		Success: true,
	}
}
