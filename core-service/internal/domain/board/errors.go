package board

import (
	"errors"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

var (
	ErrNotFound      = fmt.Errorf("%s %w", "board", common.ErrNotFound)
	ErrOwnerRequired = errors.New("owner_id is required")
	ErrOwnerMismatch = errors.New("owner_id does not match real owner")
)
