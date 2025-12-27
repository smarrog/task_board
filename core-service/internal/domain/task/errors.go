package task

import (
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

var (
	ErrNotFound = fmt.Errorf("%s %w", "task", common.ErrNotFound)
)
