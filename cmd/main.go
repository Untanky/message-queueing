package main

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	queueing "message-queueing"
	"os"
)

func main() {
	MessageID := uuid.NewString()
	Timestamp := int64(23178989641)
	Data := []byte("Hello World")
	DataHash := []byte("FGOOOO")

	file, err := os.OpenFile("data", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(fmt.Errorf("could not open file: %w", err))
	}
	defer file.Close()
	persister := queueing.NewPersister(file)

	message := &queueing.QueueMessage{
		MessageID:  &MessageID,
		Timestamp:  &Timestamp,
		Data:       Data,
		DataHash:   DataHash,
		Attributes: map[string]string{},
	}
	n, err := persister.Write(message)
	if err != nil {
		slog.Error("error writing", "err", err, "n", n)
	}

	message, err = persister.Read(n)
	if err != nil {
		slog.Error("error reading", "err", err)
	}

	slog.Info("Success!", "id", *message.MessageID)
}
