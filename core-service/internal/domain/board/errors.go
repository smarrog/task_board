package board

import (
	"fmt"

	"github.com/smarrog/task-board/shared/domain/shared"
)

var (
	ErrNotFound           = fmt.Errorf("%s %w", "board", shared.ErrNotFound)
	ErrInvalidId          = fmt.Errorf("%s %w", "board id", shared.ErrIsInvalid)
	ErrOwnerRequired      = fmt.Errorf("%s %w", "board owner id", shared.ErrIsRequired)
	ErrOwnerMismatch      = fmt.Errorf("%s %w", "board owner id", shared.ErrIsMismatch)
	ErrTitleEmpty         = fmt.Errorf("%s %w", "board title", shared.ErrIsEmpty)
	ErrTitleTooLong       = fmt.Errorf("%s %w", "board title", shared.ErrIsTooLong)
	ErrDescriptionTooLong = fmt.Errorf("%s %w", "board description", shared.ErrIsTooLong)
)
