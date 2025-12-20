package model

import (
	"github.com/google/uuid"
)

type Board struct {
	Id          uuid.UUID
	OwnerId     uuid.UUID
	Title       string
	Description string
}

type Column struct {
	Id       uuid.UUID
	BoardId  uuid.UUID
	Position int
}

type Task struct {
	Id          uuid.UUID
	ColumnId    uuid.UUID
	Position    int
	Title       string
	Description string
	AssigneeId  uuid.UUID
}
