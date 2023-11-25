package otel

import (
	"context"
	"github.com/google/uuid"
	queueing "message-queueing"
)

type otelRepository struct {
	repo queueing.Repository
}

func WrapRepository(repo queueing.Repository) queueing.Repository {
	return otelRepository{
		repo: repo,
	}
}

func (o otelRepository) GetByID(ctx context.Context, id uuid.UUID) (*queueing.QueueMessage, error) {
	ctx, span := tracer.Start(ctx, "GetByID")
	defer span.End()

	message, err := o.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
	}

	return message, err
}

func (o otelRepository) Create(ctx context.Context, message *queueing.QueueMessage) error {
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	err := o.repo.Create(ctx, message)
	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (o otelRepository) Update(ctx context.Context, message *queueing.QueueMessage) error {
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()

	err := o.repo.Update(ctx, message)
	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (o otelRepository) Delete(ctx context.Context, message *queueing.QueueMessage) error {
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()

	err := o.repo.Delete(ctx, message)
	if err != nil {
		span.RecordError(err)
	}

	return err
}
