package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/api-service/internal/config"
	"github.com/smarrog/task-board/api-service/internal/middleware"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
	"google.golang.org/grpc"
)

type Handler struct {
	log *zerolog.Logger
	cfg *config.Config

	boards  v1.BoardsServiceClient
	columns v1.ColumnsServiceClient
	tasks   v1.TasksServiceClient
}

func NewHandler(log *zerolog.Logger, cfg *config.Config, coreConn *grpc.ClientConn) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		boards:  v1.NewBoardsServiceClient(coreConn),
		columns: v1.NewColumnsServiceClient(coreConn),
		tasks:   v1.NewTasksServiceClient(coreConn),
	}
}

func (h *Handler) requesterID(c *fiber.Ctx) string {
	v := c.Locals(middleware.LocalUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (h *Handler) reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Context(), h.cfg.RequestTimeout)
}
