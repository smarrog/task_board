package board

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type CreatedEvent struct {
	Id          Id
	OwnerId     common.UserId
	Title       common.Title
	Description common.Description
	At          time.Time
}

func (e CreatedEvent) Name() string          { return "BoardCreated" }
func (e CreatedEvent) OccurredAt() time.Time { return e.At }

type UpdatedEvent struct {
	Id Id
	At time.Time
}

func (e UpdatedEvent) Name() string          { return "BoardUpdated" }
func (e UpdatedEvent) OccurredAt() time.Time { return e.At }

type DeletedEvent struct {
	Id Id
	At time.Time
}

func (e DeletedEvent) Name() string          { return "BoardDeleted" }
func (e DeletedEvent) OccurredAt() time.Time { return e.At }
