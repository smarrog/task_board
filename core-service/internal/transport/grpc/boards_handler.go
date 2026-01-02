package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	boarddo "github.com/smarrog/task-board/core-service/internal/domain/board"
	boarduc "github.com/smarrog/task-board/core-service/internal/usecase/board"
	"github.com/smarrog/task-board/shared/proto/base/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BoardsHandler struct {
	v1.UnimplementedBoardsServiceServer

	log *zerolog.Logger

	createBoard *boarduc.CreateBoardUseCase
	getBoard    *boarduc.GetBoardUseCase
	listBoards  *boarduc.ListBoardsUseCase
	updateBoard *boarduc.UpdateBoardUseCase
	deleteBoard *boarduc.DeleteBoardUseCase
}

func NewBoardsHandler(
	log *zerolog.Logger,
	createBoard *boarduc.CreateBoardUseCase,
	getBoard *boarduc.GetBoardUseCase,
	listBoards *boarduc.ListBoardsUseCase,
	updateBoard *boarduc.UpdateBoardUseCase,
	deleteBoard *boarduc.DeleteBoardUseCase,
) *BoardsHandler {
	return &BoardsHandler{
		log:         log,
		createBoard: createBoard,
		getBoard:    getBoard,
		listBoards:  listBoards,
		updateBoard: updateBoard,
		deleteBoard: deleteBoard,
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
		Board: toProtoBoard(output.Board),
	}, nil
}

func (h *BoardsHandler) GetBoard(ctx context.Context, req *v1.GetBoardRequest) (*v1.GetBoardResponse, error) {
	input := boarduc.GetBoardInput{
		BoardId: req.GetBoardId(),
	}
	output, err := h.getBoard.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.GetBoardResponse{
		Board: toProtoBoard(output.Board),
	}, nil
}

func (h *BoardsHandler) ListBoards(ctx context.Context, req *v1.ListBoardsRequest) (*v1.ListBoardsResponse, error) {
	input := boarduc.ListBoardsInput{
		OwnerId: req.GetOwnerId(),
	}
	output, err := h.listBoards.Execute(ctx, input)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	out := make([]*v1.Board, len(output.Boards))
	for i, b := range output.Boards {
		out[i] = toProtoBoard(b)
	}

	return &v1.ListBoardsResponse{
		Boards: out,
	}, nil
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

	return &v1.UpdateBoardResponse{
		Board: toProtoBoard(output.Board),
	}, nil
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
