package queueing

import (
	"encoding/binary"
	"io"
	"sync"
)

type genericPersister struct {
	lock sync.Locker

	handler io.ReadWriteSeeker
}

func NewPersister(handler io.ReadWriteSeeker) Persister {
	return &genericPersister{
		lock:    &sync.Mutex{},
		handler: handler,
	}
}

func (persister *genericPersister) Write(data []byte) (int64, error) {
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

	return off, nil
}

func (persister *genericPersister) Read(location int64) ([]byte, error) {
	persister.lock.Lock()
	defer persister.lock.Unlock()

	_, err := persister.handler.Seek(location, 0)
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

	return data, err
}
