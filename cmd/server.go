package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	queueing "message-queueing"
)

type WriteAheadLog interface {
	Write(ctx context.Context, message *queueing.QueueMessage) error
}

type MessageQueueingServer struct {
	queueing.UnimplementedQueueServiceServer

	queueService queueing.Service
}

func NewMessageQueueingServer(service queueing.Service) queueing.QueueServiceServer {
	return &MessageQueueingServer{
		queueService: service,
	}
}

func (m *MessageQueueingServer) WriteMessages(stream queueing.QueueService_WriteMessagesServer) error {
	for {
		rawQueueMessage, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		message := rawQueueMessage.ToQueueMessage()
		err = m.queueService.Enqueue(stream.Context(), message)
		if err != nil {
			return err
		}

		messageID := uuid.UUID(message.MessageID).String()
		var ok = err == nil
		var reason string
		if !ok {
			reason = err.Error()
		}

		err = stream.Send(
			&queueing.SubmitReceipt{
				MessageID: &messageID,
				Ok:        &ok,
				Reason:    &reason,
				Timestamp: message.Timestamp,
				DataHash:  message.DataHash,
			},
		)
		if err != nil {
			return err
		}
	}
}

func (m *MessageQueueingServer) SubmitMessages(
	ctx context.Context, request *queueing.SubmitMessagesRequest,
) (*queueing.SubmitMessagesResponse, error) {
	receipts := make([]*queueing.SubmitReceipt, 0, len(request.Messages))
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
			receipts, &queueing.SubmitReceipt{
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

	return &queueing.SubmitMessagesResponse{
		Receipts: receipts,
	}, joinErr
}

func (m *MessageQueueingServer) RetrieveMessages(
	ctx context.Context, request *queueing.RetrieveMessagesRequest,
) (*queueing.RetrieveMessagesResponse, error) {
	var messages = make([]*queueing.QueueMessage, *request.Count)
	n, err := m.queueService.Retrieve(ctx, messages)
	if err != nil && !errors.Is(err, queueing.NextMessageNotReady) {
		return nil, err
	}

	n32 := int32(n)

	return &queueing.RetrieveMessagesResponse{
		Count:    &n32,
		Messages: messages[:n],
	}, nil
}

func (m *MessageQueueingServer) AcknowledgeMessages(
	ctx context.Context, request *queueing.AcknowledgeMessagesRequest,
) (*queueing.AcknowledgeMessagesResponse, error) {
	receipts := make([]*queueing.AcknowledgeReceipt, 0, len(request.MessageIDs))

	for _, idString := range request.MessageIDs {
		id := uuid.MustParse(idString)

		e := m.queueService.Acknowledge(ctx, id)

		var ok = e == nil
		var reason string
		if !ok {
			reason = e.Error()
		}

		receipts = append(
			receipts, &queueing.AcknowledgeReceipt{
				MessageID: &idString,
				Ok:        &ok,
				Reason:    &reason,
			},
		)
	}

	return &queueing.AcknowledgeMessagesResponse{
		Receipts: receipts,
	}, nil
}

func (m *MessageQueueingServer) mustEmbedUnimplementedQueueServiceServer() {}
