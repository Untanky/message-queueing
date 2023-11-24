package queueing

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"sync"
	"time"
)

type MessageId uuid.UUID
type MessageLocation uint64

type BlockStorage interface {
	Write([]byte) (int64, error)
	Read(location int64) ([]byte, error)
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

func SetupQueueMessageRepository(id string) (Repository, error) {
	file, err := os.OpenFile(fmt.Sprintf("data/%s", id), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	storage := NewIOBlockStorage(file)
	index := NewNaiveIndex()

	loc := int64(0)
	for {
		message, length, err := readNextMessage(file)
		if err == io.EOF {
			break
		}
		index.Set(MessageId(uuid.Must(uuid.FromBytes(message.MessageID))), MessageLocation(loc))
		loc += length + 8
	}

	fmt.Println(index.(*naiveIndex).data)

	repo := NewQueueMessageRepository(storage, index)
	return repo, nil
}

func readNextMessage(reader io.Reader) (*QueueMessage, int64, error) {
	var length int64
	err := binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return nil, 0, err
	}

	data := make([]byte, length)
	_, err = reader.Read(data)
	if err != nil {
		return nil, 0, err
	}

	var message QueueMessage
	err = proto.Unmarshal(data, &message)
	if err != nil {
		return nil, 0, err
	}

	return &message, length, nil
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

func (q queueMessageRepository) GetByID(id uuid.UUID) (*QueueMessage, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	loc, ok := q.index.Get(MessageId(id))
	if !ok {
		return nil, NotFoundError
	}

	data, err := q.storage.Read(int64(loc))
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

func (q queueMessageRepository) GetAvailable(messages []*QueueMessage) (int, error) {
	panic("not implemented")
}

func (q queueMessageRepository) Create(message *QueueMessage) error {
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

	loc, err := q.storage.Write(data)
	if err != nil {
		return err
	}

	q.index.Set(MessageId(id), MessageLocation(loc))

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
