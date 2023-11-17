package queueing

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"os"
	"sync"
)

type FilePersister struct {
	lock sync.Locker

	file *os.File
}

func NewFilePersister() (Persister, error) {
	file, err := os.OpenFile("data", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	return &FilePersister{
		lock: &sync.Mutex{},
		file: file,
	}, nil
}

func (persister *FilePersister) Write(message *QueueMessage) (MessageLocation, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return 0, err
	}
	l := len(data)

	off, err := persister.file.Seek(0, 2)
	if err != nil {
		return 0, err
	}

	err = binary.Write(persister.file, binary.BigEndian, uint64(l))
	if err != nil {
		return 0, err
	}

	_, err = persister.file.Write(data)
	if err != nil {
		return 0, err
	}

	return MessageLocation(off), nil
}

func (persister *FilePersister) Read(location MessageLocation) (*QueueMessage, error) {
	_, err := persister.file.Seek(int64(location), 0)
	if err != nil {
		return nil, err
	}

	var length uint64
	err = binary.Read(persister.file, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = persister.file.Read(data)
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
