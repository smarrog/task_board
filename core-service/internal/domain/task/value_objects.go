package task

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type Id struct {
	value uuid.UUID
}

func NewId() Id {
	return Id{uuid.New()}
}

func IdFromUUID(id uuid.UUID) (Id, error) {
	if id == uuid.Nil {
		return Id{}, common.ErrInvalidUUID
	}
	return Id{value: id}, nil
}

func IdFromString(s string) (Id, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return Id{}, fmt.Errorf("%w: %v", common.ErrInvalidUUID, err)
	}
	return IdFromUUID(id)
}

func (id Id) UUID() uuid.UUID { return id.value }
func (id Id) String() string  { return id.value.String() }
