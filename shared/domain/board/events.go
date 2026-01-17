package board

import (
	"time"
)

const (
	EvtCreated = "BoardCreated"
	EvtUpdated = "BoardUpdated"
	EvtDeleted = "BoardDeleted"
)

type CreatedEvent struct {
	Id          string    `json:"id"`
	OwnerId     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return EvtCreated }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type UpdatedEvent struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
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
