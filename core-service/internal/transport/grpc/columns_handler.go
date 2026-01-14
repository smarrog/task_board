package grpc

import (
	"context"

	"github.com/rs/zerolog"
	columdo "github.com/smarrog/task-board/core-service/internal/domain/column"
	columnuc "github.com/smarrog/task-board/core-service/internal/usecase/column"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
)

type ColumnsHandler struct {
	v1.UnimplementedColumnsServiceServer

	log *zerolog.Logger

	createColumn  *columnuc.CreateColumnUseCase
	getColumn     *columnuc.GetColumnUseCase
	getColumnFull *columnuc.GetColumnFullUseCase
	moveColumn    *columnuc.MoveColumnUseCase
	deleteColumn  *columnuc.DeleteColumnUseCase
}

func NewColumnsHandler(
	log *zerolog.Logger,
	createColumn *columnuc.CreateColumnUseCase,
	getColumn *columnuc.GetColumnUseCase,
	getColumnFull *columnuc.GetColumnFullUseCase,
	moveColumn *columnuc.MoveColumnUseCase,
	deleteColumn *columnuc.DeleteColumnUseCase,
) *ColumnsHandler {
	return &ColumnsHandler{
		log:           log,
		createColumn:  createColumn,
		getColumn:     getColumn,
		getColumnFull: getColumnFull,
		moveColumn:    moveColumn,
		deleteColumn:  deleteColumn,
	}
}

func (h *ColumnsHandler) CreateColumn(ctx context.Context, req *v1.CreateColumnRequest) (*v1.CreateColumnResponse, error) {
	input := columnuc.CreateColumnInput{
		BoardId:  req.BoardId,
		Position: int(req.Position),
	}

	output, err := h.createColumn.Execute(ctx, input)
	if err != nil {
		return nil, mapColumnsErr(err)
	}

	return &v1.CreateColumnResponse{
		Column: toProtoColumn(output.Column),
	}, nil
}

func (h *ColumnsHandler) GetColumn(ctx context.Context, req *v1.GetColumnRequest) (*v1.GetColumnResponse, error) {
	input := columnuc.GetColumnInput{
		ColumnId: req.ColumnId,
	}

	output, err := h.getColumn.Execute(ctx, input)
	if err != nil {
		return nil, mapColumnsErr(err)
	}

	return &v1.GetColumnResponse{
		Column: toProtoColumn(output.Column),
	}, nil
}

func (h *ColumnsHandler) GetColumnFull(ctx context.Context, req *v1.GetColumnFullRequest) (*v1.GetColumnFullResponse, error) {
	input := columnuc.GetColumnFullInput{ColumnId: req.ColumnId}

	output, err := h.getColumnFull.Execute(ctx, input)
	if err != nil {
		return nil, mapColumnsErr(err)
	}

	tasks := make([]*v1.Task, 0, len(output.Tasks))
	for _, t := range output.Tasks {
		tasks = append(tasks, toProtoTask(t))
	}

	return &v1.GetColumnFullResponse{
		Data: &v1.ColumnFull{
			Column: toProtoColumn(output.Column),
			Tasks:  tasks,
		},
	}, nil
}

func (h *ColumnsHandler) MoveColumn(ctx context.Context, req *v1.MoveColumnRequest) (*v1.MoveColumnResponse, error) {
	input := columnuc.MoveColumnInput{
		ColumnId:   req.ColumnId,
		ToPosition: int(req.ToPosition),
	}

	output, err := h.moveColumn.Execute(ctx, input)
	if err != nil {
		return nil, mapColumnsErr(err)
	}

	return &v1.MoveColumnResponse{
		Column: toProtoColumn(output.Column),
	}, nil
}

func (h *ColumnsHandler) DeleteColumn(ctx context.Context, req *v1.DeleteColumnRequest) (*v1.DeleteColumnResponse, error) {
	input := columnuc.DeleteColumnInput{
		ColumnId: req.ColumnId,
	}

	_, err := h.deleteColumn.Execute(ctx, input)
	if err != nil {
		return nil, mapColumnsErr(err)
	}

	return &v1.DeleteColumnResponse{}, nil
}

func toProtoColumn(c *columdo.Column) *v1.Column {
	return &v1.Column{
		Id:       c.Id().String(),
		BoardId:  c.BoardId().String(),
		Position: int32(c.Position().Int()),
	}
}

func mapColumnsErr(err error) error {
	switch {
	default:
		return mapCommonErr(err)
	}
}
