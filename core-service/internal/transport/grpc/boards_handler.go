package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domain "github.com/smarrog/task-board/core-service/internal/domain/board"
	boarduc "github.com/smarrog/task-board/core-service/internal/usecase/board"
	"github.com/smarrog/task-board/shared/proto/v1"
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
	b, err := h.createBoard.Execute(ctx, req.GetOwnerId(), req.GetTitle(), req.GetDescription())
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.CreateBoardResponse{
		Board: toProtoBoard(b),
	}, nil
}

func (h *BoardsHandler) GetBoard(ctx context.Context, req *v1.GetBoardRequest) (*v1.GetBoardResponse, error) {
	b, err := h.getBoard.Execute(ctx, req.GetBoardId())
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.GetBoardResponse{
		Board: toProtoBoard(b),
	}, nil
}

func (h *BoardsHandler) ListBoards(ctx context.Context, req *v1.ListBoardsRequest) (*v1.ListBoardsResponse, error) {
	boards, err := h.listBoards.Execute(ctx, req.GetOwnerId())
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	out := make([]*v1.Board, 0, len(boards))
	for _, b := range boards {
		out = append(out, toProtoBoard(b))
	}

	return &v1.ListBoardsResponse{
		Boards: out,
	}, nil
}

func (h *BoardsHandler) UpdateBoard(ctx context.Context, req *v1.UpdateBoardRequest) (*v1.UpdateBoardResponse, error) {
	b, err := h.updateBoard.Execute(ctx, req.GetOwnerId(), req.GetBoardId(), req.GetTitle(), req.GetDescription())
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.UpdateBoardResponse{
		Board: toProtoBoard(b),
	}, nil
}

func (h *BoardsHandler) DeleteBoard(ctx context.Context, req *v1.DeleteBoardRequest) (*v1.DeleteBoardResponse, error) {
	if err := h.deleteBoard.Execute(ctx, req.GetBoardId()); err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.DeleteBoardResponse{}, nil
}

func toProtoBoard(b *domain.Board) *v1.Board {
	return &v1.Board{
		Id:          b.Id().String(),
		OwnerId:     b.OwnerId().String(),
		Title:       b.Title().String(),
		Description: b.Description().String(),
	}
}

func mapBoardsErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, domain.ErrInvalidUUID) ||
		errors.Is(err, domain.ErrTitleEmpty) ||
		errors.Is(err, domain.ErrTitleTooLong) ||
		errors.Is(err, domain.ErrDescriptionTooLong) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, domain.ErrBoardNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, domain.ErrOwnerMismatch) {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	return status.Error(codes.Internal, "internal error")
}
