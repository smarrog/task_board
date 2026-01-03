package column

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	uccommon "github.com/smarrog/task-board/core-service/internal/usecase/common"
)

type MoveColumnUseCase struct {
	repo column.Repository
}

type MoveColumnInput struct {
	ColumnId   string
	ToPosition int
}

type MoveColumnOutput struct {
	Column *column.Column
}

func NewMoveColumnUseCase(repo column.Repository) *MoveColumnUseCase {
	return &MoveColumnUseCase{repo: repo}
}

func (uc *MoveColumnUseCase) Execute(ctx context.Context, input MoveColumnInput) (*MoveColumnOutput, error) {
	cid, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}

	toPos, err := column.NewPosition(input.ToPosition)
	if err != nil {
		return nil, err
	}

	var out *MoveColumnOutput

	err = uc.repo.InTx(ctx, func(ctx context.Context) error {
		c, err := uc.repo.Get(ctx, cid)
		if err != nil {
			return err
		}

		bid := c.BoardId()
		fromPos := int(c.Position())

		if err := uc.repo.LockBoardColumns(ctx, bid); err != nil {
			return err
		}

		n, err := uc.repo.CountInBoard(ctx, bid)
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("invariant: board has no columns")
		}

		clampedToPos := uccommon.Clamp(int(toPos), 0, n-1)
		shift, needShift := uccommon.CalcShift(fromPos, clampedToPos)
		if needShift {
			if err := uc.repo.ShiftPositions(ctx, bid, shift.FromPosition, shift.ToPosition, shift.Delta); err != nil {
				return err
			}
		}

		c.Move(column.Position(clampedToPos))

		if err := uc.repo.Save(ctx, c); err != nil {
			return err
		}

		out = &MoveColumnOutput{Column: c}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
