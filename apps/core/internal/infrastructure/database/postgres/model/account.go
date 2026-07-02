package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountModel struct {
	ID              uuid.UUID
	Name            string
	Email           string
	Password        string
	Phone           string
	DocumentType    string
	DocumentNumber  string
	DocumentCountry string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
