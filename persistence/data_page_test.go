package persistence

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"message-queueing/testutils"
	"testing"
)

type trueIterator struct{}

func (it *trueIterator) Next() Row {
	id := uuid.New()
	return Row{
		Key:   id[:],
		Value: []byte("Hello World! SSTable are amazing and work well for Key-Value-Database"),
	}
}

func (it *trueIterator) HasNext() bool {
	return true
}

func fillPage(handler io.WriteSeeker, data Iterator[Row]) (*dataPage, error) {
	page := newDataPage()

	for data.HasNext() {
		ok := page.addRow(data.Next())
		if !ok {
			break
		}
	}

	_, err := page.WriteTo(handler)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func TestDataPage_WriteTo(t *testing.T) {
	const expectedHash = "8ay4rwfoEjUb/MbCUWS3rquqjWNY+oM+eVuN0PAk/YI="

	sliceIO := &testutils.SliceReadWriteSeeker{}
	uuid.SetRand(rand.New(rand.NewSource(10)))

	table, err := fillPage(sliceIO, &trueIterator{})

	if err != nil {
		t.Errorf("err: expected: nil, got %v", err)
	}
	if table == nil {
		t.Errorf("table: expected not nil; got %v", table)
	}

	hash := sha256.New()
	hash.Write(sliceIO.Data)
	hashBytes := hash.Sum(nil)

	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	if hashBase64 != expectedHash {
		t.Errorf("hashBytes: expected %v; got %v", expectedHash, hashBase64)
	}
}
