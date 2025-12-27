package column

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
)

type CreateColumnUseCase struct {
	repo column.Repository
}

type CreateColumnInput struct {
	BoardId  string
	Position int
}

type CreateColumnOutput struct {
	Column *column.Column
}

func NewCreateColumnUseCase(repo column.Repository) *CreateColumnUseCase {
	return &CreateColumnUseCase{repo: repo}
}

func (uc *CreateColumnUseCase) Execute(ctx context.Context, input CreateColumnInput) (*CreateColumnOutput, error) {
	bid, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}
	position := input.Position

	c := column.New(bid, position)

	err = uc.repo.Save(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("save column: %w", err)
	}

	output := &CreateColumnOutput{
		Column: c,
	}

	return output, nil
}
