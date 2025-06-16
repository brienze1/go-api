package usecases

import (
	"context"
)

type DeleteNoteUseCase interface {
	Delete(ctx context.Context, userId, noteId string) (err error)
}
