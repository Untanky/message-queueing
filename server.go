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

	wal WriteAheadLog
}

func (m MessageQueueingServer) SubmitMessages(
	ctx context.Context, request *SubmitMessagesRequest,
) (*SubmitMessagesResponse, error) {
	receipts := make([]*SubmitReceipt, len(request.Messages))

	for _, rawQueueMessage := range request.Messages {
		queueMessage := rawQueueMessage.ToQueueMessage()
		err := m.wal.Write(ctx, queueMessage)
		if err != nil {
			f := false
			reason := err.Error()

			receipts = append(
				receipts, &SubmitReceipt{
					Ok:        &f,
					MessageID: queueMessage.MessageID,
					Reason:    &reason,
				},
			)
		} else {
			t := true

			receipts = append(
				receipts, &SubmitReceipt{
					Ok:        &t,
					MessageID: queueMessage.MessageID,
					Timestamp: queueMessage.Timestamp,
					DataHash:  queueMessage.DataHash,
				},
			)
		}
	}

	return &SubmitMessagesResponse{
		Receipts: receipts,
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

func (m MessageQueueingServer) RetrieveMessages(
	ctx context.Context, request *RetrieveMessagesRequest,
) (*RetrieveMessagesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MessageQueueingServer) AcknowledgeMessages(
	ctx context.Context, request *AcknowledgeMessagesRequest,
) (*AcknowledgeMessagesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MessageQueueingServer) mustEmbedUnimplementedQueueServiceServer() {}
