package column

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type CreateColumnUseCase struct {
	repo  column.Repository
	cache cache.Invalidator
}

type CreateColumnInput struct {
	BoardId  string
	Position int
}

type CreateColumnOutput struct {
	Column *column.Column
}

func NewCreateColumnUseCase(repo column.Repository, cache cache.Invalidator) *CreateColumnUseCase {
	return &CreateColumnUseCase{repo: repo, cache: cache}
}

func (uc *CreateColumnUseCase) Execute(ctx context.Context, input CreateColumnInput) (*CreateColumnOutput, error) {
	bid, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}
	position, err := column.NewPosition(input.Position)
	if err != nil {
		return nil, err
	}

	c := column.New(bid, position)

	err = uc.repo.Save(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("save column: %w", err)
	}

	if uc.cache != nil {
		_ = uc.cache.InvalidateBoard(ctx, bid)
	}

	output := &CreateColumnOutput{
		Column: c,
	}

	return output, nil
}
