package kafka

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	evt "github.com/smarrog/task-board/notification-service/internal/domain/events"
	"github.com/smarrog/task-board/notification-service/internal/handler"
	"github.com/smarrog/task-board/shared/messaging"
)

type handlerFn func(ctx context.Context, msg *kafka.Message, env messaging.OutboxMessage) error

type OutboxHandler struct {
	log *zerolog.Logger
	app *handler.Handler
	dlq *DlqWriter

	handlers map[string]handlerFn
}

func NewOutboxHandler(log *zerolog.Logger, app *handler.Handler, dlq *DlqWriter) *OutboxHandler {
	h := &OutboxHandler{log: log, app: app, dlq: dlq}

	h.handlers = map[string]handlerFn{
		messaging.EvtBoardCreated: makeHandler[evt.BoardCreated](h.app.HandleBoardCreated, h.publishToDlq),
		messaging.EvtBoardUpdated: makeHandler[evt.BoardUpdated](h.app.HandleBoardUpdated, h.publishToDlq),
		messaging.EvtBoardDeleted: makeHandler[evt.BoardDeleted](h.app.HandleBoardDeleted, h.publishToDlq),

		messaging.EvtColumnCreated: makeHandler[evt.ColumnCreated](h.app.HandleColumnCreated, h.publishToDlq),
		messaging.EvtColumnMoved:   makeHandler[evt.ColumnMoved](h.app.HandleColumnMoved, h.publishToDlq),
		messaging.EvtColumnDeleted: makeHandler[evt.ColumnDeleted](h.app.HandleColumnDeleted, h.publishToDlq),

		messaging.EvtTaskCreated: makeHandler[evt.TaskCreated](h.app.HandleTaskCreated, h.publishToDlq),
		messaging.EvtTaskUpdated: makeHandler[evt.TaskUpdated](h.app.HandleTaskUpdated, h.publishToDlq),
		messaging.EvtTaskMoved:   makeHandler[evt.TaskMoved](h.app.HandleTaskMoved, h.publishToDlq),
		messaging.EvtTaskDeleted: makeHandler[evt.TaskDeleted](h.app.HandleTaskDeleted, h.publishToDlq),
	}

	return h
}

func (h *OutboxHandler) HandleKafkaMessage(ctx context.Context, msg *kafka.Message) error {
	var envelope messaging.OutboxMessage
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

func (h *OutboxHandler) publishToDlq(ctx context.Context, msg *kafka.Message, err error) {
	if h.dlq == nil {
		h.log.Error().Err(err).Msg("message handling failed (DLQ disabled)")
		return
	}
	if pErr := h.dlq.Publish(ctx, msg, err); pErr != nil {
		h.log.Error().Err(pErr).Msg("failed to publish to DLQ")
	}
}

func makeHandler[T any](
	handle func(context.Context, messaging.OutboxMessage, T) error,
	publishToDlq func(ctx context.Context, msg *kafka.Message, err error),
) handlerFn {
	return func(ctx context.Context, msg *kafka.Message, env messaging.OutboxMessage) error {
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
