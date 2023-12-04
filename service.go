package queueing

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
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

type walEvent struct {
	EventType string `json:"type"`
	Data      any    `json:"data"`
}

type walQueueService struct {
	log io.Writer

	service Service
}

func NewWriteAheadLogQueueService(writer io.Writer, service Service) Service {
	return &walQueueService{
		log:     writer,
		service: service,
	}
}

func (queue *walQueueService) Enqueue(ctx context.Context, message *QueueMessage) error {
	err := json.NewEncoder(queue.log).Encode(
		walEvent{
			EventType: "enqueue",
			Data:      message,
		},
	)
	if err != nil {
		return err
	}

	return queue.service.Enqueue(ctx, message)
}

func (queue *walQueueService) Retrieve(ctx context.Context, messages []*QueueMessage) (int, error) {
	n, err := queue.service.Retrieve(ctx, messages)

	walErr := json.NewEncoder(queue.log).Encode(
		walEvent{
			EventType: "retrieve",
			Data:      messages,
		},
	)

	if walErr != nil {
		return 0, errors.Join(WalError, FatalDequeueMitigationError)
	}

	return n, err
}

func (queue *walQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	err := json.NewEncoder(queue.log).Encode(
		walEvent{
			EventType: "enqueue",
			Data:      messageID,
		},
	)
	if err != nil {
		return err
	}

	return queue.service.Acknowledge(ctx, messageID)
}

type globalQueueService struct {
	repo         Repository
	timeoutQueue *timeoutQueue
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
