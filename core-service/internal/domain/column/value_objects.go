package column

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Id struct {
	value uuid.UUID
}

func NewId() Id {
	return Id{uuid.New()}
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

type Position int

func NewPosition(pos int) (Position, error) {
	if pos < 0 {
		return -1, ErrInvalidPosition
	}

	return Position(pos), nil
}

func (p Position) Int() int { return int(p) }
