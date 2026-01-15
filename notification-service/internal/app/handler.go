package app

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/notification-service/internal/domain/events"
)

type Handler struct {
	notifier Notifier
}

func NewHandler(notifier Notifier) *Handler {
	return &Handler{notifier: notifier}
}

func (h *Handler) HandleBoardCreated(ctx context.Context, e events.BoardCreated) error {
	text := fmt.Sprintf("Board created: '%s' (board_id=%s, owner_id=%s)", e.Title, e.Id, e.OwnerId)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleBoardUpdated(ctx context.Context, e events.BoardUpdated) error {
	text := fmt.Sprintf("Board updated: '%s' (board_id=%s)", e.Title, e.Id)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleBoardDeleted(ctx context.Context, e events.BoardDeleted) error {
	text := fmt.Sprintf("Board deleted: (board_id=%s)", e.Id)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleColumnCreated(ctx context.Context, e events.ColumnCreated) error {
	text := fmt.Sprintf("Column created: (column_id=%s, board_id=%s)", e.Id, e.BoardId)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleColumnMoved(ctx context.Context, e events.ColumnMoved) error {
	text := fmt.Sprintf("Column moved: (column_id=%s, from_position=%d, to_position=%d)", e.Id, e.FromPosition, e.ToPosition)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleColumnDeleted(ctx context.Context, e events.ColumnDeleted) error {
	text := fmt.Sprintf("Column deleted: (column_id=%s)", e.Id)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleTaskCreated(ctx context.Context, e events.TaskCreated) error {
	text := fmt.Sprintf("Task created: '%s' (task_id=%s, column_id=%s, assignee_id=%s)", e.Title, e.Id, e.ColumnId, e.AssigneeId)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleTaskUpdated(ctx context.Context, e events.TaskUpdated) error {
	text := fmt.Sprintf("Task updated: '%s' (task_id=%s, assignee_id=%s)", e.Title, e.Id, e.AssigneeId)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleTaskMoved(ctx context.Context, e events.TaskMoved) error {
	text := fmt.Sprintf("Task moved: (task_id=%s, from_column_id=%s, to_column_id=%s, from_position=%d, to_position=%d)", e.Id, e.FromColumnId, e.ToColumnId, e.FromPosition, e.ToPosition)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}

func (h *Handler) HandleTaskDeleted(ctx context.Context, e events.TaskDeleted) error {
	text := fmt.Sprintf("Task deleted: (task_id=%s)", e.Id)
	return h.notifier.Notify(ctx, Notification{
		Text: text,
	})
}
