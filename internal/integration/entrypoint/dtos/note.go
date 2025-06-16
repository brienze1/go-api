package dtos

import (
	"time"

	"github.com/brienze1/notes-api/internal/domain/entities"
)

// TODO arrumar entirades (json com dto no delivery) (gorm com model no integration)
type Note struct {
	Id        string     `json:"id"`
	Title     string     `json:"title" validate:"required,min=3"`
	Content   string     `json:"content" validate:"required,min=3"`
	UserId    string     `json:"user_id"  validate:"required"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (m Note) ToEntity() entities.Note {
	return entities.Note{
		Id:        m.Id,
		Title:     m.Title,
		Content:   m.Content,
		UserId:    m.UserId,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m Note) FromEntity(note *entities.Note) Note {
	return Note{
		Id:        note.Id,
		Title:     note.Title,
		Content:   note.Content,
		UserId:    note.UserId,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

type Notes []*Note

func (u Notes) FromEntities(notes []*entities.Note) Notes {
	u = make([]*Note, len(notes))
	for i, note := range notes {
		u[i] = &Note{
			Id:        note.Id,
			Title:     note.Title,
			Content:   note.Content,
			UserId:    note.UserId,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		}
	}
	return u
}
