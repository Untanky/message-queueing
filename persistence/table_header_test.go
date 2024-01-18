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

func TestTableHeader_Marshal(t *testing.T) {
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

	bytes, err := header.Marshal()
	if err != nil {
		t.Errorf("header: expected nil; got %v", header)
	}

	unmarshaledHeader := newTableHeader()
	err = unmarshaledHeader.Unmarshal(bytes)
	if err != nil {
		t.Errorf("header: expected nil; got %v", header)
	}

	if header.tableVersion != unmarshaledHeader.tableVersion {
		t.Errorf("header.tableVersion: unmarshaled %v; expected: %v", unmarshaledHeader.tableVersion, header.tableVersion)
	}

	if header.tableID != unmarshaledHeader.tableID {
		t.Errorf("header.tableID: unmarshaled %v; expected: %v", unmarshaledHeader.tableID, header.tableID)
	}

	if header.createdAt.Equal(unmarshaledHeader.createdAt) {
		t.Errorf("header.createdAt: unmarshaled %v; expected: %v", unmarshaledHeader.createdAt, header.createdAt)
	}

	if header.compactionInformation != unmarshaledHeader.compactionInformation {
		t.Errorf("header.compactionInformation: unmarshaled %v; expected: %v", unmarshaledHeader.compactionInformation, header.compactionInformation)
	}

	if len(header.spans) != len(unmarshaledHeader.spans) {
		t.Errorf("header.spans: unmarshaled %v; expected: %v", len(unmarshaledHeader.spans), len(header.spans))
	}
}

func TestTableHeader_WriteTo(t *testing.T) {
	const expectedHash = "5QBbcq9k28DBedkVO+vLFJjcoAXAMgHPoN+p8mmnwyg="

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
