package otel

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	queueing "message-queueing"
)

var (
	tracer = otel.Tracer("message-queueing")
	meter  = otel.Meter("message-queueing")
)

type otelPriorityQueue struct {
	messagesWritten      metric.Int64Counter
	messagesRetrieved    metric.Int64Counter
	messagesAcknowledged metric.Int64Counter

	service queueing.Service
}

func WrapService(service queueing.Service) (
	queueing.Service, error,
) {
	messagesWritten, err := meter.Int64Counter("messagesWritten")
	if err != nil {
		return nil, err
	}

	messagesRetrieved, err := meter.Int64Counter("messageRetrieved")
	if err != nil {
		return nil, err
	}

	messagesAcknowledged, err := meter.Int64Counter("messagesAcknowledged")
	if err != nil {
		return nil, err
	}

	return &otelPriorityQueue{
		service: service,

		messagesWritten:      messagesWritten,
		messagesRetrieved:    messagesRetrieved,
		messagesAcknowledged: messagesAcknowledged,
	}, nil
}

func (o *otelPriorityQueue) Enqueue(ctx context.Context, message *queueing.QueueMessage) error {
	ctx, span := tracer.Start(ctx, "Enqueue")
	defer span.End()

	o.messagesWritten.Add(ctx, 1)

	err := o.service.Enqueue(ctx, message)

	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (o *otelPriorityQueue) Retrieve(ctx context.Context, messages []*queueing.QueueMessage) (int, error) {
	ctx, span := tracer.Start(ctx, "Retrieve")
	defer span.End()

	n, err := o.service.Retrieve(ctx, messages)

	o.messagesRetrieved.Add(ctx, int64(n))
	if err != nil {
		span.RecordError(err)
	}

	return n, err
}

func (o *otelPriorityQueue) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "Acknowledge")
	defer span.End()

	err := o.service.Acknowledge(ctx, messageID)

	o.messagesAcknowledged.Add(ctx, 1)
	if err != nil {
		span.RecordError(err)
	}

	return err
}
