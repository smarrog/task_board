package notification

import (
	"context"
	"time"
)

type HistoryRecord struct {
	OutboxId       string
	EventType      string
	AggregateType  string
	AggregateId    string
	EventCreatedAt time.Time
	Version        int
	Payload        []byte
	Text           string
}

type HistoryRepository interface {
	Save(ctx context.Context, r HistoryRecord) error
}
