package main

import (
	"fmt"
	"github.com/google/uuid"
	queueing "message-queueing"
	"time"
)

func main() {
	repo, err := queueing.SetupQueueMessageRepository("abc")
	if err != nil {
		panic(err)
	}

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
