package column

import (
	"context"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
)

type DeleteColumnUseCase struct {
	repo column.Repository
}

type DeleteColumnInput struct {
	ColumnId string
}

type DeleteColumnOutput struct {
}

func NewDeleteColumnUseCase(repo column.Repository) *DeleteColumnUseCase {
	return &DeleteColumnUseCase{repo: repo}
}

func (uc *DeleteColumnUseCase) Execute(ctx context.Context, input DeleteColumnInput) (*DeleteColumnOutput, error) {
	cid, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Delete(ctx, cid)
	if err != nil {
		return nil, err
	}

	output := &DeleteColumnOutput{}

	return output, nil
}
