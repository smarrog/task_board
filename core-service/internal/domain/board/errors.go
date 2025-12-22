package board

import "errors"

// Domain (business) errors. Keep them sentinel so upper layers can map them.
var (
	ErrInvalidUUID   = errors.New("invalid uuid")
	ErrOwnerRequired = errors.New("owner_id is required")
	ErrOwnerMismatch = errors.New("owner_id does not match board owner")

	ErrTitleEmpty         = errors.New("title is empty")
	ErrTitleTooLong       = errors.New("title is too long")
	ErrDescriptionTooLong = errors.New("description is too long")
	
	ErrBoardNotFound = errors.New("board not found")
)
