package usecases

import (
	"context"
	"github.com/brienze1/notes-api/internal/domain/entities"
)

// TODO implementar a interface
type PublishNoteUseCase interface {
	Publish(ctx context.Context, note entities.Note) (err error)
}
