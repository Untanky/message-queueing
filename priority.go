package queueing

import (
	"errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Queue interface {
	Enqueue(messages ...*QueueMessage) error
	Dequeue(messages []*QueueMessage) (int, error)
	Acknowledge(messageID uuid.UUID) error
}

type MessageId uuid.UUID
type MessageLocation uint64

type Persister interface {
	Write([]byte) (int64, error)
	Read(location int64) ([]byte, error)
}

type Index[Key comparable, Value any] interface {
	Get(id MessageId) (MessageLocation, error)
	Set(id MessageId, location MessageLocation) error
	Delete(id MessageId) (MessageLocation, error)
}

type globalQueueService struct {
	lock              sync.Locker
	acknowledgeBuffer []uuid.UUID

	priorityQueue Queue

	persister    Persister
	primaryIndex Index[MessageId, MessageLocation]
}

func (queue *globalQueueService) Enqueue(messages ...*QueueMessage) error {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	var e error
	for _, message := range messages {
		data, err := proto.Marshal(message)
		if err != nil {
			e = errors.Join(e, err)
			continue
		}

		location, err := queue.persister.Write(data)
		if err != nil {
			e = errors.Join(e, err)
			continue
		}

		err = queue.primaryIndex.Set(MessageId(uuid.MustParse(*message.MessageID)), MessageLocation(location))
		if err != nil {
			e = errors.Join(e, err)
			continue
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
