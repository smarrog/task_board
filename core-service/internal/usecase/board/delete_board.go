package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type DeleteBoardUseCase struct {
	repo board.Repository
}

func NewDeleteBoardUseCase(repo board.Repository) *DeleteBoardUseCase {
	return &DeleteBoardUseCase{repo: repo}
}

func (uc *DeleteBoardUseCase) Execute(ctx context.Context, boardId string) error {
	id, err := board.BoardIdFromString(boardId)
	if err != nil {
		return fmt.Errorf("board_id: %w", err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete board: %w", err)
	}

	return nil
}
