package task

import (
	"time"
)

const (
	EvtCreated = "TaskCreated"
	EvtUpdated = "TaskUpdated"
	EvtMoved   = "TaskMoved"
	EvtDeleted = "TaskDeleted"
)

type CreatedEvent struct {
	Id          string    `json:"id"`
	ColumnId    string    `json:"column_id"`
	Position    int       `json:"position"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	At          time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return EvtCreated }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MovedEvent struct {
	Id           string    `json:"id"`
	FromColumnId string    `json:"from_column_id"`
	ToColumnId   string    `json:"to_column_id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At           time.Time `json:"at"`
}

func (e MovedEvent) Name() string          { return EvtMoved }
func (e MovedEvent) OccurredAt() time.Time { return e.At }

type UpdatedEvent struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	At          time.Time `json:"at"`
}

func (e UpdatedEvent) Name() string          { return EvtUpdated }
func (e UpdatedEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}

func (e DeletedEvent) Name() string          { return EvtDeleted }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
