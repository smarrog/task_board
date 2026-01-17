package board

import (
	"time"

	"github.com/google/uuid"
	"github.com/smarrog/task-board/shared/domain/board"
	"github.com/smarrog/task-board/shared/domain/shared"
)

type Board struct {
	id          Id
	ownerId     shared.UserId
	title       Title
	description Description
	createdAt   time.Time
	updatedAt   time.Time
	events      []shared.DomainEvent
}

func New(ownerID shared.UserId, title Title, description Description) (*Board, error) {
	if ownerID.UUID() == uuid.Nil {
		return nil, ErrOwnerRequired
	}

	now := time.Now().UTC()
	b := &Board{
		id:          NewId(),
		ownerId:     ownerID,
		title:       title,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}
	b.events = append(b.events, board.CreatedEvent{
		Id:          b.id.String(),
		OwnerId:     b.ownerId.String(),
		Title:       b.title.String(),
		Description: b.description.String(),
		At:          now,
	})
	return b, nil
}

func Rehydrate(
	id Id,
	ownerId shared.UserId,
	title Title,
	description Description,
	createdAt time.Time,
	updatedAt time.Time,
) *Board {
	return &Board{
		id:          id,
		ownerId:     ownerId,
		title:       title,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *Board) Id() Id                   { return b.id }
func (b *Board) OwnerId() shared.UserId   { return b.ownerId }
func (b *Board) Title() Title             { return b.title }
func (b *Board) Description() Description { return b.description }
func (b *Board) CreatedAt() time.Time     { return b.createdAt }
func (b *Board) UpdatedAt() time.Time     { return b.updatedAt }

func (b *Board) Update(title Title, description Description) {
	b.title = title
	b.description = description
	b.updatedAt = time.Now().UTC()
	b.events = append(b.events, board.UpdatedEvent{
		Id:          b.id.String(),
		Title:       b.title.String(),
		Description: b.description.String(),
		At:          b.updatedAt,
	})
}

func (b *Board) MarkDeleted() {
	at := time.Now().UTC()
	b.events = append(b.events, board.DeletedEvent{
		Id: b.id.String(),
		At: at,
	})
}

func (b *Board) PullEvents() []shared.DomainEvent {
	if len(b.events) == 0 {
		return nil
	}
	out := make([]shared.DomainEvent, len(b.events))
	copy(out, b.events)
	b.events = nil
	return out
}

type Columns struct {
	boardId Id
}
