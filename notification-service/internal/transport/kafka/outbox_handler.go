package kafka

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	kafkago "github.com/segmentio/kafka-go"
	infra "github.com/smarrog/task-board/notification-service/internal/infrastructure/kafka"
	uc "github.com/smarrog/task-board/notification-service/internal/usecase/notification"
	"github.com/smarrog/task-board/shared/domain/board"
	"github.com/smarrog/task-board/shared/domain/column"
	"github.com/smarrog/task-board/shared/domain/outbox"
	"github.com/smarrog/task-board/shared/domain/task"
)

type handlerFn func(ctx context.Context, msg *kafkago.Message, env outbox.Message) error

type OutboxHandler struct {
	log *zerolog.Logger
	uc  *uc.Handler
	dlq *infra.DlqWriter

	handlers map[string]handlerFn
}

func NewOutboxHandler(log *zerolog.Logger, ucHandler *uc.Handler, dlq *infra.DlqWriter) *OutboxHandler {
	h := &OutboxHandler{log: log, uc: ucHandler, dlq: dlq}

	h.handlers = map[string]handlerFn{
		board.EvtCreated: makeHandler[board.CreatedEvent](h.uc.HandleBoardCreated, h.publishToDlq),
		board.EvtUpdated: makeHandler[board.UpdatedEvent](h.uc.HandleBoardUpdated, h.publishToDlq),
		board.EvtDeleted: makeHandler[board.DeletedEvent](h.uc.HandleBoardDeleted, h.publishToDlq),

		column.EvtCreated: makeHandler[column.CreatedEvent](h.uc.HandleColumnCreated, h.publishToDlq),
		column.EvtMoved:   makeHandler[column.MovedEvent](h.uc.HandleColumnMoved, h.publishToDlq),
		column.EvtDeleted: makeHandler[column.DeletedEvent](h.uc.HandleColumnDeleted, h.publishToDlq),

		task.EvtCreated: makeHandler[task.CreatedEvent](h.uc.HandleTaskCreated, h.publishToDlq),
		task.EvtUpdated: makeHandler[task.UpdatedEvent](h.uc.HandleTaskUpdated, h.publishToDlq),
		task.EvtMoved:   makeHandler[task.MovedEvent](h.uc.HandleTaskMoved, h.publishToDlq),
		task.EvtDeleted: makeHandler[task.DeletedEvent](h.uc.HandleTaskDeleted, h.publishToDlq),
	}

	return h
}

func (h *OutboxHandler) HandleKafkaMessage(ctx context.Context, msg *kafkago.Message) error {
	var envelope outbox.Message
	if err := json.Unmarshal(msg.Value, &envelope); err != nil {
		h.publishToDlq(ctx, msg, err)
		return nil
	}

	if handler, ok := h.handlers[envelope.EventType]; ok {
		return handler(ctx, msg, envelope)
	}

	h.log.Debug().Str("event_type", envelope.EventType).Msg("skip unknown event type")
	return nil
}

func (h *OutboxHandler) publishToDlq(ctx context.Context, msg *kafkago.Message, err error) {
	if h.dlq == nil {
		h.log.Error().Err(err).Msg("message handling failed (DLQ disabled)")
		return
	}
	if pErr := h.dlq.Publish(ctx, msg, err); pErr != nil {
		h.log.Error().Err(pErr).Msg("failed to publish to DLQ")
	}
}

func makeHandler[T any](
	handle func(context.Context, outbox.Message, T) error,
	publishToDlq func(ctx context.Context, msg *kafkago.Message, err error),
) handlerFn {
	return func(ctx context.Context, msg *kafkago.Message, env outbox.Message) error {
		var e T

		if err := json.Unmarshal(env.Payload, &e); err != nil {
			publishToDlq(ctx, msg, err)
			return nil
		}

		if err := handle(ctx, env, e); err != nil {
			publishToDlq(ctx, msg, err)
			return nil
		}

		return nil
	}
}
