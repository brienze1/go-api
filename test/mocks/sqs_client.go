package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/google/uuid"
)

type SqsClient struct {
	messages map[string][]*sqs.SendMessageInput
}

var sqsMockInstance *SqsClient

func Sqs() *SqsClient {
	if sqsMockInstance == nil {
		sqsMockInstance = &SqsClient{
			messages: make(map[string][]*sqs.SendMessageInput),
		}
	}

	return sqsMockInstance
}

func (s *SqsClient) SendMessage(_ context.Context, input *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	messageId := uuid.New().String()

	if s.messages[*input.QueueUrl] == nil {
		s.messages[*input.QueueUrl] = make([]*sqs.SendMessageInput, 0)
	}
	s.messages[*input.QueueUrl] = append(s.messages[*input.QueueUrl], input)

	return &sqs.SendMessageOutput{
		MessageId:      &messageId,
		SequenceNumber: nil,
		ResultMetadata: middleware.Metadata{},
	}, nil
}

func (s *SqsClient) GetMessages(queue string) []*string {
	var messages []*string
	for _, message := range s.messages[queue] {
		messages = append(messages, message.MessageBody)
	}

	return messages
}

func (s *SqsClient) GetMessageAttributes(queue string) []map[string]types.MessageAttributeValue {
	var messageAttributes []map[string]types.MessageAttributeValue
	for _, message := range s.messages[queue] {
		messageAttributes = append(messageAttributes, message.MessageAttributes)
	}

	return messageAttributes
}

func (s *SqsClient) GetMessageGroupId(queue string) []*string {
	var messageGroupIds []*string
	for _, message := range s.messages[queue] {
		messageGroupIds = append(messageGroupIds, message.MessageGroupId)
	}

	return messageGroupIds
}

func (s *SqsClient) GetMessageDeduplicationId(queue string) []*string {
	var messageDeduplicationIds []*string
	for _, message := range s.messages[queue] {
		messageDeduplicationIds = append(messageDeduplicationIds, message.MessageGroupId)
	}

	return messageDeduplicationIds
}

func (s *SqsClient) Reset() *SqsClient {
	s.messages = make(map[string][]*sqs.SendMessageInput)
	return s
}
