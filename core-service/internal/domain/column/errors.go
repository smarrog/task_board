package column

import (
	"fmt"

	"github.com/smarrog/task-board/shared/domain/shared"
)

var (
	ErrNotFound        = fmt.Errorf("%s %w", "column", shared.ErrNotFound)
	ErrInvalidId       = fmt.Errorf("%s %w", "column id", shared.ErrIsInvalid)
	ErrInvalidPosition = fmt.Errorf("%s %w", "column position", shared.ErrIsInvalid)
)
