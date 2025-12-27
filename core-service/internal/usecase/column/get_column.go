package column

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
)

type GetColumnUseCase struct {
	repo column.Repository
}

type GetColumnInput struct {
	ColumnId string
}

type GetColumnOutput struct {
	Column *column.Column
}

func NewGetColumnUseCase(repo column.Repository) *GetColumnUseCase {
	return &GetColumnUseCase{repo: repo}
}

func (uc *GetColumnUseCase) Execute(ctx context.Context, input GetColumnInput) (*GetColumnOutput, error) {
	id, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}

	c, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get column: %w", err)
	}

	output := &GetColumnOutput{
		Column: c,
	}

	return output, nil
}
