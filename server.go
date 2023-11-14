package queueing

import "context"

type WriteAheadLog interface {
	Write([]byte) (int, error)
}

type MessageQueueingServer struct {
	UnimplementedQueueServiceServer
}

func (m MessageQueueingServer) SubmitMessages(
	ctx context.Context, request *SubmitMessagesRequest,
) (*SubmitMessagesResponse, error) {
	//TODO implement me
	panic("implement me")
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
