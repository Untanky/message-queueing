package queueing

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
)

type GenericPersister struct {
	lock sync.Locker

	handler io.ReadWriteSeeker
}

func NewPersister(handler io.ReadWriteSeeker) Persister {
	return &GenericPersister{
		lock:    &sync.Mutex{},
		handler: handler,
	}
}

func (persister *GenericPersister) Write(message *QueueMessage) (MessageLocation, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return 0, err
	}
	l := len(data)

	persister.lock.Lock()
	defer persister.lock.Unlock()

	off, err := persister.handler.Seek(0, 2)
	if err != nil {
		return 0, err
	}

	err = binary.Write(persister.handler, binary.BigEndian, uint64(l))
	if err != nil {
		return 0, err
	}

	_, err = persister.handler.Write(data)
	if err != nil {
		return 0, err
	}

	return MessageLocation(off), nil
}

func (persister *GenericPersister) Read(location MessageLocation) (*QueueMessage, error) {
	persister.lock.Lock()
	defer persister.lock.Unlock()

	_, err := persister.handler.Seek(int64(location), 0)
	if err != nil {
		return nil, err
	}

	var length uint64
	err = binary.Read(persister.handler, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = persister.handler.Read(data)
	if err != nil {
		return nil, err
	}

	message := new(QueueMessage)
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}

	return message, err
}
