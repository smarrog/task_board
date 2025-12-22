package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type ListBoardsUseCase struct {
	repo board.Repository
}

func NewListBoardsUseCase(repo board.Repository) *ListBoardsUseCase {
	return &ListBoardsUseCase{repo: repo}
}

func (uc *ListBoardsUseCase) Execute(ctx context.Context, ownerId string) ([]*board.Board, error) {
	oid, err := board.UserIdFromString(ownerId)
	if err != nil {
		return nil, fmt.Errorf("owner_id: %w", err)
	}

	return uc.repo.ListByOwner(ctx, oid)
}
