package column

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/shared/domain/column"
	"github.com/smarrog/task-board/shared/domain/shared"
)

type Column struct {
	id        Id
	boardId   board.Id
	position  Position
	createdAt time.Time
	updatedAt time.Time
	events    []shared.DomainEvent
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
	c.events = append(c.events, column.CreatedEvent{
		Id:      c.id.String(),
		BoardId: c.boardId.String(),
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

func (c *Column) Move(toPosition Position) {
	fromPosition := c.position

	c.position = toPosition

	c.events = append(c.events, column.MovedEvent{
		Id:           c.id.String(),
		FromPosition: fromPosition.Int(),
		ToPosition:   toPosition.Int(),
		At:           time.Now().UTC(),
	})
}

func (c *Column) PullEvents() []shared.DomainEvent {
	if len(c.events) == 0 {
		return nil
	}
	out := make([]shared.DomainEvent, len(c.events))
	copy(out, c.events)
	c.events = nil
	return out
}

type Tasks struct {
	columnId Id
}
