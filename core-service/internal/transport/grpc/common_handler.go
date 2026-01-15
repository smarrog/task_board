package grpc

import (
	"errors"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapCommonErr(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, common.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, common.ErrInvalidUserId),
		errors.Is(err, common.ErrTitleEmpty),
		errors.Is(err, common.ErrTitleTooLong),
		errors.Is(err, common.ErrDescriptionTooLong):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, "internal error")
	}
}
