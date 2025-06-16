package queues

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/brienze1/notes-api/internal/integration/entrypoint/dtos"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/domain/exceptions"
	"github.com/brienze1/notes-api/internal/infra/properties"
	"github.com/brienze1/notes-api/internal/integration/queues"
	"github.com/brienze1/notes-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotesQueue_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockSqsClient := mocks.NewMockSqsClient(ctrl)
	noteQueue := queues.NewNotesQueue(mockSqsClient)
	ctx := context.TODO()

	note := entities.Note{
		Title:   "Test Title",
		Content: "Test Content",
	}

	t.Run("successful publish", func(t *testing.T) {

		noteDTO := dtos.Note{}.FromEntity(&note)
		expectedBody, _ := json.Marshal(noteDTO)
		expectedInput := &sqs.SendMessageInput{
			MessageBody: aws.String(string(expectedBody)),
			QueueUrl:    aws.String(properties.GetNotesQueueURL()),
		}

		mockSqsClient.EXPECT().SendMessage(ctx, expectedInput).Return(&sqs.SendMessageOutput{}, nil)

		err := noteQueue.Publish(ctx, note)

		assert.NoError(t, err)
	})

	t.Run("queue error", func(t *testing.T) {

		noteDTO := dtos.Note{}.FromEntity(&note)
		expectedBody, _ := json.Marshal(noteDTO)
		expectedInput := &sqs.SendMessageInput{
			MessageBody: aws.String(string(expectedBody)),
			QueueUrl:    aws.String(properties.GetNotesQueueURL()),
		}

		queueError := errors.New("failed to publish to SQS")
		mockSqsClient.EXPECT().SendMessage(ctx, expectedInput).Return(nil, queueError)

		err := noteQueue.Publish(ctx, note)

		assert.Equal(t, exceptions.NewNotesQueueError(queueError.Error()), err)
	})
}
