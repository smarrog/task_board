package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type ListBoardsUseCase struct {
	repo board.Repository
}

type ListBoardsInput struct {
	OwnerId string
}

type ListBoardsOutput struct {
	Boards []*board.Board
}

func NewListBoardsUseCase(repo board.Repository) *ListBoardsUseCase {
	return &ListBoardsUseCase{repo: repo}
}

func (uc *ListBoardsUseCase) Execute(ctx context.Context, input ListBoardsInput) (*ListBoardsOutput, error) {
	oid, err := common.UserIdFromString(input.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("owner_id: %w", err)
	}

	l, err := uc.repo.ListByOwner(ctx, oid)
	if err != nil {
		return nil, err
	}

	output := &ListBoardsOutput{
		Boards: l,
	}

	return output, nil
}
