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

func TestTableHeader_WriteTo(t *testing.T) {
	const expectedHash = "U6eO6cx5MtW/3g11CB384RlwOe4oK8Pv8fCWxw48LGA="

	sliceIO := &testutils.SliceReadWriteSeeker{}
	random = rand.New(rand.NewSource(10))
	now = func() time.Time {
		return time.Date(2024, 1, 10, 14, 40, 0, 0, time.Local)
	}
	uuid.SetRand(random)

	header := newTableHeader()
	if header == nil {
		t.Errorf("header: expected not nil; got %v", header)
	}

	for i := uint64(1); i < 17; i++ {
		id1 := uuid.New()
		id2 := uuid.New()

		header.addPage(pageSpanWithOffset{
			pageSpan: pageSpan{
				startKey: id1[:],
				endKey:   id2[:],
			},
			offset: i * pageSize,
		})
	}

	header.WriteTo(sliceIO)

	hash := sha256.New()
	hash.Write(sliceIO.Data)
	hashBytes := hash.Sum(nil)

	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	if hashBase64 != expectedHash {
		t.Errorf("hashBytes: expected %v; got %v", expectedHash, hashBase64)
	}
}
