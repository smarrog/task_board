package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type GetBoardUseCase struct {
	repo board.Repository
}

type GetBoardInput struct {
	BoardId string
}

type GetBoardOutput struct {
	Board *board.Board
}

func NewGetBoardUseCase(repo board.Repository) *GetBoardUseCase {
	return &GetBoardUseCase{repo: repo}
}

func (uc *GetBoardUseCase) Execute(ctx context.Context, input GetBoardInput) (*GetBoardOutput, error) {
	id, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}

	b, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	output := &GetBoardOutput{
		Board: b,
	}

	return output, nil
}
