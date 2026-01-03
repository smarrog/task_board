package board

import (
	"time"

	"github.com/google/uuid"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type Board struct {
	id          Id
	ownerId     common.UserId
	title       Title
	description Description
	createdAt   time.Time
	updatedAt   time.Time
	events      []common.DomainEvent
}

func New(ownerID common.UserId, title Title, description Description) (*Board, error) {
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
	b.events = append(b.events, CreatedEvent{
		Id:      b.id,
		OwnerId: b.ownerId,
		Title:   b.title,
		At:      now,
	})
	return b, nil
}

func Rehydrate(
	id Id,
	ownerId common.UserId,
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
func (b *Board) OwnerId() common.UserId   { return b.ownerId }
func (b *Board) Title() Title             { return b.title }
func (b *Board) Description() Description { return b.description }
func (b *Board) CreatedAt() time.Time     { return b.createdAt }
func (b *Board) UpdatedAt() time.Time     { return b.updatedAt }

func (b *Board) Update(title Title, description Description) {
	b.title = title
	b.description = description
	b.updatedAt = time.Now().UTC()
	b.events = append(b.events, UpdatedEvent{
		Id: b.id,
		At: b.updatedAt,
	})
}

func (b *Board) MarkDeleted() {
	at := time.Now().UTC()
	b.events = append(b.events, DeletedEvent{
		Id: b.id,
		At: at,
	})
}

func (b *Board) PullEvents() []common.DomainEvent {
	if len(b.events) == 0 {
		return nil
	}
	out := make([]common.DomainEvent, len(b.events))
	copy(out, b.events)
	b.events = nil
	return out
}

type Columns struct {
	boardId Id
}
