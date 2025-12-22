package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type GetBoardUseCase struct {
	repo board.Repository
}

func NewGetBoardUseCase(repo board.Repository) *GetBoardUseCase {
	return &GetBoardUseCase{repo: repo}
}

func (uc *GetBoardUseCase) Execute(ctx context.Context, boardId string) (*board.Board, error) {
	id, err := board.BoardIdFromString(boardId)
	if err != nil {
		return nil, fmt.Errorf("board_id: %w", err)
	}

	b, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	return b, nil
}
