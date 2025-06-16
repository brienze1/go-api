package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/brienze1/notes-api/internal/integration/persistence/repositories/models"

	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/domain/exceptions"
	"github.com/brienze1/notes-api/internal/integration/adapters"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type noteRepository struct {
	connection *gorm.DB
}

func (n *noteRepository) Get(ctx context.Context, userId string, noteId string) (*entities.Note, error) {
	var note *models.Note
	err := n.connection.WithContext(ctx).Where(
		&models.Note{
			UserId: userId,
			Id:     noteId,
		}).First(&note).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(fmt.Sprintf("note with id %s not found", noteId))
		}
	}
	return note.ToEntity(), err
}

func (n *noteRepository) Create(ctx context.Context, note entities.Note) (*entities.Note, error) {
	var noteModel models.Note
	err := n.connection.WithContext(ctx).Create(noteModel.FromEntity(note)).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, exceptions.NewNotFoundError(fmt.Sprintf("user with id %s not found", note.UserId))
		}
		return nil, exceptions.NewInternalServerError("failed to create note", err.Error())
	}
	return noteModel.ToEntity(), nil
}

func (n *noteRepository) Delete(ctx context.Context, userId, noteId string) error {
	tx := n.connection.WithContext(ctx).Delete(&models.Note{Id: noteId, UserId: userId})
	err := tx.Error
	if err != nil {
		return exceptions.NewInternalServerError(fmt.Sprintf("failed to delete note with id %s and userId %s", noteId, userId), err.Error())
	}

	if tx.RowsAffected == 0 {
		return exceptions.NewNotFoundError(fmt.Sprintf("note with id %s and user %s not found", noteId, userId))
	}
	return nil
}

func (n *noteRepository) List(ctx context.Context, userId string) ([]*entities.Note, error) {
	var notes models.Notes
	err := n.connection.WithContext(ctx).Where(&models.Note{UserId: userId}).Find(&notes).Error
	if err != nil {
		return nil, exceptions.NewInternalServerError(fmt.Sprintf("failed to list notes from userId %s", userId), err.Error())
	}
	return notes.ToEntities(), err
}

func NewNote(connection *gorm.DB) adapters.NoteRepository {
	return &noteRepository{connection: connection}
}
