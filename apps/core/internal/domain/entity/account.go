package entity

import (
	"github.com/google/uuid"
)

type Account struct {
	ID              uuid.UUID
	Name            string
	Email           string
	Password        string
	Phone           string
	DocumentType    string
	DocumentNumber  string
	DocumentCountry string
}

func (a *Account) NewAccount() *Account {
	a.ID = uuid.New()
	return a
}
