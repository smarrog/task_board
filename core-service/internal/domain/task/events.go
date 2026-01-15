package task

import (
	"time"

	"github.com/smarrog/task-board/shared/messaging"
)

// IMPORTANT:
// Domain events are stored in a JSON-friendly shape (primitives + exported fields),
// so that outbox can serialize them without extra DTO mapping.

type CreatedEvent struct {
	Id          string    `json:"id"`
	ColumnId    string    `json:"column_id"`
	Position    int       `json:"position"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	At          time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return messaging.EvtTaskCreated }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MoveEvent struct {
	Id           string    `json:"id"`
	FromColumnId string    `json:"from_column_id"`
	ToColumnId   string    `json:"to_column_id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At           time.Time `json:"at"`
}

func (e MoveEvent) Name() string          { return messaging.EvtTaskMoved }
func (e MoveEvent) OccurredAt() time.Time { return e.At }

type UpdatedEvent struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	At          time.Time `json:"at"`
}

func (e UpdatedEvent) Name() string          { return messaging.EvtTaskUpdated }
func (e UpdatedEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}

func (e DeletedEvent) Name() string          { return messaging.EvtTaskDeleted }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
