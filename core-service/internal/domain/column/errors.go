package column

import (
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

var (
	ErrNotFound = fmt.Errorf("%s %w", "column", common.ErrNotFound)
)
