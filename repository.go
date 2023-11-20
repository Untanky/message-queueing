package queueing

import (
	"encoding/binary"
	"errors"
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

type Persister interface {
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

	persister    Persister
	index        Index[MessageId, MessageLocation]
	timeoutQueue *timeoutQueue
}

func SetupQueueMessageRepository(id string) (Repository, error) {
	file, err := os.OpenFile(fmt.Sprintf("data/%s", id), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	persister := NewPersister(file)
	index := NewNaiveIndex()

	loc := int64(0)
	for {
		message, length, err := readNextMessage(file)
		if err == io.EOF {
			break
		}
		index.Set(MessageId(uuid.MustParse(*message.MessageID)), MessageLocation(loc))
		loc += length + 8
	}

	fmt.Println(index.(*naiveIndex).data)

	repo := NewQueueMessageRepository(persister, index)
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

	loc, ok := q.index.Get(MessageId(id))
	if !ok {
		return nil, NotFoundError
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

func (q queueMessageRepository) GetActive(messages []*QueueMessage) (int, error) {
	locations := make([]MessageLocation, len(messages))
	n, err := q.timeoutQueue.DequeueMultiple(locations)
	actual := n
	for i := 0; i < n; i++ {
		data, e := q.persister.Read(int64(locations[i]))
		if e != nil {
			actual -= 1
			err = errors.Join(err, e)
			continue
		}

		e = proto.Unmarshal(data, messages[actual])
		if e != nil {
			actual -= 1
			err = errors.Join(err, e)
		}
	}

	return n, err
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

	q.index.Set(MessageId(id), MessageLocation(loc))
	q.timeoutQueue.Enqueue(time.Now().Add(defaultDelay), MessageLocation(loc))

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
