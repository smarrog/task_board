package task

import (
	"fmt"

	"github.com/smarrog/task-board/shared/domain/shared"
)

var (
	ErrNotFound           = fmt.Errorf("%s %w", "task", shared.ErrNotFound)
	ErrInvalidId          = fmt.Errorf("%s %w", "task id", shared.ErrIsInvalid)
	ErrTitleEmpty         = fmt.Errorf("%s %w", "task title", shared.ErrIsEmpty)
	ErrInvalidPosition    = fmt.Errorf("%s %w", "task position", shared.ErrIsInvalid)
	ErrTitleTooLong       = fmt.Errorf("%s %w", "task title", shared.ErrIsTooLong)
	ErrDescriptionTooLong = fmt.Errorf("%s %w", "task description", shared.ErrIsTooLong)
)
