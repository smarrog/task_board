package board

import (
	"errors"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

var (
	ErrNotFound      = fmt.Errorf("%s %w", "board", common.ErrNotFound)
	ErrInvalidId     = fmt.Errorf("%s %w", "board", common.ErrInvalidId)
	ErrInvalidUserId = fmt.Errorf("%s %w", "board", common.ErrInvalidUserId)
	ErrIsEmpty       = fmt.Errorf("%s %w", "board", common.ErrIsEmpty)

	ErrOwnerRequired = errors.New("owner_id is required")
	ErrOwnerMismatch = errors.New("owner_id does not match real owner")

	ErrTitleEmpty         = fmt.Errorf("%s %w", "column", common.ErrTitleEmpty)
	ErrTitleTooLong       = fmt.Errorf("%s %w", "column", common.ErrTitleTooLong)
	ErrDescriptionTooLong = fmt.Errorf("%s %w", "column", common.ErrDescriptionTooLong)
)
