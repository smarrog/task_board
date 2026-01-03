package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	uccommon "github.com/smarrog/task-board/core-service/internal/usecase/common"
)

type MoveTaskUseCase struct {
	repo task.Repository
}

type MoveTaskInput struct {
	TaskId     string
	ToColumnId string
	ToPosition int
}

type MoveTaskOutput struct {
	Task *task.Task
}

func NewMoveTaskUseCase(repo task.Repository) *MoveTaskUseCase {
	return &MoveTaskUseCase{repo: repo}
}

func (uc *MoveTaskUseCase) Execute(ctx context.Context, input MoveTaskInput) (*MoveTaskOutput, error) {
	tid, err := task.IdFromString(input.TaskId)
	if err != nil {
		return nil, err
	}
	toCol, err := column.IdFromString(input.ToColumnId)
	if err != nil {
		return nil, err
	}
	toPos, err := task.NewPosition(input.ToPosition)
	if err != nil {
		return nil, err
	}

	var out *MoveTaskOutput

	err = uc.repo.InTx(ctx, func(ctx context.Context) error {
		t, err := uc.repo.Get(ctx, tid)
		if err != nil {
			return err
		}

		fromCol := t.ColumnId()
		fromPos := int(t.Position())

		// смена позиции внутри одной колонки
		if fromCol == toCol {
			if err := uc.repo.LockColumnTasks(ctx, toCol); err != nil {
				return err
			}

			n, err := uc.repo.CountInColumn(ctx, toCol)
			if err != nil {
				return err
			}
			if n == 0 {
				return fmt.Errorf("column has no tasks")
			}

			clampedToPos := uccommon.Clamp(int(toPos), 0, n-1)
			shift, needShift := uccommon.CalcShift(fromPos, clampedToPos)
			if needShift {
				if err := uc.repo.ShiftPositions(ctx, toCol, shift.FromPosition, shift.ToPosition, shift.Delta); err != nil {
					return err
				}
			}

			t.Move(toCol, task.Position(clampedToPos))
			if err := uc.repo.Save(ctx, t); err != nil {
				return err
			}

			out = &MoveTaskOutput{Task: t}
			return nil
		}

		// перенос между колонками
		first, second := fromCol, toCol
		if first.UUID().String() > second.UUID().String() {
			first, second = second, first
		}
		if err := uc.repo.LockColumnTasks(ctx, first); err != nil {
			return err
		}
		if err := uc.repo.LockColumnTasks(ctx, second); err != nil {
			return err
		}

		nd, err := uc.repo.CountInColumn(ctx, toCol)
		if err != nil {
			return err
		}

		insertPos := uccommon.Clamp(int(toPos), 0, nd)

		if err := uc.repo.ShiftAfterRemove(ctx, fromCol, fromPos); err != nil {
			return err
		}
		if err := uc.repo.ShiftForInsert(ctx, toCol, insertPos); err != nil {
			return err
		}

		t.Move(toCol, task.Position(insertPos))
		if err := uc.repo.Save(ctx, t); err != nil {
			return err
		}

		out = &MoveTaskOutput{Task: t}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
