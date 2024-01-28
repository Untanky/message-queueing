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
	const expectedHash = "Lom6kaNn5Xu3js23pVVXtkk0OBtV+zK+geAorMicoLM="

	sliceIO := &testutils.SliceReadWriteSeeker{}
	uuid.SetRand(rand.New(rand.NewSource(10)))

	err := CreateSSTable(10, time.Date(2024, 0, 28, 11, 6, 39, 0, time.Local), sliceIO, &countIterator{maxCount: 500})
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
	uuid.SetRand(rand.New(rand.NewSource(10)))

	err := CreateSSTable(10, time.Date(2024, 0, 28, 11, 6, 39, 0, time.Local), sliceIO, &countIterator{maxCount: 500})
	if err != nil {
		t.Errorf("err: expected nil; got %v", err)
	}

	table, err := NewSSTable(sliceIO)
	if err != nil {
		t.Errorf("err: expected nil; got %v", err)
	}

	type testCase struct {
		key     uuid.UUID
		value   string
		wantErr bool
	}

	cases := []testCase{
		{
			key:     uuid.MustParse("00000000-0000-0000-0000-000000000800"),
			value:   demoString,
			wantErr: false,
		},
		{
			key:     uuid.MustParse("00000000-0000-0000-0000-000000000700"),
			value:   "",
			wantErr: true,
		},
		{
			key:     uuid.MustParse("00000000-0000-0000-0000-000000000900"),
			value:   "",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run("FromKey", func(tt *testing.T) {
			row, err := table.Get(c.key[:])

			if err == nil == c.wantErr {
				t.Errorf("err: expected %v, got: %v", c.wantErr, err)
			}
			if string(row.Value) != c.value {
				t.Errorf("row.Value: expected %v; got %v", c.value, string(row.Value))
			}
		})
	}
}
