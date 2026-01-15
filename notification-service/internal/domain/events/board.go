package events

import "time"

// NOTE: These structs are *contracts* for payloads produced by core-service domain events
// and stored in the outbox as JSON. Keep JSON tags aligned with core-service.

type BoardCreated struct {
	Id          string    `json:"id"`
	OwnerId     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

type BoardUpdated struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

type BoardDeleted struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}
