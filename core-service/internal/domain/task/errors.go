package task

import (
	"errors"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

var (
	ErrNotFound  = fmt.Errorf("%s %w", "task", common.ErrNotFound)
	ErrInvalidId = fmt.Errorf("%s %w", "task", common.ErrInvalidId)

	ErrInvalidPosition = errors.New("invalid task position")

	ErrTitleEmpty         = fmt.Errorf("%s %w", "task", common.ErrTitleEmpty)
	ErrTitleTooLong       = fmt.Errorf("%s %w", "task", common.ErrTitleTooLong)
	ErrDescriptionTooLong = fmt.Errorf("%s %w", "task", common.ErrDescriptionTooLong)
)
