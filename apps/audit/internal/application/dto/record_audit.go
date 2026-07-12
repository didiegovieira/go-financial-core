package dto

import "encoding/json"

type RecordAuditInput struct {
	EventName   string
	AggregateID string
	Payload     json.RawMessage
}
