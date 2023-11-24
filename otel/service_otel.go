package otel

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	queueing "message-queueing"
)

type otelPriorityQueue struct {
	messagesWritten      metric.Int64Counter
	messagesRetrieved    metric.Int64Counter
	messagesAcknowledged metric.Int64Counter
	tracer               trace.Tracer

	service queueing.Service
}

func WrapOTelPriorityQueue(service queueing.Service) (
	queueing.Service, error,
) {
	meter := otel.Meter("message-queueing")

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

		tracer: otel.Tracer("message-queueing"),
	}, nil
}

func (o *otelPriorityQueue) Enqueue(ctx context.Context, messages ...*queueing.QueueMessage) error {
	ctx, span := o.tracer.Start(ctx, "Enqueue")
	defer span.End()

	o.messagesWritten.Add(ctx, int64(len(messages)))

	return o.service.Enqueue(ctx, messages...)
}

func (o *otelPriorityQueue) Dequeue(ctx context.Context, messages []*queueing.QueueMessage) (int, error) {
	ctx, span := o.tracer.Start(ctx, "Dequeue")
	defer span.End()

	n, err := o.service.Dequeue(ctx, messages)

	return n, err
}

func (o *otelPriorityQueue) Retrieve(ctx context.Context, messages []*queueing.QueueMessage) (int, error) {
	ctx, span := o.tracer.Start(ctx, "Retrieve")
	defer span.End()

	n, err := o.service.Retrieve(ctx, messages)

	o.messagesRetrieved.Add(ctx, int64(n))

	return n, err
}

func (o *otelPriorityQueue) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	ctx, span := o.tracer.Start(ctx, "Acknowledge")
	defer span.End()

	err := o.service.Acknowledge(ctx, messageID)

	o.messagesAcknowledged.Add(ctx, 1)

	return err
}
