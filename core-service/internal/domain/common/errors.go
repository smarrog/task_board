package common

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidUUID        = errors.New("invalid uuid")
	ErrTitleEmpty         = errors.New("title is empty")
	ErrTitleTooLong       = errors.New("title is too long")
	ErrDescriptionTooLong = errors.New("description is too long")
)
