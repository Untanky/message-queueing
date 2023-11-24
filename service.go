package queueing

import (
	"context"
	"crypto/sha256"
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

func (m *RawQueueMessage) ToQueueMessage() *QueueMessage {
	messageID := uuid.New()
	now := time.Now().Unix()
	hash := sha256.New()
	hash.Write(m.Data)

	return &QueueMessage{
		MessageID:  messageID[:],
		Timestamp:  &now,
		Attributes: m.Attributes,
		Data:       m.Data,
		DataHash:   hash.Sum(nil),
	}
}

type Service interface {
	Enqueue(ctx context.Context, message *QueueMessage) error
	Dequeue(ctx context.Context, messages []*QueueMessage) (int, error)
	Retrieve(ctx context.Context, messages []*QueueMessage) (int, error)
	Acknowledge(ctx context.Context, messageID uuid.UUID) error
}

type Repository interface {
	GetByID(id uuid.UUID) (*QueueMessage, error)
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

func (queue *walQueueService) Dequeue(ctx context.Context, messages []*QueueMessage) (int, error) {
	n, err := queue.service.Dequeue(ctx, messages)

	walErr := json.NewEncoder(queue.log).Encode(
		walEvent{
			EventType: "dequeue",
			Data:      messages,
		},
	)

	if walErr != nil {
		return 0, errors.Join(WalError, FatalDequeueMitigationError)
	}

	return n, err
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
	//TODO implement me
	panic("implement me")
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
	err := queue.repo.Create(message)
	if err != nil {
		return err
	}

	queue.timeoutQueue.Enqueue(time.Now(), MessageId(uuid.Must(uuid.FromBytes(message.MessageID))))

	return nil
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
		message, e := queue.repo.GetByID(uuid.UUID(locations[i]))
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		queue.timeoutQueue.Enqueue(now.Add(defaultDelay), locations[i])
		slice = append(slice, message)
	}

	return len(slice), err
}

func (queue *globalQueueService) Acknowledge(ctx context.Context, messageID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
