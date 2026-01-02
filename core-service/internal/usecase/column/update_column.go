package column

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
)

type UpdateColumnUseCase struct {
	repo column.Repository
}

type UpdateColumnInput struct {
	ColumnId string
	BoardId  string
	Position int
}

type UpdateColumnOutput struct {
	Column *column.Column
}

func NewUpdateColumnUseCase(repo column.Repository) *UpdateColumnUseCase {
	return &UpdateColumnUseCase{repo: repo}
}

func (uc *UpdateColumnUseCase) Execute(ctx context.Context, input UpdateColumnInput) (output *UpdateColumnOutput, err error) {
	cid, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}
	bid, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}
	position, err := column.NewPosition(input.Position)
	if err != nil {
		return nil, err
	}

	c, err := uc.repo.Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	c.Update(bid, position)

	err = uc.repo.Save(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("save column: %w", err)
	}

	output = &UpdateColumnOutput{
		Column: c,
	}

	return output, nil
}
