package events

import "time"

// Contract payloads for column-related events.

type ColumnCreated struct {
	Id      string    `json:"id"`
	BoardId string    `json:"board_id"`
	At      time.Time `json:"at"`
}

type ColumnMoved struct {
	Id           string    `json:"id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At           time.Time `json:"at"`
}

type ColumnDeleted struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}
