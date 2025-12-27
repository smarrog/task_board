package column

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type CreatedEvent struct {
	Id      Id
	BoardId board.Id
	At      time.Time
}

func (e CreatedEvent) Name() string          { return "ColumnCreated" }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type MoveEvent struct {
	Id       Id
	BoardId  board.Id
	Position int
	At       time.Time
}

func (e MoveEvent) Name() string          { return "ColumnMoved" }
func (e MoveEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id Id
	At time.Time
}

func (e DeletedEvent) Name() string          { return "ColumnDeleted" }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
