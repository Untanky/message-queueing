package otel

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	queueing "message-queueing"
)

type otelQueueServiceServer struct {
	numberOfRequests metric.Int64Counter
	numberOfError    metric.Int64Counter

	server queueing.QueueServiceServer

	queueing.UnimplementedQueueServiceServer
}

func WrapQueueServiceServer(server queueing.QueueServiceServer) (queueing.QueueServiceServer, error) {
	numberOfRequests, err := meter.Int64Counter("numberOfRequests")
	if err != nil {
		return nil, err
	}

	return otelQueueServiceServer{
		numberOfRequests: numberOfRequests,
		server:           server,
	}, nil
}

func (o otelQueueServiceServer) SubmitMessages(
	ctx context.Context, request *queueing.SubmitMessagesRequest,
) (*queueing.SubmitMessagesResponse, error) {
	ctx, span := tracer.Start(ctx, "SubmitMessages")
	defer span.End()

	o.numberOfRequests.Add(ctx, 1)

	response, err := o.server.SubmitMessages(ctx, request)
	if err != nil {
		o.numberOfError.Add(ctx, 1)
		span.RecordError(err)
	}

	return response, err
}

func (o otelQueueServiceServer) RetrieveMessages(
	ctx context.Context, request *queueing.RetrieveMessagesRequest,
) (*queueing.RetrieveMessagesResponse, error) {
	ctx, span := tracer.Start(ctx, "RetrieveMessages")
	defer span.End()

	o.numberOfRequests.Add(ctx, 1)

	response, err := o.server.RetrieveMessages(ctx, request)
	if err != nil {
		o.numberOfError.Add(ctx, 1)
		span.RecordError(err)
	}

	return response, err
}

func (o otelQueueServiceServer) AcknowledgeMessages(
	ctx context.Context, request *queueing.AcknowledgeMessagesRequest,
) (*queueing.AcknowledgeMessagesResponse, error) {
	ctx, span := tracer.Start(ctx, "AcknowledgeMessages")
	defer span.End()

	o.numberOfRequests.Add(ctx, 1)

	response, err := o.server.AcknowledgeMessages(ctx, request)
	if err != nil {
		o.numberOfError.Add(ctx, 1)
		span.RecordError(err)
	}

	return response, err
}

func (o otelQueueServiceServer) mustEmbedUnimplementedQueueServiceServer() {
}
