package board

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	id          BoardId
	ownerId     UserId
	title       Title
	description Description
	createdAt   time.Time
	updatedAt   time.Time
	events      []DomainEvent
}

func NewBoard(ownerID UserId, title Title, description Description) (*Board, error) {
	if ownerID.UUID() == uuid.Nil {
		return nil, ErrOwnerRequired
	}

	now := time.Now().UTC()
	b := &Board{
		id:          NewBoardID(),
		ownerId:     ownerID,
		title:       title,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}
	b.events = append(b.events, BoardCreated{BoardID: b.id, OwnerID: b.ownerId, Title: b.title, At: now})
	return b, nil
}

func RehydrateBoard(
	id BoardId,
	ownerId UserId,
	title Title,
	description Description,
	createdAt, updatedAt time.Time,
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

func (b *Board) Id() BoardId              { return b.id }
func (b *Board) OwnerId() UserId          { return b.ownerId }
func (b *Board) Title() Title             { return b.title }
func (b *Board) Description() Description { return b.description }
func (b *Board) CreatedAt() time.Time     { return b.createdAt }
func (b *Board) UpdatedAt() time.Time     { return b.updatedAt }

func (b *Board) Update(title Title, description Description) {
	b.title = title
	b.description = description
	b.updatedAt = time.Now().UTC()
	b.events = append(b.events, BoardUpdated{BoardID: b.id, At: b.updatedAt})
}

func (b *Board) MarkDeleted() {
	at := time.Now().UTC()
	b.events = append(b.events, BoardDeleted{BoardID: b.id, At: at})
}

func (b *Board) PullEvents() []DomainEvent {
	if len(b.events) == 0 {
		return nil
	}
	out := make([]DomainEvent, len(b.events))
	copy(out, b.events)
	b.events = nil
	return out
}
