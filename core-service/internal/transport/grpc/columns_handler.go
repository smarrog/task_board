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

	createColumn *columnuc.CreateColumnUseCase
	getColumn    *columnuc.GetColumnUseCase
	moveColumn   *columnuc.MoveColumnUseCase
	deleteColumn *columnuc.DeleteColumnUseCase
}

func NewColumnsHandler(
	log *zerolog.Logger,
	createColumn *columnuc.CreateColumnUseCase,
	getColumn *columnuc.GetColumnUseCase,
	moveColumn *columnuc.MoveColumnUseCase,
	deleteColumn *columnuc.DeleteColumnUseCase,
) *ColumnsHandler {
	return &ColumnsHandler{
		log:          log,
		createColumn: createColumn,
		getColumn:    getColumn,
		moveColumn:   moveColumn,
		deleteColumn: deleteColumn,
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

func (h *ColumnsHandler) MoveColumn(ctx context.Context, req *v1.MoveColumnRequest) (*v1.MoveColumnResponse, error) {
	input := columnuc.MoveColumnRequest{
		ColumnId: req.ColumnId,
		BoardId:  req.BoardId,
		Position: int(req.Position),
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
