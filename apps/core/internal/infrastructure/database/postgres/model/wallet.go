package model

import (
	"time"

	"github.com/google/uuid"
)

type WalletModel struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Balance   float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
