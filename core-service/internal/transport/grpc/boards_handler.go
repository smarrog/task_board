package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smarrog/task-board/shared/proto/v1"

	boardrepo "github.com/smarrog/task-board/core-service/internal/repo/boards"
	boardsvc "github.com/smarrog/task-board/core-service/internal/service/boards"
)

type BoardsHandler struct {
	v1.UnimplementedBoardsServiceServer

	log *zerolog.Logger
	svc *boardsvc.Service
}

func NewBoardsHandler(log *zerolog.Logger, svc *boardsvc.Service) *BoardsHandler {
	return &BoardsHandler{
		log: log,
		svc: svc,
	}
}

func (h *BoardsHandler) CreateBoard(ctx context.Context, req *v1.CreateBoardRequest) (*v1.CreateBoardResponse, error) {
	ownerId, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is invalid")
	}

	b, err := h.svc.Create(ctx, ownerId, req.GetTitle(), req.GetDescription())
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.CreateBoardResponse{
		Board: toProtoBoard(b),
	}, nil
}

func (h *BoardsHandler) GetBoard(ctx context.Context, req *v1.GetBoardRequest) (*v1.GetBoardResponse, error) {
	boardID, err := uuid.Parse(req.GetBoardId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "board_id is invalid")
	}

	b, err := h.svc.Get(ctx, boardID)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.GetBoardResponse{Board: toProtoBoard(b)}, nil
}

func (h *BoardsHandler) ListBoards(ctx context.Context, req *v1.ListBoardsRequest) (*v1.ListBoardsResponse, error) {
	ownerID, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is invalid")
	}

	boards, err := h.svc.List(ctx, ownerID)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	out := make([]*v1.Board, 0, len(boards))
	for i := range boards {
		b := boards[i]
		out = append(out, toProtoBoard(&b))
	}

	return &v1.ListBoardsResponse{Boards: out}, nil
}

func (h *BoardsHandler) UpdateBoard(ctx context.Context, req *v1.UpdateBoardRequest) (*v1.UpdateBoardResponse, error) {
	ownerID, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is invalid")
	}

	boardID, err := uuid.Parse(req.GetBoardId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "board_id is invalid")
	}

	if err := h.svc.Update(ctx, ownerID, boardID, req.GetTitle(), req.GetDescription()); err != nil {
		return nil, mapBoardsErr(err)
	}

	b, err := h.svc.Get(ctx, boardID)
	if err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.UpdateBoardResponse{Board: toProtoBoard(b)}, nil
}

func (h *BoardsHandler) DeleteBoard(ctx context.Context, req *v1.DeleteBoardRequest) (*v1.DeleteBoardResponse, error) {
	boardId, err := uuid.Parse(req.GetBoardId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "board_id is invalid")
	}

	if err := h.svc.Delete(ctx, boardId); err != nil {
		return nil, mapBoardsErr(err)
	}

	return &v1.DeleteBoardResponse{}, nil
}

func toProtoBoard(b *model.Board) *v1.Board {
	return &v1.Board{
		Id:          b.Id.String(),
		OwnerId:     b.OwnerId.String(),
		Title:       b.Title,
		Description: b.Description,
	}
}

func mapBoardsErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, boardsvc.ErrInvalidArgument) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, boardrepo.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, "internal error")
}
