package boards

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/smarrog/task-board/core-service/internal/model"
	boardrepo "github.com/smarrog/task-board/core-service/internal/repo/boards"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service struct {
	repo boardrepo.Repo
}

func NewSvc(repo boardrepo.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, ownerId uuid.UUID, title, description string) (*model.Board, error) {
	if ownerId == uuid.Nil || title == "" {
		return nil, ErrInvalidArgument
	}

	b := &model.Board{
		OwnerId:     ownerId,
		Title:       title,
		Description: description,
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) Get(ctx context.Context, boardId uuid.UUID) (*model.Board, error) {
	if boardId == uuid.Nil {
		return nil, ErrInvalidArgument
	}

	return s.repo.Get(ctx, boardId)
}

func (s *Service) List(ctx context.Context, ownerId uuid.UUID) ([]model.Board, error) {
	if ownerId == uuid.Nil {
		return nil, ErrInvalidArgument
	}

	return s.repo.List(ctx, ownerId)
}

func (s *Service) Update(ctx context.Context, ownerId, boardId uuid.UUID, title, description string) error {
	if ownerId == uuid.Nil || boardId == uuid.Nil || title == "" {
		return ErrInvalidArgument
	}

	return s.repo.Update(ctx, &model.Board{
		Id:          boardId,
		OwnerId:     ownerId,
		Title:       title,
		Description: description,
	})
}

func (s *Service) Delete(ctx context.Context, boardId uuid.UUID) error {
	if boardId == uuid.Nil {
		return ErrInvalidArgument
	}
	return s.repo.Delete(ctx, boardId)
}
