package column

import (
	"time"
)

const (
	EvtCreated = "ColumnCreated"
	EvtMoved   = "ColumnMoved"
	EvtDeleted = "ColumnDeleted"
)

type CreatedEvent struct {
	Id      string    `json:"id"`
	BoardId string    `json:"board_id"`
	At      time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return EvtCreated }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MovedEvent struct {
	Id           string    `json:"id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At           time.Time `json:"at"`
}

func (e MovedEvent) Name() string          { return EvtMoved }
func (e MovedEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}

func (e DeletedEvent) Name() string          { return EvtDeleted }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
