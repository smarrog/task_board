package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
	"github.com/smarrog/task-board/core-service/internal/usecase/common"
)

type MoveTaskUseCase struct {
	repo    task.Repository
	columns column.Repository
	cache   cache.Invalidator
}

type MoveTaskInput struct {
	TaskId     string
	ToColumnId string
	ToPosition int
}

type MoveTaskOutput struct {
	Task *task.Task
}

func NewMoveTaskUseCase(repo task.Repository, columns column.Repository, cache cache.Invalidator) *MoveTaskUseCase {
	return &MoveTaskUseCase{repo: repo, columns: columns, cache: cache}
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
	var fromColId column.Id
	var toColId = toCol

	err = uc.repo.InTx(ctx, func(ctx context.Context) error {
		t, err := uc.repo.Get(ctx, tid)
		if err != nil {
			return err
		}

		fromCol := t.ColumnId()
		fromColId = fromCol
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

			clampedToPos := common.Clamp(int(toPos), 0, n-1)
			shift, needShift := common.CalcShift(fromPos, clampedToPos)
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

		insertPos := common.Clamp(int(toPos), 0, nd)

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

	if uc.cache != nil {
		if c, err := uc.columns.Get(ctx, fromColId); err == nil {
			_ = uc.cache.InvalidateBoard(ctx, c.BoardId())
		}
		if c, err := uc.columns.Get(ctx, toColId); err == nil {
			_ = uc.cache.InvalidateBoard(ctx, c.BoardId())
		}
	}

	return out, nil
}
