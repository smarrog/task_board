package shared

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrIsEmpty    = errors.New("empty")
	ErrIsInvalid  = errors.New("invalid")
	ErrIsTooLong  = fmt.Errorf("%s %w", "too long", ErrIsInvalid)
	ErrIsRequired = errors.New("required")
	ErrIsMismatch = errors.New("mismatch")
)
