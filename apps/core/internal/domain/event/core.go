package event

import (
	"time"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/google/uuid"
)

const Version = 1

type Envelope[T any] struct {
	EventID       uuid.UUID `json:"event_id"`
	EventName     string    `json:"event_name"`
	Version       int       `json:"version"`
	OccurredAt    time.Time `json:"occurred_at"`
	AggregateType string    `json:"aggregate_type"`
	AggregateID   uuid.UUID `json:"aggregate_id"`
	Payload       T         `json:"payload"`
}

type AccountCreatedPayload struct {
	AccountID       uuid.UUID `json:"account_id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	DocumentType    string    `json:"document_type"`
	DocumentCountry string    `json:"document_country"`
}

type AccountDeletedPayload struct {
	AccountID uuid.UUID `json:"account_id"`
	Phone     string    `json:"phone"`
}

type WalletCreditedPayload struct {
	AccountID uuid.UUID `json:"account_id"`
	Amount    int64     `json:"amount"`
	Balance   int64     `json:"balance"`
	Operation string    `json:"operation"`
}

type WalletDebitedPayload struct {
	AccountID uuid.UUID `json:"account_id"`
	Amount    int64     `json:"amount"`
	Balance   int64     `json:"balance"`
	Operation string    `json:"operation"`
}

func NewAccountCreated(account entity.Account, occurredAt time.Time) Envelope[AccountCreatedPayload] {
	return Envelope[AccountCreatedPayload]{
		EventID:       uuid.New(),
		EventName:     "account.created",
		Version:       Version,
		OccurredAt:    occurredAt.UTC(),
		AggregateType: "account",
		AggregateID:   account.ID,
		Payload: AccountCreatedPayload{
			AccountID:       account.ID,
			Name:            account.Name,
			Email:           account.Email,
			Phone:           account.Phone,
			DocumentType:    account.DocumentType,
			DocumentCountry: account.DocumentCountry,
		},
	}
}

func NewAccountDeleted(account entity.Account, occurredAt time.Time) Envelope[AccountDeletedPayload] {
	return Envelope[AccountDeletedPayload]{
		EventID:       uuid.New(),
		EventName:     "account.deleted",
		Version:       Version,
		OccurredAt:    occurredAt.UTC(),
		AggregateType: "account",
		AggregateID:   account.ID,
		Payload: AccountDeletedPayload{
			AccountID: account.ID,
			Phone:     account.Phone,
		},
	}
}

func NewWalletCredited(accountID uuid.UUID, amount int64, balance int64, occurredAt time.Time) Envelope[WalletCreditedPayload] {
	return Envelope[WalletCreditedPayload]{
		EventID:       uuid.New(),
		EventName:     "wallet.credited",
		Version:       Version,
		OccurredAt:    occurredAt.UTC(),
		AggregateType: "wallet",
		AggregateID:   accountID,
		Payload: WalletCreditedPayload{
			AccountID: accountID,
			Amount:    amount,
			Balance:   balance,
			Operation: "credit",
		},
	}
}

func NewWalletDebited(accountID uuid.UUID, amount int64, balance int64, occurredAt time.Time) Envelope[WalletDebitedPayload] {
	return Envelope[WalletDebitedPayload]{
		EventID:       uuid.New(),
		EventName:     "wallet.debited",
		Version:       Version,
		OccurredAt:    occurredAt.UTC(),
		AggregateType: "wallet",
		AggregateID:   accountID,
		Payload: WalletDebitedPayload{
			AccountID: accountID,
			Amount:    amount,
			Balance:   balance,
			Operation: "debit",
		},
	}
}
