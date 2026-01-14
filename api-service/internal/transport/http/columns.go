package http

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
)

type createColumnBody struct {
	Position int32 `json:"position"`
}

func (h *Handler) CreateColumn(c *fiber.Ctx) error {
	boardId := c.Params("boardId")
	var body createColumnBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.columns.CreateColumn(ctx, &v1.CreateColumnRequest{
		Base:     &v1.BaseRequest{RequesterId: h.requesterID(c)},
		BoardId:  boardId,
		Position: body.Position,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	col := resp.GetColumn()
	return c.Status(fiber.StatusCreated).JSON(ColumnDTO{
		Id:       col.GetId(),
		BoardId:  col.GetBoardId(),
		Position: col.GetPosition(),
		Tasks:    []TaskDTO{},
	})
}

func (h *Handler) GetColumn(c *fiber.Ctx) error {
	columnId := c.Params("columnId")
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.columns.GetColumnFull(ctx, &v1.GetColumnFullRequest{
		Base:     &v1.BaseRequest{RequesterId: h.requesterID(c)},
		ColumnId: columnId,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	columnDTO := toColumnDTO(resp.GetData())
	return c.JSON(columnDTO)
}

type moveColumnBody struct {
	ToPosition int32 `json:"to_position"`
}

func (h *Handler) MoveColumn(c *fiber.Ctx) error {
	columnId := c.Params("columnId")
	var body moveColumnBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	_, err := h.columns.MoveColumn(ctx, &v1.MoveColumnRequest{
		Base:       &v1.BaseRequest{RequesterId: h.requesterID(c)},
		ColumnId:   columnId,
		ToPosition: body.ToPosition,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	resp, err := h.columns.GetColumnFull(ctx, &v1.GetColumnFullRequest{
		Base:     &v1.BaseRequest{RequesterId: h.requesterID(c)},
		ColumnId: columnId,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	columnDTO := toColumnDTO(resp.GetData())
	return c.JSON(columnDTO)
}

func (h *Handler) DeleteColumn(c *fiber.Ctx) error {
	columnId := c.Params("columnId")
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	_, err := h.columns.DeleteColumn(ctx, &v1.DeleteColumnRequest{
		Base:     &v1.BaseRequest{RequesterId: h.requesterID(c)},
		ColumnId: columnId,
	})
	if err != nil {
		return grpcToHTTP(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func toColumnDTO(full *v1.ColumnFull) ColumnDTO {
	col := full.GetColumn()
	tasks := make([]TaskDTO, 0, len(full.GetTasks()))
	for _, t := range full.GetTasks() {
		tasks = append(tasks, TaskDTO{
			Id:          t.GetId(),
			ColumnId:    t.GetColumnId(),
			Position:    t.GetPosition(),
			Title:       t.GetTitle(),
			Description: t.GetDescription(),
			AssigneeId:  t.GetAssigneeId(),
		})
	}

	return ColumnDTO{
		Id:       col.GetId(),
		BoardId:  col.GetBoardId(),
		Position: col.GetPosition(),
		Tasks:    tasks,
	}
}
