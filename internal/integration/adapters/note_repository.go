package adapters

import (
	"context"

	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/domain/usecases"
)

type NoteRepository interface {
	usecases.GetNoteUseCase
	usecases.ListNotesUseCase
	usecases.DeleteNoteUseCase

	Create(ctx context.Context, note entities.Note) (createdNote *entities.Note, err error)
}
