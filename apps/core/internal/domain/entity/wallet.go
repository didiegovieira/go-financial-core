package entity

import (
	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Balance   int64
	Status    string
}

func NewWallet(accountID uuid.UUID) *Wallet {
	return &Wallet{
		ID:        uuid.New(),
		AccountID: accountID,
		Balance:   0,
		Status:    "active",
	}
}
