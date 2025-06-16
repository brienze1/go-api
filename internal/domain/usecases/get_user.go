package usecases

import (
	"context"
	"github.com/brienze1/notes-api/internal/domain/entities"
)

type GetUserUseCase interface {
	Get(ctx context.Context, id string) (user *entities.User, err error)
}
