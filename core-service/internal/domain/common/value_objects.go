package common

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxTitleLength       = 255
	MaxDescriptionLength = 1_024
)

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
