package mapper

import (
	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/model"
)

func AccountToEntity(u model.AccountModel) *entity.Account {
	return &entity.Account{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Password:        u.Password,
		Phone:           u.Phone,
		DocumentType:    u.DocumentType,
		DocumentNumber:  u.DocumentNumber,
		DocumentCountry: u.DocumentCountry,
	}
}

func AccountToModel(u entity.Account) model.AccountModel {
	return model.AccountModel{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Password:        u.Password,
		Phone:           u.Phone,
		DocumentType:    u.DocumentType,
		DocumentNumber:  u.DocumentNumber,
		DocumentCountry: u.DocumentCountry,
	}
}
