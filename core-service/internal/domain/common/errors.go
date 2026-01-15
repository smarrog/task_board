package common

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidId     = errors.New("invalid id")
	ErrInvalidUserId = errors.New("invalid user id")
	ErrIsEmpty       = errors.New("is empty")

	ErrTitleEmpty         = errors.New("title is empty")
	ErrTitleTooLong       = errors.New("title is too long")
	ErrDescriptionTooLong = errors.New("description is too long")
)
