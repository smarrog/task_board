package task

import (
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/shared/domain/shared"
	"github.com/smarrog/task-board/shared/domain/task"
)

type Task struct {
	id          Id
	columnId    column.Id
	position    Position
	title       Title
	description Description
	assigneeId  shared.UserId
	createdAt   time.Time
	updatedAt   time.Time
	events      []shared.DomainEvent
}

func New(
	columnId column.Id,
	position Position,
	title Title,
	desc Description,
	assigneeId shared.UserId,
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

	t.events = append(t.events, task.CreatedEvent{
		Id:          t.id.String(),
		ColumnId:    columnId.String(),
		Position:    position.Int(),
		Title:       title.String(),
		Description: desc.String(),
		AssigneeId:  assigneeId.String(),
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
	assigneeId shared.UserId,
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
func (t *Task) AssigneeId() shared.UserId { return t.assigneeId }
func (t *Task) CreatedAt() time.Time      { return t.createdAt }
func (t *Task) UpdatedAt() time.Time      { return t.updatedAt }

func (t *Task) Update(title Title, desc Description, assigneeId shared.UserId) {
	now := time.Now().UTC()

	t.title = title
	t.description = desc
	t.assigneeId = assigneeId
	t.updatedAt = now

	t.events = append(t.events, task.UpdatedEvent{
		Id:          t.id.String(),
		Title:       t.title.String(),
		Description: t.description.String(),
		AssigneeId:  assigneeId.String(),
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

	t.events = append(t.events, task.MovedEvent{
		Id:           t.id.String(),
		FromColumnId: fromColumnId.String(),
		ToColumnId:   toColumnId.String(),
		FromPosition: fromPosition.Int(),
		ToPosition:   toPosition.Int(),
		At:           t.updatedAt,
	})
}

func (t *Task) PullEvents() []shared.DomainEvent {
	if len(t.events) == 0 {
		return nil
	}
	out := make([]shared.DomainEvent, len(t.events))
	copy(out, t.events)
	t.events = nil
	return out
}
