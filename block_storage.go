package queueing

import (
	"encoding/binary"
	"io"
	"sync"
)

type ioBlockStorage struct {
	lock sync.Locker

	handler io.ReadWriteSeeker
}

func NewIOBlockStorage(handler io.ReadWriteSeeker) BlockStorage {
	return &ioBlockStorage{
		lock:    &sync.Mutex{},
		handler: handler,
	}
}

func (storage *ioBlockStorage) Write(data []byte) (int64, error) {
	l := len(data)

	storage.lock.Lock()
	defer storage.lock.Unlock()

	off, err := storage.handler.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	err = binary.Write(storage.handler, binary.BigEndian, uint64(l))
	if err != nil {
		return 0, err
	}

	_, err = storage.handler.Write(data)
	if err != nil {
		return 0, err
	}

	return off, nil
}

func (storage *ioBlockStorage) Read(location int64) ([]byte, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	_, err := storage.handler.Seek(location, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var length uint64
	err = binary.Read(storage.handler, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = storage.handler.Read(data)
	if err != nil {
		return nil, err
	}

	return data, err
}
