package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type CreateBoardUseCase struct {
	repo board.Repository
}

func NewCreateBoardUseCase(repo board.Repository) *CreateBoardUseCase {
	return &CreateBoardUseCase{repo: repo}
}

func (uc *CreateBoardUseCase) Execute(ctx context.Context, ownerId, title, description string) (*board.Board, error) {
	userId, err := board.UserIdFromString(ownerId)
	if err != nil {
		return nil, fmt.Errorf("owner_id: %w", err)
	}

	t, err := board.NewTitle(title)
	if err != nil {
		return nil, fmt.Errorf("title: %w", err)
	}

	d, err := board.NewDescription(description)
	if err != nil {
		return nil, fmt.Errorf("description: %w", err)
	}

	b, err := board.NewBoard(userId, t, d)
	if err != nil {
		return nil, fmt.Errorf("create board: %w", err)
	}

	if err := uc.repo.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("save board: %w", err)
	}

	return b, nil
}
