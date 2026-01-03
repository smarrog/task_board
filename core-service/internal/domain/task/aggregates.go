package task

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type Task struct {
	id          Id
	columnId    column.Id
	position    Position
	title       Title
	description Description
	assigneeId  common.UserId
	createdAt   time.Time
	updatedAt   time.Time
	events      []common.DomainEvent
}

func New(
	columnId column.Id,
	position Position,
	title Title,
	desc Description,
	assigneeId common.UserId,
) *Task {
	now := time.Now().UTC()

	t := &Task{
		id:          NewId(),
		columnId:    columnId,
		position:    position,
		title:       title,
		description: desc,
		assigneeId:  assigneeId,
		createdAt:   now,
		updatedAt:   now,
	}

	t.events = append(t.events, CreatedEvent{
		Id:          t.id,
		ColumnId:    columnId,
		Position:    position,
		Title:       title,
		Description: desc,
		AssigneeId:  assigneeId,
		At:          now,
	})
	return t
}

func Rehydrate(
	id Id,
	columnId column.Id,
	position Position,
	title Title,
	desc Description,
	assigneeId common.UserId,
	createdAt time.Time,
	updatedAt time.Time,
) *Task {
	return &Task{
		id:          id,
		columnId:    columnId,
		position:    position,
		title:       title,
		description: desc,
		assigneeId:  assigneeId,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (t *Task) Id() Id                    { return t.id }
func (t *Task) ColumnId() column.Id       { return t.columnId }
func (t *Task) Position() Position        { return t.position }
func (t *Task) Title() Title              { return t.title }
func (t *Task) Description() Description  { return t.description }
func (t *Task) AssigneeId() common.UserId { return t.assigneeId }
func (t *Task) CreatedAt() time.Time      { return t.createdAt }
func (t *Task) UpdatedAt() time.Time      { return t.updatedAt }

func (t *Task) Update(title Title, desc Description, assigneeId common.UserId) {
	now := time.Now().UTC()

	t.title = title
	t.description = desc
	t.assigneeId = assigneeId
	t.updatedAt = now

	t.events = append(t.events, UpdatedEvent{
		Id:          t.id,
		Title:       t.title,
		Description: t.description,
		AssigneeId:  assigneeId,
		At:          t.updatedAt,
	})
}

func (t *Task) Move(toColumnId column.Id, toPosition Position) {
	now := time.Now().UTC()

	fromColumnId := t.columnId
	fromPosition := t.position

	t.columnId = toColumnId
	t.position = toPosition
	t.updatedAt = now

	t.events = append(t.events, MoveEvent{
		Id:           t.id,
		FromColumnId: fromColumnId,
		ToColumnId:   toColumnId,
		FromPosition: fromPosition,
		ToPosition:   toPosition,
		At:           t.updatedAt,
	})
}

func (t *Task) PullEvents() []common.DomainEvent {
	if len(t.events) == 0 {
		return nil
	}
	out := make([]common.DomainEvent, len(t.events))
	copy(out, t.events)
	t.events = nil
	return out
}
