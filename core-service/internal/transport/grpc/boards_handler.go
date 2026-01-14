package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	boarddo "github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	boarduc "github.com/smarrog/task-board/core-service/internal/usecase/board"
	"github.com/smarrog/task-board/shared/proto/base/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BoardsHandler struct {
	v1.UnimplementedBoardsServiceServer

	log *zerolog.Logger

	createBoard    *boarduc.CreateBoardUseCase
	getBoardFull   *boarduc.GetBoardFullUseCase
	listBoards     *boarduc.ListBoardsUseCase
	listBoardsFull *boarduc.ListBoardsFullUseCase
	updateBoard    *boarduc.UpdateBoardUseCase
	deleteBoard    *boarduc.DeleteBoardUseCase
}

func NewBoardsHandler(
	log *zerolog.Logger,
	createBoard *boarduc.CreateBoardUseCase,
	getBoardFull *boarduc.GetBoardFullUseCase,
	listBoards *boarduc.ListBoardsUseCase,
	listBoardsFull *boarduc.ListBoardsFullUseCase,
	updateBoard *boarduc.UpdateBoardUseCase,
	deleteBoard *boarduc.DeleteBoardUseCase,
) *BoardsHandler {
	return &BoardsHandler{
		log:            log,
		createBoard:    createBoard,
		getBoardFull:   getBoardFull,
		listBoards:     listBoards,
		listBoardsFull: listBoardsFull,
		updateBoard:    updateBoard,
		deleteBoard:    deleteBoard,
	}
}

func (h *BoardsHandler) CreateBoard(ctx context.Context, req *v1.CreateBoardRequest) (*v1.CreateBoardResponse, error) {
	input := boarduc.CreateBoardInput{
		OwnerId:     req.OwnerId,
		Title:       req.Title,
		Description: req.Description,
	}
	output, err := h.createBoard.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.CreateBoardResponse{
		Data: &v1.BoardFull{
			Board:   toProtoBoard(output.Board),
			Columns: []*v1.ColumnFull{},
		},
	}, nil
}

func (h *BoardsHandler) GetBoard(ctx context.Context, req *v1.GetBoardRequest) (*v1.GetBoardResponse, error) {
	input := boarduc.GetBoardFullInput{BoardId: req.GetBoardId()}
	output, err := h.getBoardFull.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.GetBoardResponse{
		Data: toProtoBoardFull(output.Board, output.Columns, output.Tasks),
	}, nil
}

func (h *BoardsHandler) ListBoards(ctx context.Context, req *v1.ListBoardsRequest) (*v1.ListBoardsResponse, error) {
	input := boarduc.ListBoardsFullInput{
		OwnerId: req.GetOwnerId(),
	}
	output, err := h.listBoardsFull.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	full := make([]*v1.BoardFull, 0, len(output.Items))
	for _, it := range output.Items {
		full = append(full, toProtoBoardFull(it.Board, it.Columns, it.Tasks))
	}

	return &v1.ListBoardsResponse{Boards: full}, nil
}

func (h *BoardsHandler) UpdateBoard(ctx context.Context, req *v1.UpdateBoardRequest) (*v1.UpdateBoardResponse, error) {
	input := boarduc.UpdateBoardInput{
		OwnerId:     req.GetOwnerId(),
		BoardId:     req.GetBoardId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}
	output, err := h.updateBoard.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	// Return updated full board.
	fo, e := h.getBoardFull.Execute(ctx, boarduc.GetBoardFullInput{BoardId: output.Board.Id().String()})
	if e != nil {
		return nil, mapBoardsErr(e)
	}
	return &v1.UpdateBoardResponse{Data: toProtoBoardFull(fo.Board, fo.Columns, fo.Tasks)}, nil
}

func (h *BoardsHandler) DeleteBoard(ctx context.Context, req *v1.DeleteBoardRequest) (*v1.DeleteBoardResponse, error) {
	input := boarduc.DeleteBoardInput{
		BoardId: req.GetBoardId(),
	}
	_, err := h.deleteBoard.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.DeleteBoardResponse{}, nil
}

func toProtoBoardFull(b *boarddo.Board, cols []*column.Column, ts []*task.Task) *v1.BoardFull {
	tasksByColumn := make(map[string][]*v1.Task)
	for _, t := range ts {
		colID := t.ColumnId().String()
		tasksByColumn[colID] = append(tasksByColumn[colID], &v1.Task{
			Id:          t.Id().String(),
			ColumnId:    t.ColumnId().String(),
			Position:    int32(t.Position()),
			Title:       t.Title().String(),
			Description: t.Description().String(),
			AssigneeId:  t.AssigneeId().String(),
		})
	}

	colWithTasks := make([]*v1.ColumnFull, 0, len(cols))
	for _, c := range cols {
		pc := &v1.Column{Id: c.Id().String(), BoardId: c.BoardId().String(), Position: int32(c.Position())}
		colWithTasks = append(colWithTasks, &v1.ColumnFull{Column: pc, Tasks: tasksByColumn[c.Id().String()]})
	}

	return &v1.BoardFull{Board: toProtoBoard(b), Columns: colWithTasks}
}

func toProtoBoard(b *boarddo.Board) *v1.Board {
	return &v1.Board{
		Id:          b.Id().String(),
		OwnerId:     b.OwnerId().String(),
		Title:       b.Title().String(),
		Description: b.Description().String(),
	}
}

func mapBoardsErr(err error) error {
	switch {
	case errors.Is(err, boarddo.ErrOwnerMismatch):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return mapCommonErr(err)
	}
}
