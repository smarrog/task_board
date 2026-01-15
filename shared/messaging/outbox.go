package messaging

import (
	"encoding/json"
	"time"
)

type OutboxMessage struct {
	Id            string          `json:"id"`
	EventType     string          `json:"event_type"`
	AggregateType string          `json:"aggregate_type"`
	AggregateId   string          `json:"aggregate_id"`
	CreatedAt     time.Time       `json:"created_at"`
	Payload       json.RawMessage `json:"payload"`
	Version       int             `json:"version"`
}
