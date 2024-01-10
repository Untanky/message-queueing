package persistence_test

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"io"
	"math/rand"
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

type trueIterator struct{}

func (it *trueIterator) Next() persistence.Row {
	uuid.SetRand(rand.New(rand.NewSource(10)))
	return persistence.Row{
		Key:   uuid.New(),
		Value: []byte("Hello World! SSTable are amazing and work well for Key-Value-Database"),
	}
}

func (it *trueIterator) HasNext() bool {
	return true
}

func TestSSTableFromIterator(t *testing.T) {
	sliceIO := &sliceReadWriteSeeker{}

	table, err := persistence.SSTableFromIterator(sliceIO, &trueIterator{})

	if err != nil {
		t.Errorf("err: expected: nil, got %v", err)
	}
	if table == nil {
		t.Errorf("table: expected not nil; got %v", table)
	}

	hash := sha256.New()
	hash.Write(sliceIO.data)
	hashBytes := hash.Sum(nil)

	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	file, _ := os.OpenFile("fromIteratorHash.sha256", os.O_RDONLY, 0600)
	fileBytes := make([]byte, 44)
	file.Read(fileBytes)
	expectedBase64 := string(fileBytes)

	if hashBase64 != expectedBase64 {
		t.Errorf("hashBytes: expected %v; got %v", expectedBase64, hashBase64)
	}
}
