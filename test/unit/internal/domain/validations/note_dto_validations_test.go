package main

import (
	"github.com/brienze1/notes-api/internal/integration/entrypoint/dtos"
	"github.com/brienze1/notes-api/internal/integration/entrypoint/validations"
	"testing"

	"github.com/brienze1/notes-api/internal/domain/exceptions"
	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {

	tests := []struct {
		name    string
		note    dtos.Note
		wantErr bool
		err     error
	}{
		{
			name: "valid note",
			note: dtos.Note{
				Id:      "1",
				Title:   "Valid Title",
				Content: "Valid Content",
				UserId:  "user123",
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "missing title",
			note: dtos.Note{
				Id:      "2",
				Content: "Valid Content",
				UserId:  "user123",
			},
			wantErr: true,
			err:     exceptions.NewValidationError("'title' is required"),
		},
		{
			name: "short title",
			note: dtos.Note{
				Id:      "3",
				Title:   "Hi",
				Content: "Valid Content",
				UserId:  "user123",
			},
			wantErr: true,
			err:     exceptions.NewValidationError("'title' should be greater in length"),
		},
		{
			name: "missing content",
			note: dtos.Note{
				Id:     "4",
				Title:  "Valid Title",
				UserId: "user123",
			},
			wantErr: true,
			err:     exceptions.NewValidationError("'content' is required"),
		},
		{
			name: "missing user_id",
			note: dtos.Note{
				Id:      "5",
				Title:   "Valid Title",
				Content: "Valid Content",
			},
			wantErr: true,
			err:     exceptions.NewValidationError("'user_id' is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validations.Validate(&tt.note)

			if tt.wantErr {
				assert.Equal(t, err, tt.err)
			}

		})
	}
}
