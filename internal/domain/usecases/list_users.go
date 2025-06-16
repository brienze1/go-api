package usecases

import (
	"context"
	"github.com/brienze1/notes-api/internal/domain/entities"
)

type ListUsersUseCase interface {
	List(ctx context.Context) (users []*entities.User, err error)
}
