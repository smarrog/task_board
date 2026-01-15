package events

import "time"

// Contract payloads for task-related events.

type TaskCreated struct {
	Id          string    `json:"id"`
	ColumnId    string    `json:"column_id"`
	Position    int       `json:"position"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	At          time.Time `json:"at"`
}

type TaskUpdated struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	At          time.Time `json:"at"`
}

type TaskMoved struct {
	Id           string    `json:"id"`
	FromColumnId string    `json:"from_column_id"`
	ToColumnId   string    `json:"to_column_id"`
	FromPosition int       `json:"from_position"`
	ToPosition   int       `json:"to_position"`
	At          time.Time `json:"at"`
}

type TaskDeleted struct {
	Id string    `json:"id"`
	At time.Time `json:"at"`
}
