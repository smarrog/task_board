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

type Id struct {
	value uuid.UUID
}

func NewId() Id {
	return Id{
		value: uuid.New(),
	}
}

func IdFromUUID(id uuid.UUID) (Id, error) {
	if id == uuid.Nil {
		return Id{}, ErrInvalidId
	}

	return Id{value: id}, nil
}

func IdFromString(s string) (Id, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return Id{}, fmt.Errorf("%w: %v", ErrInvalidId, err)
	}

	return IdFromUUID(id)
}

func (id Id) UUID() uuid.UUID { return id.value }
func (id Id) String() string  { return id.value.String() }

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
