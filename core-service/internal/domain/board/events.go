package board

import "time"

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type BoardCreated struct {
	BoardID BoardId
	OwnerID UserId
	Title   Title
	At      time.Time
}

func (e BoardCreated) EventName() string     { return "BoardCreated" }
func (e BoardCreated) OccurredAt() time.Time { return e.At }

type BoardUpdated struct {
	BoardID BoardId
	At      time.Time
}

func (e BoardUpdated) EventName() string     { return "BoardUpdated" }
func (e BoardUpdated) OccurredAt() time.Time { return e.At }

type BoardDeleted struct {
	BoardID BoardId
	At      time.Time
}

func (e BoardDeleted) EventName() string     { return "BoardDeleted" }
func (e BoardDeleted) OccurredAt() time.Time { return e.At }
