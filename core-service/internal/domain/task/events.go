package task

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type CreatedEvent struct {
	Id          Id
	ColumnId    column.Id
	Position    Position
	Title       Title
	Description Description
	AssigneeId  common.UserId
	At          time.Time
}

func (e CreatedEvent) Name() string          { return "TaskCreated" }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MoveEvent struct {
	Id       Id
	ColumnId column.Id
	Position int
	At       time.Time
}

func (e MoveEvent) Name() string          { return "TaskMoved" }
func (e MoveEvent) OccurredAt() time.Time { return e.At }

type UpdatedEvent struct {
	Id          Id
	Title       Title
	Description Description
	AssigneeId  common.UserId
	At          time.Time
}

func (e UpdatedEvent) Name() string          { return "TaskUpdated" }
func (e UpdatedEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id Id
	At time.Time
}

func (e DeletedEvent) Name() string          { return "TaskDeleted" }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
