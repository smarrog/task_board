package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/api-service/internal/config"
	authv1 "github.com/smarrog/task-board/shared/proto/auth/v1"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	log  *zerolog.Logger
	cfg  *config.Config
	auth authv1.AuthServiceClient
}

func NewAuthHandler(log *zerolog.Logger, cfg *config.Config, conn *grpc.ClientConn) *AuthHandler {
	return &AuthHandler{log: log, cfg: cfg, auth: authv1.NewAuthServiceClient(conn)}
}

type registerBody struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body registerBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}
	ctx, cancel := h.reqCtxFromCfg()
	defer cancel()

	resp, err := h.auth.Register(ctx, &authv1.RegisterRequest{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		return grpcToHTTP(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": fiber.Map{
			"id":       resp.GetUser().GetId(),
			"email":    resp.GetUser().GetEmail(),
			"username": resp.GetUser().GetUsername(),
		},
		"access_token": resp.GetAccessToken(),
	})
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body loginBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid_json")
	}
	ctx, cancel := h.reqCtxFromCfg()
	defer cancel()

	resp, err := h.auth.Login(ctx, &authv1.LoginRequest{Email: body.Email, Password: body.Password})
	if err != nil {
		return grpcToHTTP(err)
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":       resp.GetUser().GetId(),
			"email":    resp.GetUser().GetEmail(),
			"username": resp.GetUser().GetUsername(),
		},
		"access_token": resp.GetAccessToken(),
	})
}

func (h *AuthHandler) reqCtxFromCfg() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), h.cfg.RequestTimeout)
}
