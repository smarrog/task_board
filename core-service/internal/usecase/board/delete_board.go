package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type DeleteBoardUseCase struct {
	repo  board.Repository
	cache cache.Invalidator
}

type DeleteBoardInput struct {
	BoardId string
}

type DeleteBoardOutput struct {
}

func NewDeleteBoardUseCase(repo board.Repository, cache cache.Invalidator) *DeleteBoardUseCase {
	return &DeleteBoardUseCase{repo: repo, cache: cache}
}

func (uc *DeleteBoardUseCase) Execute(ctx context.Context, input DeleteBoardInput) (*DeleteBoardOutput, error) {
	id, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete board: %w", err)
	}

	if uc.cache != nil {
		_ = uc.cache.InvalidateBoard(ctx, id)
	}

	output := &DeleteBoardOutput{}

	return output, nil
}
