package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type DeleteBoardUseCase struct {
	repo board.Repository
}

type DeleteBoardInput struct {
	BoardId string
}

type DeleteBoardOutput struct {
}

func NewDeleteBoardUseCase(repo board.Repository) *DeleteBoardUseCase {
	return &DeleteBoardUseCase{repo: repo}
}

func (uc *DeleteBoardUseCase) Execute(ctx context.Context, input DeleteBoardInput) (*DeleteBoardOutput, error) {
	id, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, fmt.Errorf("board_id: %w", err)
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete board: %w", err)
	}

	output := &DeleteBoardOutput{}

	return output, nil
}
