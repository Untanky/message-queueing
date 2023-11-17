package main

import (
	"github.com/google/uuid"
	"log/slog"
	queueing "message-queueing"
)

func main() {
	MessageID := uuid.NewString()
	Timestamp := int64(23178989641)
	Data := []byte("Hello World")
	DataHash := []byte("FGOOOO")

	persister, err := queueing.NewFilePersister()
	if err != nil {
		slog.Error("could not create persister", "err", err)
	}

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
