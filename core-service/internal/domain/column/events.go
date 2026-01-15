package column

import (
	"time"

	"github.com/smarrog/task-board/shared/messaging"
)

// Domain events are JSON-serializable (primitives + exported fields),
// so the outbox can marshal them directly.

type CreatedEvent struct {
	Id      string    `json:"id"`
	BoardId string    `json:"board_id"`
	At      time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return messaging.EvtColumnCreated }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MoveEvent struct {
	Id           string    `json:"id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At           time.Time `json:"at"`
}

func (e MoveEvent) Name() string          { return messaging.EvtColumnMoved }
func (e MoveEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}

func (e DeletedEvent) Name() string          { return messaging.EvtColumnDeleted }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
