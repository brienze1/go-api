package queues

import (
	"context"
	"encoding/json"
	"github.com/brienze1/notes-api/internal/integration/entrypoint/dtos"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/brienze1/notes-api/internal/domain/adapters"
	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/domain/exceptions"
	"github.com/brienze1/notes-api/internal/infra/properties"
)

type SqsClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type notes struct {
	sqsClient SqsClient
}

func (n *notes) Publish(ctx context.Context, note entities.Note) (err error) {
	var noteDTO dtos.Note

	bytes, err := json.Marshal(noteDTO.FromEntity(&note))
	if err != nil {
		return exceptions.NewNotesQueueError(err.Error())
	}

	if _, err = n.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(bytes)),
		QueueUrl:    aws.String(properties.GetNotesQueueURL()),
	}); err != nil {
		return exceptions.NewNotesQueueError(err.Error())
	}

	return nil
}

func NewNotesQueue(sqsClient SqsClient) adapters.NoteQueue {
	return &notes{sqsClient}
}
