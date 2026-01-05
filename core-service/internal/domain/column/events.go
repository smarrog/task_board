package column

import "time"

// Domain events are JSON-serializable (primitives + exported fields),
// so the outbox can marshal them directly.

type CreatedEvent struct {
	Id      string    `json:"id"`
	BoardId string    `json:"board_id"`
	At      time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return "ColumnCreated" }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MoveEvent struct {
	Id           string    `json:"id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At           time.Time `json:"at"`
}

func (e MoveEvent) Name() string          { return "ColumnMoved" }
func (e MoveEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}

func (e DeletedEvent) Name() string          { return "ColumnDeleted" }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
