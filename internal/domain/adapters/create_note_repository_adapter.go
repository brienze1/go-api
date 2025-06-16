package adapters

import (
	"context"

	"github.com/brienze1/notes-api/internal/domain/entities"
)

type CreateNoteRepository interface {
	Create(ctx context.Context, note entities.Note) (createdNote *entities.Note, err error)
}
