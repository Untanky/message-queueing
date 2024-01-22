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

const demoString = "Hello World! SSTable are amazing and work well for Key-Value-Database"

type countIterator struct {
	count    int
	maxCount int
}

func (it *countIterator) Next() Row {
	if it.count >= it.maxCount {
		return Row{}
	}

	it.count++

	return Row{
		Key:   byteOrder.AppendUint64(make([]byte, 8, 16), uint64(it.count)*2048),
		Value: []byte(demoString),
	}
}

func (it *countIterator) HasNext() bool {
	return it.count < it.maxCount
}

func TestSSTableFromIterator(t *testing.T) {
	const expectedHash = "2Xw+UeQK1M1LeefwnFwPKchyOwGqrFCG7y4lTdMqUD8="

	sliceIO := &testutils.SliceReadWriteSeeker{}
	random = rand.New(rand.NewSource(10))
	now = func() time.Time {
		return time.Date(2024, 1, 10, 14, 40, 0, 0, time.Local)
	}
	uuid.SetRand(random)

	err := CreateSSTable(sliceIO, &countIterator{maxCount: 500})
	if err != nil {
		t.Errorf("err: expected nil; got %v", err)
	}

	hash := sha256.New()
	hash.Write(sliceIO.Data)
	hashBytes := hash.Sum(nil)

	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	if hashBase64 != expectedHash {
		t.Errorf("hashBytes: expected %v; got %v", expectedHash, hashBase64)
	}
}

func TestSSTable_Get(t *testing.T) {
	sliceIO := &testutils.SliceReadWriteSeeker{}
	random = rand.New(rand.NewSource(10))
	now = func() time.Time {
		return time.Date(2024, 1, 10, 14, 40, 0, 0, time.Local)
	}
	uuid.SetRand(random)

	err := CreateSSTable(sliceIO, &countIterator{maxCount: 500})
	if err != nil {
		t.Errorf("err: expected nil; got %v", err)
	}

	table, err := NewSSTable(sliceIO)
	if err != nil {
		t.Errorf("err: expected nil; got %v", err)
	}

	id := uuid.MustParse("00000000-0000-0000-0000-000000000800")
	row, err := table.Get(id[:])

	if string(row.Value) != demoString {
		t.Errorf("row.Value: expected %v; got %v", demoString, string(row.Value))
	}
}
