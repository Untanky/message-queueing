package queueing

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"sync"
)

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

type queueMessageRepository struct {
	lock sync.Locker

	persister Persister
	index     Index[MessageId, MessageLocation]
}

func NewQueueMessageRepository(persister Persister, index Index[MessageId, MessageLocation]) Repository {
	return &queueMessageRepository{
		lock: &sync.Mutex{},

		persister: persister,
		index:     index,
	}
}

func (q queueMessageRepository) GetByID(id uuid.UUID) (*QueueMessage, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	loc, err := q.index.Get(MessageId(id))
	if err != nil {
		return nil, err
	}

	data, err := q.persister.Read(int64(loc))
	if err != nil {
		return nil, err
	}

	var queueMessage QueueMessage
	err = proto.Unmarshal(data, &queueMessage)
	if err != nil {
		return nil, err
	}

	return &queueMessage, nil
}

func (q queueMessageRepository) Create(message *QueueMessage) error {
	id, err := uuid.Parse(*message.MessageID)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	loc, err := q.persister.Write(data)
	if err != nil {
		return err
	}

	err = q.index.Set(MessageId(id), MessageLocation(loc))
	if err != nil {
		return err
	}

	return nil
}

func (q queueMessageRepository) Update(message *QueueMessage) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	//TODO implement me
	panic("implement me")
}

func (q queueMessageRepository) Delete(message *QueueMessage) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	//TODO implement me
	panic("implement me")
}
