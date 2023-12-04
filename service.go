package queueing

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"io"
	"time"
)

var (
	WalError                    = errors.New("error writing to write ahead log")
	FatalDequeueMitigationError = errors.New("fatal error! writing to write ahead log failed AND could not enqueue messages again! potential data loss")
)

type Service interface {
	Enqueue(ctx context.Context, message *QueueMessage) error
	Retrieve(ctx context.Context, messages []*QueueMessage) (int, error)
	Acknowledge(ctx context.Context, messageID uuid.UUID) error
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*QueueMessage, error)
	Create(ctx context.Context, message *QueueMessage) error
	Update(ctx context.Context, message *QueueMessage) error
	Delete(ctx context.Context, message *QueueMessage) error
}

type globalQueueService struct {
	repo         Repository
	timeoutQueue *timeoutQueue
}

func SetupQueueMessageRepository(
	file io.ReadWriteSeeker, index Index[MessageId, MessageLocation], queue *timeoutQueue,
) {
	loc := int64(0)
	for {
		message, length, err := readNextMessage(file)
		if err == io.EOF {
			break
		}

		if *message.Acknowledged != true {
			mId := MessageId(uuid.Must(uuid.FromBytes(message.MessageID)))
			index.Set(mId, MessageLocation(loc))

			timeout := time.UnixMilli(*message.Timestamp)
			if message.LastRetrieved != nil {
				timeout = time.UnixMilli(*message.LastRetrieved).Add(defaultDelay)
			}
			queue.Enqueue(timeout, mId)
		}
		loc += length + 8
	}

	fmt.Println(queue.heap)
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

func NewQueueService(repo Repository, timeoutQueue *timeoutQueue) Service {
	return &globalQueueService{
		repo:         repo,
		timeoutQueue: timeoutQueue,
	}
}

func (queue *globalQueueService) Enqueue(ctx context.Context, message *QueueMessage) error {
	err := queue.repo.Create(ctx, message)
	if err != nil {
		return err
	}

	queue.timeoutQueue.Enqueue(time.Now(), MessageId(uuid.Must(uuid.FromBytes(message.MessageID))))

	return nil
}

func (queue *globalQueueService) Retrieve(ctx context.Context, messages []*QueueMessage) (int, error) {
	desired := len(messages)
	count := 0
	var err error

	for count < desired && err == nil {
		additional, e := queue.retrieveMessages(ctx, messages[count:])
		count += additional
		err = errors.Join(err, e)
	}

	return count, err
}

func (queue *globalQueueService) retrieveMessages(ctx context.Context, messages []*QueueMessage) (int, error) {
	locations := make([]MessageId, len(messages))
	slice := messages[:0]

	now := time.Now()
	n, err := queue.timeoutQueue.DequeueMultiple(locations, now)

	for i := 0; i < n; i++ {
		message, e := queue.repo.GetByID(ctx, uuid.UUID(locations[i]))
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		if !*message.Acknowledged {
			queue.timeoutQueue.Enqueue(now.Add(defaultDelay), locations[i])
			slice = append(slice, message)
		}
	}

	return len(slice), err
}

func (queue *globalQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	message, err := queue.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	t := true

	message.Acknowledged = &t

	return queue.repo.Update(ctx, message)
}
