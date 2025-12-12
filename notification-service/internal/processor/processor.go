package processor

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	v1 "github.com/smarrog/notification-service/proto/v1"
	"google.golang.org/protobuf/proto"
)

type Processor interface {
	Handle(ctx context.Context, msg *kafka.Message) error
}

type processor struct {
	logger *zerolog.Logger
}

func NewProcessor(logger *zerolog.Logger) Processor {
	return &processor{
		logger: logger,
	}
}

func (p *processor) Handle(ctx context.Context, msg *kafka.Message) error {
	var evt v1.Event

	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}

	switch payload := evt.Payload.(type) {
	case *v1.Event_TaskCreated:
		return p.handleTaskCreatedEvent(ctx, &evt, payload)
	case *v1.Event_TaskUpdated:
		return p.handleTaskUpdatedEvent(ctx, &evt, payload)
	case *v1.Event_TaskMoved:
		return p.handleTaskMovedEvent(ctx, &evt, payload)
	default:
		return errors.New("invalid payload")
	}
}

func (p *processor) handleTaskCreatedEvent(ctx context.Context, event *v1.Event, payload *v1.Event_TaskCreated) error {
	s := payload.TaskCreated.GetSnapshot()

	p.logger.Info().
		Str("event_id", event.GetEventId()).
		Int64("task_id", s.GetTaskId()).
		Str("title", s.GetTitle()).
		Int64("board_id", s.GetBoardId()).
		Int64("column_id", s.GetColumnId()).
		Int64("assignee_id", s.GetAssigneeId()).
		Msgf("Task created: '%s' assigned to user %d", s.GetTitle(), s.GetAssigneeId())

	return nil
}

func (p *processor) handleTaskUpdatedEvent(ctx context.Context, event *v1.Event, payload *v1.Event_TaskUpdated) error {
	s := payload.TaskUpdated.GetSnapshot()

	p.logger.Info().
		Str("event_id", event.GetEventId()).
		Int64("task_id", s.GetTaskId()).
		Str("title", s.GetTitle()).
		Str("description", s.GetDescription()).
		Int64("assignee_id", s.GetAssigneeId()).
		Msgf("Task updated: '%s' (task_id=%d)", s.GetTitle(), s.GetTaskId())

	return nil
}

func (p *processor) handleTaskMovedEvent(ctx context.Context, event *v1.Event, payload *v1.Event_TaskMoved) error {
	s := payload.TaskMoved.GetSnapshot()

	p.logger.Info().
		Str("event_id", event.GetEventId()).
		Int64("task_id", s.GetTaskId()).
		Str("title", s.GetTitle()).
		Int64("from_column_id", s.GetColumnId()).
		Int64("to_column_id", payload.TaskMoved.GetToColumnId()).
		Int64("from_board_id", s.GetBoardId()).
		Int64("to_board_id", payload.TaskMoved.GetToBoardId()).
		Msgf("Task moved: '%s' moved from column %d to %d", s.GetTitle(), s.GetColumnId(), payload.TaskMoved.GetToColumnId())

	return nil
}
