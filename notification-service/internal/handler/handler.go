package handler

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/notification-service/internal/domain/events"
	"github.com/smarrog/task-board/notification-service/internal/infrastructure/persistence"
	"github.com/smarrog/task-board/shared/messaging"
)

type Handler struct {
	notifier Notifier
	repo     *persistence.NotificationsRepo
}

func NewHandler(notifier Notifier, repo *persistence.NotificationsRepo) *Handler {
	return &Handler{notifier: notifier, repo: repo}
}

func (h *Handler) HandleBoardCreated(ctx context.Context, env messaging.OutboxMessage, e events.BoardCreated) error {
	text := fmt.Sprintf("Board created: '%s' (board_id=%s, owner_id=%s)", e.Title, e.Id, e.OwnerId)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleBoardUpdated(ctx context.Context, env messaging.OutboxMessage, e events.BoardUpdated) error {
	text := fmt.Sprintf("Board updated: '%s' (board_id=%s)", e.Title, e.Id)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleBoardDeleted(ctx context.Context, env messaging.OutboxMessage, e events.BoardDeleted) error {
	text := fmt.Sprintf("Board deleted: (board_id=%s)", e.Id)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleColumnCreated(ctx context.Context, env messaging.OutboxMessage, e events.ColumnCreated) error {
	text := fmt.Sprintf("Column created: (column_id=%s, board_id=%s)", e.Id, e.BoardId)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleColumnMoved(ctx context.Context, env messaging.OutboxMessage, e events.ColumnMoved) error {
	text := fmt.Sprintf("Column moved: (column_id=%s, from_position=%d, to_position=%d)", e.Id, e.FromPosition, e.ToPosition)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleColumnDeleted(ctx context.Context, env messaging.OutboxMessage, e events.ColumnDeleted) error {
	text := fmt.Sprintf("Column deleted: (column_id=%s)", e.Id)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleTaskCreated(ctx context.Context, env messaging.OutboxMessage, e events.TaskCreated) error {
	text := fmt.Sprintf("Task created: '%s' (task_id=%s, column_id=%s, assignee_id=%s)", e.Title, e.Id, e.ColumnId, e.AssigneeId)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleTaskUpdated(ctx context.Context, env messaging.OutboxMessage, e events.TaskUpdated) error {
	text := fmt.Sprintf("Task updated: '%s' (task_id=%s, assignee_id=%s)", e.Title, e.Id, e.AssigneeId)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleTaskMoved(ctx context.Context, env messaging.OutboxMessage, e events.TaskMoved) error {
	text := fmt.Sprintf("Task moved: (task_id=%s, from_column_id=%s, to_column_id=%s, from_position=%d, to_position=%d)", e.Id, e.FromColumnId, e.ToColumnId, e.FromPosition, e.ToPosition)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) HandleTaskDeleted(ctx context.Context, env messaging.OutboxMessage, e events.TaskDeleted) error {
	text := fmt.Sprintf("Task deleted: (task_id=%s)", e.Id)
	if err := h.saveHistory(ctx, env, text); err != nil {
		return err
	}
	return h.notifier.Notify(ctx, Notification{Text: text})
}

func (h *Handler) saveHistory(ctx context.Context, env messaging.OutboxMessage, text string) error {
	if h.repo == nil {
		return nil
	}
	return h.repo.Save(ctx, persistence.NotificationRecord{
		OutboxId:       env.Id,
		EventType:      env.EventType,
		AggregateType:  env.AggregateType,
		AggregateId:    env.AggregateId,
		EventCreatedAt: env.CreatedAt,
		Version:        env.Version,
		Payload:        env.Payload,
		Text:           text,
	})
}
