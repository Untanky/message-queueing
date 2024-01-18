package persistence

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"message-queueing/testutils"
	"testing"
	"time"
)

func areTableHeaderEqual(t *testing.T, headerA, headerB *tableHeader) {
	if headerA.tableVersion != headerB.tableVersion {
		t.Errorf("tableVersion: tableA %v; tableB: %v", headerA.tableVersion, headerB.tableVersion)
	}

	if headerA.tableID != headerB.tableID {
		t.Errorf("tableID: tableA %v; tableB: %v", headerA.tableID, headerB.tableID)
	}

	if !headerA.createdAt.Equal(headerB.createdAt) {
		t.Errorf("createdAt: tableA %v; tableB: %v", headerA.createdAt, headerB.createdAt)
	}

	if headerA.compactionInformation != headerB.compactionInformation {
		t.Errorf("compactionInformation: tableA %v; tableB: %v", headerA.compactionInformation, headerB.compactionInformation)
	}

	if len(headerA.spans) != len(headerB.spans) {
		t.Errorf("len(spans): tableA %v; tableB: %v", len(headerA.spans), len(headerB.spans))
	}
}

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

	areTableHeaderEqual(t, header, unmarshaledHeader)
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

	sliceIO.Seek(0, io.SeekStart)
	readHeader := newTableHeader()
	if readHeader == nil {
		t.Errorf("header: expected not nil; got %v", header)
	}
	readHeader.ReadFrom(sliceIO)

	areTableHeaderEqual(t, header, readHeader)
}
