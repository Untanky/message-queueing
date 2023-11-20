package queueing

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"io"
)

type Queue[Value any] interface {
	Enqueue(ctx context.Context, messages ...Value) error
	Dequeue(ctx context.Context, messages []Value) (int, error)
	Acknowledge(ctx context.Context, messageID uuid.UUID) error
}

type Repository interface {
	GetByID(id uuid.UUID) (*QueueMessage, error)
	GetActive(messages []*QueueMessage) (int, error)
	Create(message *QueueMessage) error
	Update(message *QueueMessage) error
	Delete(message *QueueMessage) error
}

type walQueueService struct {
	log io.Writer

	service Queue[*QueueMessage]
}

func NewWriteAheadLogQueueService(writer io.Writer, service Queue[*QueueMessage]) Queue[*QueueMessage] {
	return &walQueueService{
		log:     writer,
		service: service,
	}
}

func (queue *walQueueService) Enqueue(ctx context.Context, messages ...*QueueMessage) error {
	queue.log.Write([]byte{})

	return queue.service.Enqueue(ctx, messages...)
}

func (queue *walQueueService) Dequeue(ctx context.Context, messages []*QueueMessage) (int, error) {
	n, err := queue.service.Dequeue(ctx, messages)

	queue.log.Write([]byte{})

	return n, err
}

func (queue *walQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

type globalQueueService struct {
	repo Repository
}

func NewQueueService(repo Repository) Queue[*QueueMessage] {
	return &globalQueueService{
		repo: repo,
	}
}

func (queue *globalQueueService) Enqueue(ctx context.Context, messages ...*QueueMessage) error {
	var e error
	for _, message := range messages {
		err := queue.repo.Create(message)
		if err != nil {
			e = errors.Join(e, err)
		}
	}

	return e
}

func (queue *globalQueueService) Dequeue(ctx context.Context, messages []*QueueMessage) (int, error) {
	return queue.repo.GetActive(messages)
}

func (queue *globalQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
