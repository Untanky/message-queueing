package persistence_test

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"message-queueing/persistence"
	"os"
	"testing"
)

type sliceReadWriteSeeker struct {
	loc  int64
	data []byte
}

func (s *sliceReadWriteSeeker) Read(p []byte) (n int, err error) {
	copy(p, s.data[s.loc:])
	s.loc += int64(len(p))
	return len(p), nil
}

func (s *sliceReadWriteSeeker) Write(p []byte) (n int, err error) {
	s.data = append(s.data[:s.loc], p...)
	s.loc += int64(len(p))
	return len(p), nil
}

func (s *sliceReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var newLoc int64
	if whence == io.SeekStart {
		newLoc = offset
	} else if whence == io.SeekCurrent {
		newLoc = s.loc + offset
	} else if whence == io.SeekEnd {
		newLoc = int64(len(s.data)) + offset
	}

	if newLoc < 0 || newLoc > int64(len(s.data)) {
		return -1, errors.New("out of range")
	}

	s.loc = newLoc

	return newLoc, nil
}

func (s *sliceReadWriteSeeker) Close() error {
	return nil
}

type helloWorld struct{}

func (hw helloWorld) Marshal() ([]byte, error) {
	return []byte("Hello World"), nil
}

type trueIterator struct{}

func (it *trueIterator) Next() persistence.Row {
	return persistence.Row{
		Key:   uuid.New(),
		Value: []byte("Hello World! SSTable are amazing and work well for Key-Value-Database"),
	}
}

func (it *trueIterator) HasNext() bool {
	return true
}

func TestSSTableFromIterator(t *testing.T) {
	file, _ := os.OpenFile("test.data", os.O_CREATE|os.O_RDWR, 0600)

	table, err := persistence.SSTableFromIterator(file, &trueIterator{})

	if err != nil {
		t.Errorf("err: expected: nil, got %v", err)
	}
	if table == nil {
		t.Errorf("table: expected not nil; got %v", table)
	}
}
