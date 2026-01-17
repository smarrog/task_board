package shared

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	Name() string
	OccurredAt() time.Time
}

type UserId struct {
	value uuid.UUID
}

func UserIdFromUUID(id uuid.UUID) (UserId, error) {
	if id == uuid.Nil {
		return UserId{}, ErrIsInvalid
	}
	return UserId{value: id}, nil
}

func UserIdFromString(s string) (UserId, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return UserId{}, fmt.Errorf("%w: %v", ErrIsInvalid, err)
	}

	return UserIdFromUUID(id)
}

func (id UserId) UUID() uuid.UUID { return id.value }
func (id UserId) String() string  { return id.value.String() }
