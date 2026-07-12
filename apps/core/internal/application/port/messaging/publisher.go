package messaging

import "context"

type EventName string

const (
	EventAccountCreated EventName = "account.created"
	EventAccountDeleted EventName = "account.deleted"
	EventWalletCredited EventName = "wallet.credited"
	EventWalletDebited  EventName = "wallet.debited"
)

type Event struct {
	Name    EventName
	Key     string
	Payload any
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}
