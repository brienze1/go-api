package usecases

import (
	"context"
	"github.com/brienze1/notes-api/internal/domain/entities"
)

type ListNotesUseCase interface {
	List(ctx context.Context, userId string) (notes []*entities.Note, err error)
}
