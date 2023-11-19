package main

import (
	"fmt"
	queueing "message-queueing"
	"time"
)

func main() {
	queue := queueing.NewHeapQueue()
	queue.Enqueue(time.Now().Add(time.Duration(-1_000_000_000)), queueing.MessageLocation(0))
	queue.Enqueue(time.Now().Add(time.Duration(-2_000_000_000)), queueing.MessageLocation(70))
	queue.Enqueue(time.Now().Add(time.Duration(-3_000_000_000)), queueing.MessageLocation(140))

	fmt.Println(queue.Dequeue())
}

//
//func main() {
//	repo, err := queueing.SetupQueueMessageRepository("abc")
//	if err != nil {
//		panic(err)
//	}
//
//	messageA := NewQueueMessage()
//	messageB := NewQueueMessage()
//
//	err = repo.Create(messageA)
//	if err != nil {
//		panic(err)
//	}
//
//	err = repo.Create(messageB)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(repo.GetByID(uuid.MustParse(*messageA.MessageID)))
//	fmt.Println(repo.GetByID(uuid.MustParse(*messageB.MessageID)))
//}
//
//func NewQueueMessage() *queueing.QueueMessage {
//	messageID := uuid.NewString()
//	timestamp := time.Now().UnixMicro()
//	data := []byte("Hello World")
//	dataHash := []byte("abc")
//
//	return &queueing.QueueMessage{
//		MessageID:  &messageID,
//		Timestamp:  &timestamp,
//		Data:       data,
//		DataHash:   dataHash,
//		Attributes: map[string]string{},
//	}
//}
