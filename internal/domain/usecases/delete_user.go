package usecases

import (
	"context"
)

type DeleteUserUseCase interface {
	Delete(ctx context.Context, id string) (err error)
}
