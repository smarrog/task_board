package http

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcToHTTP(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "internal_error")
	}

	switch st.Code() {
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return fiber.NewError(fiber.StatusBadRequest, st.Message())
	case codes.NotFound:
		return fiber.NewError(fiber.StatusNotFound, st.Message())
	case codes.AlreadyExists:
		return fiber.NewError(fiber.StatusConflict, st.Message())
	case codes.Unauthenticated:
		return fiber.NewError(fiber.StatusUnauthorized, st.Message())
	case codes.PermissionDenied:
		return fiber.NewError(fiber.StatusForbidden, st.Message())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, st.Message())
	}
}

func NotImplemented(msg string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotImplemented, msg)
	}
}
