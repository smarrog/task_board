package column

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type DeleteColumnUseCase struct {
	repo  column.Repository
	cache cache.Invalidator
}

type DeleteColumnInput struct {
	ColumnId string
}

type DeleteColumnOutput struct {
}

func NewDeleteColumnUseCase(repo column.Repository, cache cache.Invalidator) *DeleteColumnUseCase {
	return &DeleteColumnUseCase{repo: repo, cache: cache}
}

func (uc *DeleteColumnUseCase) Execute(ctx context.Context, input DeleteColumnInput) (*DeleteColumnOutput, error) {
	cid, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}

	c, err := uc.repo.Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Delete(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("delete column: %w", err)
	}

	if uc.cache != nil {
		_ = uc.cache.InvalidateBoard(ctx, c.BoardId())
	}

	output := &DeleteColumnOutput{}

	return output, nil
}
