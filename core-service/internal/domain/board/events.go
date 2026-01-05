package board

import "time"

// Domain events are JSON-serializable (primitives + exported fields),
// so the outbox can marshal them directly.

type CreatedEvent struct {
	Id          string    `json:"id"`
	OwnerId     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

func (e CreatedEvent) Name() string          { return "BoardCreated" }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type UpdatedEvent struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

func (e UpdatedEvent) Name() string          { return "BoardUpdated" }
func (e UpdatedEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}

func (e DeletedEvent) Name() string          { return "BoardDeleted" }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
