package queueing

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
	"time"
)

type MessageId uuid.UUID
type MessageLocation uint64

type BlockStorage interface {
	WriteBlock([]byte) (int64, error)
	OverwriteBlock(location int64, data []byte) error
	ReadBlock(location int64) ([]byte, error)

	Length() int64
	io.ReaderFrom
	io.WriterTo
}

type Index[Key comparable, Value any] interface {
	Get(id MessageId) (MessageLocation, bool)
	Set(id MessageId, location MessageLocation)
	Delete(id MessageId) (MessageLocation, bool)
}

const defaultDelay = time.Duration(1 * time.Minute)

type queueMessageRepository struct {
	lock sync.Locker

	storage BlockStorage
	index   Index[MessageId, MessageLocation]
}

func NewQueueMessageRepository(
	storage BlockStorage, index Index[MessageId, MessageLocation],
) Repository {
	return &queueMessageRepository{
		lock: &sync.Mutex{},

		storage: storage,
		index:   index,
	}
}

func (q queueMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*QueueMessage, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	loc, ok := q.index.Get(MessageId(id))
	if !ok {
		return nil, NotFoundError
	}

	data, err := q.storage.ReadBlock(int64(loc))
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

func (q queueMessageRepository) Create(ctx context.Context, message *QueueMessage) error {
	id, err := uuid.FromBytes(message.MessageID)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	loc, err := q.storage.WriteBlock(data)
	if err != nil {
		return err
	}

	q.index.Set(MessageId(id), MessageLocation(loc))

	return nil
}

func (q queueMessageRepository) Update(ctx context.Context, message *QueueMessage) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	loc, ok := q.index.Get(MessageId(message.MessageID))
	if !ok {
		return NotFoundError
	}

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	return q.storage.OverwriteBlock(int64(loc), data)
}

func (q queueMessageRepository) Delete(ctx context.Context, message *QueueMessage) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	//TODO implement me
	panic("implement me")
}
