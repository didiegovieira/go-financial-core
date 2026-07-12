package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLogModel struct {
	ID          uuid.UUID
	EventName   string
	AggregateID string
	Payload     json.RawMessage
	OccurredAt  time.Time
	CreatedAt   time.Time
}
