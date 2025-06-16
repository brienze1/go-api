package usecases

import (
	"context"

	"github.com/brienze1/notes-api/internal/domain/adapters"
	"github.com/brienze1/notes-api/internal/domain/entities"
)

type CreateNoteUseCase struct {
	createNoteRepository adapters.CreateNoteRepository
	noteQueue            adapters.NoteQueue
}

func (uc CreateNoteUseCase) Create(ctx context.Context, note entities.Note) (createdNote *entities.Note, err error) {
	createdNote, err = uc.createNoteRepository.Create(ctx, note)
	if err != nil {
		return
	}

	err = uc.noteQueue.Publish(ctx, *createdNote)
	return
}

func NewCreateNoteUseCase(createNoteRepository adapters.CreateNoteRepository, notesQueue adapters.NoteQueue) *CreateNoteUseCase {
	return &CreateNoteUseCase{
		createNoteRepository: createNoteRepository,
		noteQueue:            notesQueue,
	}
}
