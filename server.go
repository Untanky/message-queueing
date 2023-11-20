package queueing

import (
	"context"
	"crypto/sha256"
	"github.com/google/uuid"
	"time"
)

type WriteAheadLog interface {
	Write(ctx context.Context, message *QueueMessage) error
}

type MessageQueueingServer struct {
	UnimplementedQueueServiceServer

	queueService Queue[*QueueMessage]
}

func NewMessageQueueingServer(service Queue[*QueueMessage]) QueueServiceServer {
	return &MessageQueueingServer{
		queueService: service,
	}
}

func (m *MessageQueueingServer) SubmitMessages(
	ctx context.Context, request *SubmitMessagesRequest,
) (*SubmitMessagesResponse, error) {
	messages := make([]*QueueMessage, 0, len(request.Messages))

	for _, rawQueueMessage := range request.Messages {
		messages = append(messages, rawQueueMessage.ToQueueMessage())
	}

	err := m.queueService.Enqueue(ctx, messages...)
	if err != nil {
		return nil, err
	}

	return &SubmitMessagesResponse{
		Receipts: []*SubmitReceipt{},
	}, nil
}

func (m *RawQueueMessage) ToQueueMessage() *QueueMessage {
	messageID := uuid.NewString()
	now := time.Now().Unix()
	hash := sha256.New()
	hash.Write(m.RawMessage)

	return &QueueMessage{
		MessageID:  &messageID,
		Timestamp:  &now,
		Attributes: m.Attributes,
		Data:       m.RawMessage,
		DataHash:   hash.Sum(nil),
	}
}

func (m *MessageQueueingServer) RetrieveMessages(
	ctx context.Context, request *RetrieveMessagesRequest,
) (*RetrieveMessagesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageQueueingServer) AcknowledgeMessages(
	ctx context.Context, request *AcknowledgeMessagesRequest,
) (*AcknowledgeMessagesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageQueueingServer) mustEmbedUnimplementedQueueServiceServer() {}
