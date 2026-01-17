package outbox

import (
	"encoding/json"
	"time"
)

type Message struct {
	Id            string          `json:"id"`
	EventType     string          `json:"event_type"`
	AggregateType string          `json:"aggregate_type"`
	AggregateId   string          `json:"aggregate_id"`
	CreatedAt     time.Time       `json:"created_at"`
	Payload       json.RawMessage `json:"payload"`
	Version       int             `json:"version"`
}
