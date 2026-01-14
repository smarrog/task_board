package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	do "github.com/smarrog/task-board/auth-service/internal/domain"
	uc "github.com/smarrog/task-board/auth-service/internal/usecase"
	v1 "github.com/smarrog/task-board/shared/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	v1.UnimplementedAuthServiceServer
	log      *zerolog.Logger
	register *uc.RegisterUseCase
	login    *uc.LoginUseCase
}

func NewAuthHandler(log *zerolog.Logger, register *uc.RegisterUseCase, login *uc.LoginUseCase) *AuthHandler {
	return &AuthHandler{log: log, register: register, login: login}
}

func (h *AuthHandler) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	out, err := h.register.Execute(ctx, uc.RegisterInput{Email: req.GetEmail(), Username: req.GetUsername(), Pwd: req.GetPassword()})
	if err != nil {
		switch {
		case errors.Is(err, do.ErrEmailAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "email_exists")
		default:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return &v1.RegisterResponse{
		User: &v1.User{
			Id:       out.User.Id().String(),
			Email:    out.User.Email().String(),
			Username: out.User.Username().String(),
		},
		AccessToken: out.AccessToken.String(),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	out, err := h.login.Execute(ctx, uc.LoginInput{Email: req.GetEmail(), Pwd: req.GetPassword()})
	if err != nil {
		switch {
		case errors.Is(err, do.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid_credentials")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &v1.LoginResponse{
		User: &v1.User{
			Id:       out.User.Id().String(),
			Email:    out.User.Email().String(),
			Username: out.User.Username().String(),
		},
		AccessToken: out.AccessToken.String(),
	}, nil
}
