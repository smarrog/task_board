package column

import (
	"errors"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

var (
	ErrNotFound      = fmt.Errorf("%s %w", "column", common.ErrNotFound)
	ErrInvalidId     = fmt.Errorf("%s %w", "column", common.ErrInvalidId)
	ErrInvalidUserId = fmt.Errorf("%s %w", "column", common.ErrInvalidUserId)
	ErrIsEmpty       = fmt.Errorf("%s %w", "column", common.ErrIsEmpty)

	ErrInvalidPosition = errors.New("invalid column position")
)
