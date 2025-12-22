package board

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxTitleLength       = 255
	MaxDescriptionLength = 1_024
)

type BoardId struct {
	value uuid.UUID
}

func NewBoardID() BoardId {
	return BoardId{
		value: uuid.New(),
	}
}

func BoardIdFromUUID(id uuid.UUID) (BoardId, error) {
	if id == uuid.Nil {
		return BoardId{}, ErrInvalidUUID
	}

	return BoardId{value: id}, nil
}

func BoardIdFromString(s string) (BoardId, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return BoardId{}, fmt.Errorf("%w: %v", ErrInvalidUUID, err)
	}

	return BoardIdFromUUID(id)
}

func (id BoardId) UUID() uuid.UUID { return id.value }
func (id BoardId) String() string  { return id.value.String() }

type ColumnId struct {
	value uuid.UUID
}

func NewColumnID() ColumnId {
	return ColumnId{uuid.New()}
}

func ColumnIdFromUUID(id uuid.UUID) (ColumnId, error) {
	if id == uuid.Nil {
		return ColumnId{}, ErrInvalidUUID
	}
	return ColumnId{value: id}, nil
}

func ColumnIdFromString(s string) (ColumnId, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return ColumnId{}, fmt.Errorf("%w: %v", ErrInvalidUUID, err)
	}
	return ColumnIdFromUUID(id)
}

func (id ColumnId) UUID() uuid.UUID { return id.value }
func (id ColumnId) String() string  { return id.value.String() }

type TaskId struct {
	value uuid.UUID
}

func NewTaskID() TaskId {
	return TaskId{uuid.New()}
}

func TaskIdFromUUID(id uuid.UUID) (TaskId, error) {
	if id == uuid.Nil {
		return TaskId{}, ErrInvalidUUID
	}
	return TaskId{value: id}, nil
}

func TaskIdFromString(s string) (TaskId, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return TaskId{}, fmt.Errorf("%w: %v", ErrInvalidUUID, err)
	}
	return TaskIdFromUUID(id)
}

func (id TaskId) UUID() uuid.UUID { return id.value }
func (id TaskId) String() string  { return id.value.String() }

type UserId struct {
	value uuid.UUID
}

func UserIdFromUUID(id uuid.UUID) (UserId, error) {
	if id == uuid.Nil {
		return UserId{}, ErrInvalidUUID
	}
	return UserId{value: id}, nil
}

func UserIdFromString(s string) (UserId, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return UserId{}, fmt.Errorf("%w: %v", ErrInvalidUUID, err)
	}

	return UserIdFromUUID(id)
}

func (id UserId) UUID() uuid.UUID { return id.value }
func (id UserId) String() string  { return id.value.String() }

type Title struct {
	value string
}

func NewTitle(raw string) (Title, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return Title{}, ErrTitleEmpty
	}
	if len(v) > MaxTitleLength {
		return Title{}, ErrTitleTooLong
	}
	return Title{value: v}, nil
}

func (t Title) String() string {
	return t.value
}

type Description struct {
	value string
}

func NewDescription(raw string) (Description, error) {
	v := strings.TrimSpace(raw)
	if len(v) > MaxDescriptionLength {
		return Description{}, ErrDescriptionTooLong
	}
	return Description{value: v}, nil
}

func (d Description) String() string { return d.value }
