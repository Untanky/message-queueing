package persistence

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"math/rand"
	"message-queueing/testutils"
	"testing"
	"time"
)

type countIterator struct {
	count    int
	maxCount int
}

func (it *countIterator) Next() Row {
	if it.count >= it.maxCount {
		return Row{}
	}

	it.count++
	id := uuid.New()
	return Row{
		Key:   id[:],
		Value: []byte("Hello World! SSTable are amazing and work well for Key-Value-Database"),
	}
}

func (it *countIterator) HasNext() bool {
	return it.count < it.maxCount
}

func TestSSTableFromIterator(t *testing.T) {
	const expectedHash = "DTSdg/IM9KCnGIhIS95xu6Cgb4sXB7yTA8B78GA2R4k="

	sliceIO := &testutils.SliceReadWriteSeeker{}
	random = rand.New(rand.NewSource(10))
	now = func() time.Time {
		return time.Date(2024, 1, 10, 14, 40, 0, 0, time.Local)
	}
	uuid.SetRand(random)

	table, err := SSTableFromIterator(sliceIO, &countIterator{maxCount: 500})
	if table == nil {
		t.Errorf("table: expected not nil; got %v", table)
	}
	if err != nil {
		t.Errorf("table: expected nil; got %v", err)
	}

	hash := sha256.New()
	hash.Write(sliceIO.Data)
	hashBytes := hash.Sum(nil)

	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	if hashBase64 != expectedHash {
		t.Errorf("hashBytes: expected %v; got %v", expectedHash, hashBase64)
	}
}
