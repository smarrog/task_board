package http

import "github.com/gofiber/fiber/v2"

func (h *Handler) Register(r fiber.Router) {
	// Boards
	r.Post("/boards", h.CreateBoard)
	r.Get("/boards", h.ListBoards)
	r.Get("/boards/:boardId", h.GetBoard)
	r.Put("/boards/:boardId", h.UpdateBoard)
	r.Delete("/boards/:boardId", h.DeleteBoard)

	// Columns
	r.Post("/boards/:boardId/columns", h.CreateColumn)
	r.Get("/columns/:columnId", h.GetColumn)
	r.Post("/columns/:columnId/move", h.MoveColumn)
	r.Delete("/columns/:columnId", h.DeleteColumn)

	// Tasks
	r.Post("/columns/:columnId/tasks", h.CreateTask)
	r.Get("/tasks/:taskId", h.GetTask)
	r.Put("/tasks/:taskId", h.UpdateTask)
	r.Post("/tasks/:taskId/move", h.MoveTask)
	r.Delete("/tasks/:taskId", h.DeleteTask)
}
