package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/smarrog/task-board/auth-service/internal/config"
	do "github.com/smarrog/task-board/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	repo do.Repository
	cfg  *config.Config
}

type LoginInput struct {
	Email string
	Pwd   string
}

type LoginOutput struct {
	User        *do.User
	AccessToken do.AccessToken
}

func NewLoginUseCase(repo do.Repository, cfg *config.Config) *LoginUseCase {
	return &LoginUseCase{repo: repo, cfg: cfg}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	email, err := do.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	u, err := uc.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", do.ErrInvalidCredentials, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PwdHash().String()), []byte(input.Pwd)); err != nil {
		return nil, fmt.Errorf("%w: %v", do.ErrInvalidCredentials, err)
	}

	token, err := getAccessToken(u.Id().String(), uc.cfg.JWTSecret, uc.cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", do.ErrInvalidCredentials, err)
	}

	return &LoginOutput{User: u, AccessToken: token}, nil
}
