package http

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
)

type createTaskBody struct {
	Position    int32  `json:"position"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee_id"`
}

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	columnID := c.Params("columnId")
	var body createTaskBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}

	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.tasks.CreateTask(ctx, &v1.CreateTaskRequest{
		Base:        &v1.BaseRequest{RequesterId: h.requesterID(c)},
		ColumnId:    columnID,
		Position:    body.Position,
		Title:       body.Title,
		Description: body.Description,
		AssigneeId:  body.AssigneeID,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	t := resp.GetTask()
	return c.Status(fiber.StatusCreated).JSON(TaskDTO{
		Id:          t.GetId(),
		ColumnId:    t.GetColumnId(),
		Position:    t.GetPosition(),
		Title:       t.GetTitle(),
		Description: t.GetDescription(),
		AssigneeId:  t.GetAssigneeId(),
	})
}

func (h *Handler) GetTask(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.tasks.GetTask(ctx, &v1.GetTaskRequest{
		Base:   &v1.BaseRequest{RequesterId: h.requesterID(c)},
		TaskId: taskID,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	t := resp.GetTask()
	return c.JSON(TaskDTO{
		Id:          t.GetId(),
		ColumnId:    t.GetColumnId(),
		Position:    t.GetPosition(),
		Title:       t.GetTitle(),
		Description: t.GetDescription(),
		AssigneeId:  t.GetAssigneeId(),
	})
}

type updateTaskBody struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	AssigneeID  *string `json:"assignee_id"`
}

func (h *Handler) UpdateTask(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	var body updateTaskBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}

	title := ""
	desc := ""
	assignee := ""
	if body.Title != nil {
		title = *body.Title
	}
	if body.Description != nil {
		desc = *body.Description
	}
	if body.AssigneeID != nil {
		assignee = *body.AssigneeID
	}

	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.tasks.UpdateTask(ctx, &v1.UpdateTaskRequest{
		Base:        &v1.BaseRequest{RequesterId: h.requesterID(c)},
		TaskId:      taskID,
		Title:       title,
		Description: desc,
		AssigneeId:  assignee,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	t := resp.GetTask()
	return c.JSON(TaskDTO{
		Id:          t.GetId(),
		ColumnId:    t.GetColumnId(),
		Position:    t.GetPosition(),
		Title:       t.GetTitle(),
		Description: t.GetDescription(),
		AssigneeId:  t.GetAssigneeId(),
	})
}

type moveTaskBody struct {
	ToColumnID string `json:"to_column_id"`
	ToPosition int32  `json:"to_position"`
}

func (h *Handler) MoveTask(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	var body moveTaskBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}

	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.tasks.MoveTask(ctx, &v1.MoveTaskRequest{
		Base:       &v1.BaseRequest{RequesterId: h.requesterID(c)},
		TaskId:     taskID,
		ToColumnId: body.ToColumnID,
		ToPosition: body.ToPosition,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	t := resp.GetTask()
	return c.JSON(TaskDTO{
		Id:          t.GetId(),
		ColumnId:    t.GetColumnId(),
		Position:    t.GetPosition(),
		Title:       t.GetTitle(),
		Description: t.GetDescription(),
		AssigneeId:  t.GetAssigneeId(),
	})
}

func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	_, err := h.tasks.DeleteTask(ctx, &v1.DeleteTaskRequest{
		Base:   &v1.BaseRequest{RequesterId: h.requesterID(c)},
		TaskId: taskID,
	})
	if err != nil {
		return grpcToHTTP(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
