package grpc

import (
	"errors"

	domain "github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapCommonErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, common.ErrInvalidUUID) ||
		errors.Is(err, common.ErrTitleEmpty) ||
		errors.Is(err, common.ErrTitleTooLong) ||
		errors.Is(err, common.ErrDescriptionTooLong) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, domain.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, domain.ErrOwnerMismatch) {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	return status.Error(codes.Internal, "internal error")
}
