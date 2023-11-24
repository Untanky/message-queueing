package queueing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

var (
	WalError                    = errors.New("error writing to write ahead log")
	FatalDequeueMitigationError = errors.New("fatal error! writing to write ahead log failed AND could not enqueue messages again! potential data loss")
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

type walEvent struct {
	EventType string `json:"type"`
	Data      any    `json:"data"`
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
	err := json.NewEncoder(queue.log).Encode(
		walEvent{
			EventType: "enqueue",
			Data:      messages,
		},
	)
	if err != nil {
		return err
	}

	return queue.service.Enqueue(ctx, messages...)
}

func (queue *walQueueService) Dequeue(ctx context.Context, messages []*QueueMessage) (int, error) {
	n, err := queue.service.Dequeue(ctx, messages)

	walErr := json.NewEncoder(queue.log).Encode(
		walEvent{
			EventType: "dequeue",
			Data:      messages,
		},
	)

	if walErr != nil {
		enqErr := queue.service.Enqueue(ctx, messages...)
		if enqErr != nil {
			return 0, errors.Join(walErr, enqErr, WalError, FatalDequeueMitigationError)
		}
		return 0, errors.Join(WalError, FatalDequeueMitigationError)
	}

	return n, err
}

func (queue *walQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

type globalQueueService struct {
	repo         Repository
	timeoutQueue *timeoutQueue
}

func NewQueueService(repo Repository, timeoutQueue *timeoutQueue) Queue[*QueueMessage] {
	return &globalQueueService{
		repo:         repo,
		timeoutQueue: timeoutQueue,
	}
}

func (queue *globalQueueService) Enqueue(ctx context.Context, messages ...*QueueMessage) error {
	var e error
	for _, message := range messages {
		err := queue.repo.Create(message)
		if err != nil {
			e = errors.Join(e, err)
		}
		id, err := uuid.FromBytes(message.MessageID)
		fmt.Println(message.MessageID, id, err)
		queue.timeoutQueue.Enqueue(time.Now(), MessageId(uuid.Must(uuid.FromBytes(message.MessageID))))
	}

	return e
}

func (queue *globalQueueService) Dequeue(ctx context.Context, messages []*QueueMessage) (int, error) {
	locations := make([]MessageId, len(messages))
	messages = messages[:0]
	n, err := queue.timeoutQueue.DequeueMultiple(locations, time.Now())
	for i := 0; i < n; i++ {
		message, e := queue.repo.GetByID(uuid.UUID(locations[i]))
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		messages = append(messages, message)
	}

	return len(messages), err
}

func (queue *globalQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
