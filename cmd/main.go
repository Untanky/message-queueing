package main

import (
	"fmt"
	"github.com/google/uuid"
	queueing "message-queueing"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("data", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	persister := queueing.NewPersister(file)
	index := queueing.NewNaiveIndex()

	repo := queueing.NewQueueMessageRepository(persister, index)

	messageA := NewQueueMessage()
	messageB := NewQueueMessage()

	err = repo.Create(messageA)
	if err != nil {
		panic(err)
	}

	err = repo.Create(messageB)
	if err != nil {
		panic(err)
	}

	fmt.Println(repo.GetByID(uuid.MustParse(*messageA.MessageID)))
	fmt.Println(repo.GetByID(uuid.MustParse(*messageB.MessageID)))
}

func NewQueueMessage() *queueing.QueueMessage {
	messageID := uuid.NewString()
	timestamp := time.Now().Unix()
	data := []byte("Hello World")
	dataHash := []byte("abc")

	return &queueing.QueueMessage{
		MessageID:  &messageID,
		Timestamp:  &timestamp,
		Data:       data,
		DataHash:   dataHash,
		Attributes: map[string]string{},
	}
}
