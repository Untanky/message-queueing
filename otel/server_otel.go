package otel

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	queueing "message-queueing"
	"time"
)

type otelQueueServiceServer struct {
	numberOfRequests metric.Int64Counter
	numberOfErrors   metric.Int64Counter
	latency          metric.Int64Histogram

	server queueing.QueueServiceServer

	queueing.UnimplementedQueueServiceServer
}

func WrapQueueServiceServer(server queueing.QueueServiceServer) (queueing.QueueServiceServer, error) {
	numberOfRequests, err := meter.Int64Counter("numberOfRequests")
	if err != nil {
		return nil, err
	}

	numberOfErrors, err := meter.Int64Counter("numberOfErrors")
	if err != nil {
		return nil, err
	}

	latency, err := meter.Int64Histogram("requestLatency", metric.WithUnit("μs"))
	if err != nil {
		return nil, err
	}

	return otelQueueServiceServer{
		numberOfRequests: numberOfRequests,
		numberOfErrors:   numberOfErrors,
		latency:          latency,

		server: server,
	}, nil
}

func (o otelQueueServiceServer) SubmitMessages(
	ctx context.Context, request *queueing.SubmitMessagesRequest,
) (*queueing.SubmitMessagesResponse, error) {
	start := time.Now()
	ctx, span := tracer.Start(ctx, "SubmitMessages")
	defer span.End()

	o.numberOfRequests.Add(ctx, 1)

	response, err := o.server.SubmitMessages(ctx, request)
	if err != nil {
		o.numberOfErrors.Add(ctx, 1)
		span.RecordError(err)
	}

	end := time.Now()
	o.latency.Record(ctx, end.Sub(start).Microseconds())

	return response, err
}

func (o otelQueueServiceServer) RetrieveMessages(
	ctx context.Context, request *queueing.RetrieveMessagesRequest,
) (*queueing.RetrieveMessagesResponse, error) {
	start := time.Now()
	ctx, span := tracer.Start(ctx, "RetrieveMessages")
	defer span.End()

	o.numberOfRequests.Add(ctx, 1)

	response, err := o.server.RetrieveMessages(ctx, request)
	if err != nil {
		o.numberOfErrors.Add(ctx, 1)
		span.RecordError(err)
	}

	end := time.Now()
	o.latency.Record(ctx, end.Sub(start).Microseconds())

	return response, err
}

func (o otelQueueServiceServer) AcknowledgeMessages(
	ctx context.Context, request *queueing.AcknowledgeMessagesRequest,
) (*queueing.AcknowledgeMessagesResponse, error) {
	start := time.Now()
	ctx, span := tracer.Start(ctx, "AcknowledgeMessages")
	defer span.End()

	o.numberOfRequests.Add(ctx, 1)

	response, err := o.server.AcknowledgeMessages(ctx, request)
	if err != nil {
		o.numberOfErrors.Add(ctx, 1)
		span.RecordError(err)
	}

	end := time.Now()
	o.latency.Record(ctx, end.Sub(start).Microseconds())

	return response, err
}

func (o otelQueueServiceServer) mustEmbedUnimplementedQueueServiceServer() {
}
