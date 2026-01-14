package usecase

import (
	"context"

	"github.com/smarrog/task-board/auth-service/internal/config"
	do "github.com/smarrog/task-board/auth-service/internal/domain"
)

type RegisterUseCase struct {
	repo do.Repository
	cfg  *config.Config
}

type RegisterInput struct {
	Email    string
	Username string
	Pwd      string
}

type RegisterOutput struct {
	User        *do.User
	AccessToken do.AccessToken
}

func NewRegisterUseCase(repo do.Repository, cfg *config.Config) *RegisterUseCase {
	return &RegisterUseCase{repo: repo, cfg: cfg}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	email, err := do.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}
	userName, err := do.NewUserName(input.Username)
	if err != nil {
		return nil, err
	}
	pwd, err := do.NewPwd(input.Pwd)
	if err != nil {
		return nil, err
	}
	pwdHash, err := getPwdHash(pwd)
	if err != nil {
		return nil, err
	}

	u := do.NewUser(email, userName, pwdHash)
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	token, err := getAccessToken(u.Id().String(), uc.cfg.JWTSecret, uc.cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{User: u, AccessToken: token}, nil
}
