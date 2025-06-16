package adapters

import (
	"context"
)

type NotesController interface {
	ListNotes(ctx context.Context, userId string) (response []byte, status int)
	GetNote(ctx context.Context, userId, noteId string) (response []byte, status int)
	CreateNote(ctx context.Context, userId string, body []byte) (response []byte, status int)
	DeleteNote(ctx context.Context, noteId, userId string) (response []byte, status int)
}
