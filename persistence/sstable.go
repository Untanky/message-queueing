package persistence

import (
	"encoding/binary"
	"io"
	"time"
)

const (
	endOfPage = byte(0x01)
	deleted   = byte(0x02)
	retrieved = byte(0x04)
	dlq       = byte(0x08)

	pageSize       = uint64(16 * 1024)
	headerSize     = uint64(64)
	indexEntrySize = uint64(20)
)

var byteOrder = binary.BigEndian

type RetrieveInfo struct {
	retrieved       uint32
	lastRetrievedAt time.Time
}

type DeadLetterQueueInfo struct {
	movedAt     time.Time
	originQueue [16]byte
}

type Row struct {
	Key                 []byte
	DeletedAt           *time.Time
	RetrieveInfo        *RetrieveInfo
	DeadLetterQueueInfo *DeadLetterQueueInfo
	Value               []byte
}

func (row *Row) getByteCount() uint64 {
	bytes := uint64(1)

	if row.DeletedAt != nil {
		return bytes + 8
	}

	bytes += uint64(len(row.Value)) + 8

	if row.RetrieveInfo != nil {
		bytes += 12
	}
	if row.DeadLetterQueueInfo != nil {
		bytes += 24
	}

	return bytes
}

func (row *Row) Marshal() ([]byte, error) {
	if row.DeletedAt != nil {
		data := make([]byte, 0, 9)
		data = append(data, deleted)
		data = byteOrder.AppendUint64(data, uint64(row.DeletedAt.UnixMilli()))
		return data, nil
	}

	data := make([]byte, 0, len(row.Value)+45)

	data = append(data, 0)
	if row.RetrieveInfo != nil {
		data[0] |= retrieved
		data = byteOrder.AppendUint32(data, row.RetrieveInfo.retrieved)
		data = byteOrder.AppendUint64(data, uint64(row.RetrieveInfo.lastRetrievedAt.UnixMilli()))
	}
	if row.DeadLetterQueueInfo != nil {
		data[0] |= dlq
		data = byteOrder.AppendUint64(data, uint64(row.DeadLetterQueueInfo.movedAt.UnixMilli()))
		data = append(data, row.DeadLetterQueueInfo.originQueue[:]...)
	}
	data = byteOrder.AppendUint64(data, uint64(len(row.Value)))
	data = append(data, row.Value...)
	return data, nil
}

type Iterator[Value any] interface {
	Next() Value
	HasNext() bool
}

type SSTable struct {
}

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

type pageSpan struct {
	startKey []byte
	endKey   []byte
}

func compareUUID(a, b []byte) byte {
	l := min(len(a), len(b))
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return a[i] - b[i]
		}
	}

	return byte(len(a) - len(b))
}

func (span pageSpan) containsKey(key []byte) bool {
	return compareUUID(key, span.startKey) > 0 && compareUUID(key, span.endKey) < 0
}

func SSTableFromIterator(handler ReadWriteSeekCloser, data Iterator[Row]) (*SSTable, error) {
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

	return &SSTable{}, nil
}
