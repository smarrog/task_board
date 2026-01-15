package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	commonuc "github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

const (
	boardPrefix = "board:"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (c *RedisCache) GetBoard(ctx context.Context, id board.Id) (out *commonuc.BoardData, hit bool, err error) {
	if c == nil || c.rdb == nil {
		return nil, false, nil
	}

	raw, err := c.rdb.Get(ctx, c.boardKey(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var dto getBoardDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		// corrupted entry -> treat as miss
		_ = c.rdb.Del(ctx, c.boardKey(id)).Err()
		return nil, false, nil
	}

	out, err = dto.toOutput()
	if err != nil {
		_ = c.rdb.Del(ctx, c.boardKey(id)).Err()
		return nil, false, nil
	}

	return out, true, nil
}

func (c *RedisCache) SetBoard(ctx context.Context, id board.Id, out *commonuc.BoardData, ttl time.Duration) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	if out == nil || out.Board == nil {
		return nil
	}

	dto, err := dtoFromOutput(out)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	if ttl <= 0 {
		return c.rdb.Set(ctx, c.boardKey(id), raw, 0).Err()
	}
	return c.rdb.Set(ctx, c.boardKey(id), raw, ttl).Err()
}

func (c *RedisCache) InvalidateBoard(ctx context.Context, id board.Id) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Del(ctx, c.boardKey(id)).Err()
}

func (c *RedisCache) boardKey(id board.Id) string { return boardPrefix + id.String() }

type getBoardDTO struct {
	Board   boardDTO    `json:"board"`
	Columns []columnDTO `json:"columns"`
	Tasks   []taskDTO   `json:"tasks"`
}

type boardDTO struct {
	Id          string    `json:"id"`
	OwnerId     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type columnDTO struct {
	Id        string    `json:"id"`
	BoardId   string    `json:"board_id"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type taskDTO struct {
	Id          string    `json:"id"`
	ColumnId    string    `json:"column_id"`
	Position    int       `json:"position"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeId  string    `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func dtoFromOutput(out *commonuc.BoardData) (*getBoardDTO, error) {
	b := out.Board
	dto := &getBoardDTO{
		Board: boardDTO{
			Id:          b.Id().String(),
			OwnerId:     b.OwnerId().String(),
			Title:       b.Title().String(),
			Description: b.Description().String(),
			CreatedAt:   b.CreatedAt(),
			UpdatedAt:   b.UpdatedAt(),
		},
		Columns: make([]columnDTO, 0, len(out.Columns)),
		Tasks:   make([]taskDTO, 0, len(out.Tasks)),
	}

	for _, c := range out.Columns {
		dto.Columns = append(dto.Columns, columnDTO{
			Id:        c.Id().String(),
			BoardId:   c.BoardId().String(),
			Position:  int(c.Position()),
			CreatedAt: c.CreatedAt(),
			UpdatedAt: c.UpdatedAt(),
		})
	}

	for _, t := range out.Tasks {
		dto.Tasks = append(dto.Tasks, taskDTO{
			Id:          t.Id().String(),
			ColumnId:    t.ColumnId().String(),
			Position:    int(t.Position()),
			Title:       t.Title().String(),
			Description: t.Description().String(),
			AssigneeId:  t.AssigneeId().String(),
			CreatedAt:   t.CreatedAt(),
			UpdatedAt:   t.UpdatedAt(),
		})
	}

	return dto, nil
}

func (d getBoardDTO) toOutput() (*commonuc.BoardData, error) {
	bid, err := board.IdFromString(d.Board.Id)
	if err != nil {
		return nil, err
	}
	oid, err := common.UserIdFromString(d.Board.OwnerId)
	if err != nil {
		return nil, err
	}
	bt, err := board.NewTitle(d.Board.Title)
	if err != nil {
		return nil, err
	}
	bd, err := board.NewDescription(d.Board.Description)
	if err != nil {
		return nil, err
	}

	b := board.Rehydrate(bid, oid, bt, bd, d.Board.CreatedAt, d.Board.UpdatedAt)

	cols := make([]*column.Column, 0, len(d.Columns))
	for _, c := range d.Columns {
		cid, err := column.IdFromString(c.Id)
		if err != nil {
			return nil, err
		}
		cb, err := board.IdFromString(c.BoardId)
		if err != nil {
			return nil, err
		}
		pos, err := column.NewPosition(c.Position)
		if err != nil {
			return nil, err
		}
		cols = append(cols, column.Rehydrate(cid, cb, pos, c.CreatedAt, c.UpdatedAt))
	}

	tasksOut := make([]*task.Task, 0, len(d.Tasks))
	for _, t := range d.Tasks {
		tid, err := task.IdFromString(t.Id)
		if err != nil {
			return nil, err
		}
		tc, err := column.IdFromString(t.ColumnId)
		if err != nil {
			return nil, err
		}
		pos, err := task.NewPosition(t.Position)
		if err != nil {
			return nil, err
		}
		tt, err := task.NewTitle(t.Title)
		if err != nil {
			return nil, err
		}
		td, err := task.NewDescription(t.Description)
		if err != nil {
			return nil, err
		}
		aid, err := common.UserIdFromString(t.AssigneeId)
		if err != nil {
			return nil, fmt.Errorf("assignee_id: %w", err)
		}

		tasksOut = append(tasksOut, task.Rehydrate(tid, tc, pos, tt, td, aid, t.CreatedAt, t.UpdatedAt))
	}

	return &commonuc.BoardData{Board: b, Columns: cols, Tasks: tasksOut}, nil
}
