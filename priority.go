package queueing

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type Queue interface {
	Enqueue(ctx context.Context, messages ...*QueueMessage) error
	Dequeue(ctx context.Context, messages []*QueueMessage) (int, error)
	Acknowledge(ctx context.Context, messageID uuid.UUID) error
}

type Repository interface {
	GetByID(id uuid.UUID) (*QueueMessage, error)
	GetActive(messages []*QueueMessage) (int, error)
	Create(message *QueueMessage) error
	Update(message *QueueMessage) error
	Delete(message *QueueMessage) error
}

type globalQueueService struct {
	repo Repository
}

func NewQueue(repo Repository) Queue {
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
