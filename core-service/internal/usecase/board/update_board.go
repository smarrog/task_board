package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type UpdateBoardUseCase struct {
	repo board.Repository
}

type UpdateBoardInput struct {
	BoardId     string
	OwnerId     string
	Title       string
	Description string
}

type UpdateBoardOutput struct {
	Board *board.Board
}

func NewUpdateBoardUseCase(repo board.Repository) *UpdateBoardUseCase {
	return &UpdateBoardUseCase{repo: repo}
}

func (uc *UpdateBoardUseCase) Execute(ctx context.Context, input UpdateBoardInput) (*UpdateBoardOutput, error) {
	bid, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}

	oid, err := common.UserIdFromString(input.Title)
	if err != nil {
		return nil, fmt.Errorf("board owner_id: %w", err)
	}

	t, err := board.NewTitle(input.Title)
	if err != nil {
		return nil, fmt.Errorf("board: %w", err)
	}

	d, err := board.NewDescription(input.Description)
	if err != nil {
		return nil, fmt.Errorf("board: %w", err)
	}

	b, err := uc.repo.Get(ctx, bid)
	if err != nil {
		return nil, err
	}

	if b.OwnerId().UUID() != oid.UUID() {
		return nil, board.ErrOwnerMismatch
	}

	b.Update(t, d)

	err = uc.repo.Save(ctx, b)
	if err != nil {
		return nil, fmt.Errorf("save board: %w", err)
	}

	output := &UpdateBoardOutput{
		Board: b,
	}

	return output, nil
}
