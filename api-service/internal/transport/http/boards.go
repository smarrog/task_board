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

	return c.Status(fiber.StatusCreated).JSON(BoardDTO{
		Id:          resp.GetBoard().GetId(),
		OwnerId:     resp.GetBoard().GetOwnerId(),
		Title:       resp.GetBoard().GetTitle(),
		Description: resp.GetBoard().GetDescription(),
	})
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

	b := resp.GetBoard()
	return c.JSON(BoardDTO{
		Id:          b.GetId(),
		OwnerId:     b.GetOwnerId(),
		Title:       b.GetTitle(),
		Description: b.GetDescription(),
	})
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

	out := make([]BoardDTO, 0, len(resp.GetBoards()))
	for _, b := range resp.GetBoards() {
		out = append(out, BoardDTO{
			Id:          b.GetId(),
			OwnerId:     b.GetOwnerId(),
			Title:       b.GetTitle(),
			Description: b.GetDescription(),
		})
	}
	return c.JSON(out)
}

type updateBoardBody struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (h *Handler) UpdateBoard(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	var body updateBoardBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}

	ctx, cancel := h.reqCtx(c)
	defer cancel()

	// Core service expects full values, so send empty string if not provided.
	title := ""
	desc := ""
	if body.Title != nil {
		title = *body.Title
	}
	if body.Description != nil {
		desc = *body.Description
	}

	resp, err := h.boards.UpdateBoard(ctx, &v1.UpdateBoardRequest{
		Base:        &v1.BaseRequest{RequesterId: h.requesterID(c)},
		BoardId:     boardID,
		OwnerId:     h.requesterID(c),
		Title:       title,
		Description: desc,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	b := resp.GetBoard()
	return c.JSON(BoardDTO{
		Id:          b.GetId(),
		OwnerId:     b.GetOwnerId(),
		Title:       b.GetTitle(),
		Description: b.GetDescription(),
	})
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
