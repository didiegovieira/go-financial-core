package dto

import "encoding/json"

type SendNotificationInput struct {
	EventName   string
	AggregateID string
	Payload     json.RawMessage
}
