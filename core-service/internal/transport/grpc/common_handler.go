package grpc

import (
	"errors"

	"github.com/smarrog/task-board/shared/domain/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapCommonErr(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, shared.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, shared.ErrIsEmpty),
		errors.Is(err, shared.ErrIsInvalid),
		errors.Is(err, shared.ErrIsRequired),
		errors.Is(err, shared.ErrIsMismatch):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, "internal error")
	}
}
