package column

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type Column struct {
	id        Id
	boardId   board.Id
	position  Position
	createdAt time.Time
	updatedAt time.Time
	events    []common.DomainEvent
}

func New(boardId board.Id, position Position) *Column {
	now := time.Now().UTC()
	c := &Column{
		id:        NewId(),
		boardId:   boardId,
		position:  position,
		createdAt: now,
		updatedAt: now,
	}
	c.events = append(c.events, CreatedEvent{
		Id:      c.id,
		BoardId: c.boardId,
		At:      c.createdAt,
	})
	return c
}

func Rehydrate(
	id Id,
	boardId board.Id,
	position Position,
	createdAt time.Time,
	updatedAt time.Time,
) *Column {
	return &Column{
		id:        id,
		boardId:   boardId,
		position:  position,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *Column) Id() Id               { return c.id }
func (c *Column) BoardId() board.Id    { return c.boardId }
func (c *Column) Position() Position   { return c.position }
func (c *Column) CreatedAt() time.Time { return c.createdAt }
func (c *Column) UpdatedAt() time.Time { return c.updatedAt }

func (c *Column) Update(boardId board.Id, position Position) {
	c.boardId = boardId
	c.position = position
	c.events = append(c.events, MoveEvent{
		Id:       c.id,
		BoardId:  c.boardId,
		Position: position,
	})
}

func (c *Column) PullEvents() []common.DomainEvent {
	if len(c.events) == 0 {
		return nil
	}
	out := make([]common.DomainEvent, len(c.events))
	copy(out, c.events)
	c.events = nil
	return out
}
