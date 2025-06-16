package queues

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/integration/queues"
	"github.com/brienze1/notes-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func FuzzNotesQueue_Publish(f *testing.F) {

	f.Add("Valid Title", "Valid Content")
	f.Add("", "")
	f.Add("A very long title with multiple characters and symbols!@#$%^", "Another long content string with special chars *()<>?")

	ctx := context.TODO()

	f.Fuzz(func(t *testing.T, title, content string) {
		ctrl := gomock.NewController(t)
		mockSqsClient := mocks.NewMockSqsClient(ctrl)
		noteQueue := queues.NewNotesQueue(mockSqsClient)

		note := entities.Note{
			Title:   title,
			Content: content,
		}

		mockSqsClient.EXPECT().SendMessage(ctx, gomock.Any()).Return(&sqs.SendMessageOutput{}, nil)

		err := noteQueue.Publish(ctx, note)
		assert.NoError(t, err)
	})
}
