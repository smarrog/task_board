package http

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
)

type createBoardBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) CreateBoard(c *fiber.Ctx) error {
	var body createBoardBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.boards.CreateBoard(ctx, &v1.CreateBoardRequest{
		Base:        &v1.BaseRequest{RequesterId: h.requesterID(c)},
		Title:       body.Title,
		Description: body.Description,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	boardDTO := buildBoardDTO(resp.GetData())
	return c.Status(fiber.StatusCreated).JSON(boardDTO)
}

func (h *Handler) GetBoard(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.boards.GetBoard(ctx, &v1.GetBoardRequest{
		Base:    &v1.BaseRequest{RequesterId: h.requesterID(c)},
		BoardId: boardID,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	boardDTO := buildBoardDTO(resp.GetData())
	return c.JSON(boardDTO)
}

func (h *Handler) ListBoards(c *fiber.Ctx) error {
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	ownerID := h.requesterID(c)
	resp, err := h.boards.ListBoards(ctx, &v1.ListBoardsRequest{
		Base:    &v1.BaseRequest{RequesterId: ownerID},
		OwnerId: ownerID,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	boards := resp.GetBoards()
	out := make([]BoardDTO, 0, len(boards))
	for _, b := range boards {
		out = append(out, buildBoardDTO(b))
	}
	return c.JSON(out)
}

type updateBoardBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) UpdateBoard(c *fiber.Ctx) error {
	boardID := c.Params("boardId")

	var body updateBoardBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}

	ctx, cancel := h.reqCtx(c)
	defer cancel()

	resp, err := h.boards.UpdateBoard(ctx, &v1.UpdateBoardRequest{
		Base:        &v1.BaseRequest{RequesterId: h.requesterID(c)},
		BoardId:     boardID,
		OwnerId:     h.requesterID(c),
		Title:       body.Title,
		Description: body.Description,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	boardDTO := buildBoardDTO(resp.GetData())
	return c.JSON(boardDTO)
}

func (h *Handler) DeleteBoard(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	ctx, cancel := h.reqCtx(c)
	defer cancel()

	_, err := h.boards.DeleteBoard(ctx, &v1.DeleteBoardRequest{
		Base:    &v1.BaseRequest{RequesterId: h.requesterID(c)},
		BoardId: boardID,
	})
	if err != nil {
		return grpcToHTTP(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func buildBoardDTO(full *v1.BoardFull) BoardDTO {
	b := full.GetBoard()
	cols := make([]ColumnDTO, 0, len(full.GetColumns()))
	for _, cwt := range full.GetColumns() {
		c := cwt.GetColumn()
		tasks := make([]TaskDTO, 0, len(cwt.GetTasks()))
		for _, t := range cwt.GetTasks() {
			tasks = append(tasks, TaskDTO{
				Id:          t.GetId(),
				ColumnId:    t.GetColumnId(),
				Position:    t.GetPosition(),
				Title:       t.GetTitle(),
				Description: t.GetDescription(),
				AssigneeId:  t.GetAssigneeId(),
			})
		}
		cols = append(cols, ColumnDTO{
			Id:       c.GetId(),
			BoardId:  c.GetBoardId(),
			Position: c.GetPosition(),
			Tasks:    tasks,
		})
	}

	return BoardDTO{
		Id:          b.GetId(),
		OwnerId:     b.GetOwnerId(),
		Title:       b.GetTitle(),
		Description: b.GetDescription(),
		Columns:     cols,
	}
}
