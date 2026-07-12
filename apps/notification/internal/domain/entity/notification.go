package entity

import "time"

type Notification struct {
	ID          string
	EventName   string
	AggregateID string
	Channel     string
	Recipient   string
	Subject     string
	Body        string
	Status      string
	CreatedAt   time.Time
}
