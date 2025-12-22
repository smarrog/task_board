package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type UpdateBoardUseCase struct {
	repo board.Repository
}

func NewUpdateBoardUseCase(repo board.Repository) *UpdateBoardUseCase {
	return &UpdateBoardUseCase{repo: repo}
}

func (uc *UpdateBoardUseCase) Execute(ctx context.Context, ownerID, boardID, title, description string) (*board.Board, error) {
	oid, err := board.UserIdFromString(ownerID)
	if err != nil {
		return nil, fmt.Errorf("owner_id: %w", err)
	}

	bid, err := board.BoardIdFromString(boardID)
	if err != nil {
		return nil, fmt.Errorf("board_id: %w", err)
	}

	t, err := board.NewTitle(title)
	if err != nil {
		return nil, fmt.Errorf("title: %w", err)
	}

	d, err := board.NewDescription(description)
	if err != nil {
		return nil, fmt.Errorf("description: %w", err)
	}

	b, err := uc.repo.Get(ctx, bid)
	if err != nil {
		return nil, err
	}

	if b.OwnerId().UUID() != oid.UUID() {
		return nil, board.ErrOwnerMismatch
	}

	b.Update(t, d)
	if err := uc.repo.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("save board: %w", err)
	}
	return b, nil
}
