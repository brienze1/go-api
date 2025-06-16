package adapters

import (
	"context"

	"github.com/brienze1/notes-api/internal/domain/entities"
)

type NoteQueue interface {
	Publish(ctx context.Context, note entities.Note) (err error)
}
