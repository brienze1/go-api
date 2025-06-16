package usecases

import (
	"context"
	"github.com/brienze1/notes-api/internal/domain/entities"
)

type GetNoteUseCase interface {
	Get(ctx context.Context, userId, noteId string) (note *entities.Note, err error)
}
