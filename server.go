package queueing

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type WriteAheadLog interface {
	Write(ctx context.Context, message *QueueMessage) error
}

type MessageQueueingServer struct {
	UnimplementedQueueServiceServer

	queueService Service
}

func NewMessageQueueingServer(service Service) QueueServiceServer {
	return &MessageQueueingServer{
		queueService: service,
	}
}

func (m *MessageQueueingServer) SubmitMessages(
	ctx context.Context, request *SubmitMessagesRequest,
) (*SubmitMessagesResponse, error) {
	receipts := make([]*SubmitReceipt, 0, len(request.Messages))
	var joinErr error

	for _, rawQueueMessage := range request.Messages {
		message := rawQueueMessage.ToQueueMessage()
		err := m.queueService.Enqueue(ctx, message)

		joinErr = errors.Join(joinErr, err)

		messageID := uuid.UUID(message.MessageID).String()
		var ok = err == nil
		var reason string
		if !ok {
			reason = err.Error()
		}

		receipts = append(
			receipts, &SubmitReceipt{
				MessageID: &messageID,
				Ok:        &ok,
				Reason:    &reason,
				Timestamp: message.Timestamp,
				DataHash:  message.DataHash,
			},
		)
	}

	if joinErr != nil {
		joinErr = fmt.Errorf("there were some errors: %w", joinErr)
	}

	return &SubmitMessagesResponse{
		Receipts: receipts,
	}, joinErr
}

func (m *RawQueueMessage) ToQueueMessage() *QueueMessage {
	messageID := uuid.New()
	now := time.Now().Unix()
	hash := sha256.New()
	hash.Write(m.Data)

	return &QueueMessage{
		MessageID:  messageID[:],
		Timestamp:  &now,
		Attributes: m.Attributes,
		Data:       m.Data,
		DataHash:   hash.Sum(nil),
	}
}

func (m *MessageQueueingServer) RetrieveMessages(
	ctx context.Context, request *RetrieveMessagesRequest,
) (*RetrieveMessagesResponse, error) {
	var messages = make([]*QueueMessage, *request.Count)
	n, err := m.queueService.Retrieve(ctx, messages)
	if err != nil && !errors.Is(err, NextMessageNotReady) {
		return nil, err
	}

	n32 := int32(n)

	return &RetrieveMessagesResponse{
		Count:    &n32,
		Messages: messages[:n],
	}, nil
}

func (m *MessageQueueingServer) AcknowledgeMessages(
	ctx context.Context, request *AcknowledgeMessagesRequest,
) (*AcknowledgeMessagesResponse, error) {
	receipts := make([]*AcknowledgeReceipt, 0, len(request.MessageIDs))

	for _, idString := range request.MessageIDs {
		id := uuid.MustParse(idString)

		e := m.queueService.Acknowledge(ctx, id)

		var ok = e == nil
		var reason string
		if !ok {
			reason = e.Error()
		}

		receipts = append(
			receipts, &AcknowledgeReceipt{
				MessageID: &idString,
				Ok:        &ok,
				Reason:    &reason,
			},
		)
	}

	return &AcknowledgeMessagesResponse{
		Receipts: receipts,
	}, nil
}

func (m *MessageQueueingServer) mustEmbedUnimplementedQueueServiceServer() {}
