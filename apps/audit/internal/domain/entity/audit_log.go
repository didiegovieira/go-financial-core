package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID
	EventName   string
	AggregateID string
	Payload     json.RawMessage
	OccurredAt  time.Time
}
