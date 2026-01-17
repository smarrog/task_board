package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/shared/domain/shared"
)

type CreateBoardUseCase struct {
	repo board.Repository
}

type CreateBoardInput struct {
	OwnerId     string
	Title       string
	Description string
}

type CreateBoardOutput struct {
	Board *board.Board
}

func NewCreateBoardUseCase(repo board.Repository) *CreateBoardUseCase {
	return &CreateBoardUseCase{repo: repo}
}

func (uc *CreateBoardUseCase) Execute(ctx context.Context, input CreateBoardInput) (*CreateBoardOutput, error) {
	userId, err := shared.UserIdFromString(input.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("board owner_id: %w", err)
	}

	t, err := board.NewTitle(input.Title)
	if err != nil {
		return nil, err
	}

	d, err := board.NewDescription(input.Description)
	if err != nil {
		return nil, err
	}

	b, err := board.New(userId, t, d)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Save(ctx, b)
	if err != nil {
		return nil, fmt.Errorf("save board: %w", err)
	}

	output := &CreateBoardOutput{
		Board: b,
	}

	return output, nil
}
