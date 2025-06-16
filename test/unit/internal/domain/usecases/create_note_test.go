package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/domain/usecases"
	"github.com/brienze1/notes-api/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateNote(t *testing.T) {
	var (
		ctrl             = gomock.NewController(t)
		noteQueueMocked  = mocks.NewMockNoteQueue(ctrl)
		createNoteMocked = mocks.NewMockCreateNoteRepository(ctrl)
		createNote       = usecases.NewCreateNoteUseCase(createNoteMocked, noteQueueMocked)
		ctx              = context.TODO()
		noteEntity       = entities.Note{
			Title:   "Title",
			Content: "Content",
		}
		expectedCreatedNote = entities.Note{
			Id:        uuid.NewString(),
			Title:     "Title",
			Content:   "Content",
			CreatedAt: ptr.Time(time.Now()),
			UpdatedAt: ptr.Time(time.Now()),
		}
	)

	t.Run("should CreateNote with success", func(t *testing.T) {

		createNoteMocked.EXPECT().Create(ctx, noteEntity).Return(&expectedCreatedNote, nil).Times(1)
		noteQueueMocked.EXPECT().Publish(ctx, expectedCreatedNote).Return(nil).Times(1)

		createdNote, err := createNote.Create(ctx, noteEntity)

		assert.Nil(t, err)
		assert.NotNil(t, createdNote)
		assert.Equal(t, expectedCreatedNote, *createdNote)
	})

	t.Run("should return an error when noteQueue returns an error", func(t *testing.T) {

		expectedErr := errors.New(uuid.NewString())
		createNoteMocked.EXPECT().Create(ctx, noteEntity).Return(&expectedCreatedNote, nil).Times(1)
		noteQueueMocked.EXPECT().Publish(ctx, expectedCreatedNote).Return(expectedErr).Times(1)

		_, err := createNote.Create(ctx, noteEntity)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return an error when createNote returns an error", func(t *testing.T) {

		expectedErr := errors.New(uuid.NewString())
		createNoteMocked.EXPECT().Create(ctx, noteEntity).Return(nil, expectedErr).Times(1)
		noteQueueMocked.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)

		_, err := createNote.Create(ctx, noteEntity)

		assert.Equal(t, expectedErr, err)
	})

}
