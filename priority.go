package queueing

import (
	"errors"
	"github.com/google/uuid"
)

type Queue interface {
	Enqueue(messages ...*QueueMessage) error
	Dequeue(messages []*QueueMessage) (int, error)
	Acknowledge(messageID uuid.UUID) error
}

type Repository interface {
	GetByID(id uuid.UUID) (*QueueMessage, error)
	Create(message *QueueMessage) error
	Update(message *QueueMessage) error
	Delete(message *QueueMessage) error
}

type globalQueueService struct {
	acknowledgeBuffer []uuid.UUID

	priorityQueue Queue

	repo Repository
}

func (queue *globalQueueService) Enqueue(messages ...*QueueMessage) error {
	var e error
	for _, message := range messages {
		err := queue.repo.Create(message)
		if err != nil {
			e = errors.Join(e, err)
		}
	}

	return e
}

func (queue *globalQueueService) Dequeue(messages []*QueueMessage) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (queue *globalQueueService) Acknowledge(messageID uuid.UUID) error {
	queue.acknowledgeBuffer = append(queue.acknowledgeBuffer, messageID)

	return nil
}
