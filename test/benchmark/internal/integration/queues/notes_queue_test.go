package queues

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go/ptr"
	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/integration/queues"
	"github.com/brienze1/notes-api/test/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func randomString(length int) string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func BenchmarkNotesQueue_Publish(b *testing.B) {
	var (
		ctx           = context.TODO()
		ctrl          = gomock.NewController(b)
		mockSqsClient = mocks.NewMockSqsClient(ctrl)
		noteQueue     = queues.NewNotesQueue(mockSqsClient)
		validNote     = entities.Note{
			Id:        uuid.NewString(),
			Title:     randomString(100000000),
			Content:   randomString(100000000),
			CreatedAt: ptr.Time(time.Now()),
			UpdatedAt: ptr.Time(time.Now()),
		}
	)

	mockSqsClient.EXPECT().
		SendMessage(ctx, gomock.Any()).Return(&sqs.SendMessageOutput{}, nil).AnyTimes()

	b.Run("Publish valid note", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := noteQueue.Publish(ctx, validNote)
			if err != nil {
				b.Fatalf("failed to publish valid note: %v", err)
			}
		}
	})
}
