package account

import "github.com/google/uuid"

type CreateInput struct {
	Name            string `json:"name" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	Phone           string `json:"phone" validate:"required"`
	DocumentType    string `json:"document_type" validate:"required"`
	DocumentNumber  string `json:"document_number" validate:"required"`
	DocumentCountry string `json:"document_country" validate:"required"`
}

type CreateOutput struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	DocumentType    string    `json:"document_type"`
	DocumentNumber  string    `json:"document_number"`
	DocumentCountry string    `json:"document_country"`
}
