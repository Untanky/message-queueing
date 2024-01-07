package persistence

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	"io"
	"time"
)

type BinaryMarshaler interface {
	Marshal() ([]byte, error)
}

type BinaryUnmarshaler interface {
	Unmarshal([]byte) error
}

const (
	endOfPage = 0x01
	deleted   = 0x02
	retrieved = 0x04
	dlq       = 0x08
)

type RetrieveInfo struct {
	retrieved       uint32
	lastRetrievedAt time.Time
}

type DeadLetterQueueInfo struct {
	movedAt     time.Time
	originQueue [16]byte
}

type Row struct {
	Key                 uuid.UUID
	DeletedAt           *time.Time
	RetrieveInfo        *RetrieveInfo
	DeadLetterQueueInfo *DeadLetterQueueInfo
	Value               []byte
}

func (row *Row) Marshal() ([]byte, error) {
	if row.DeletedAt != nil {
		data := make([]byte, 0, 9)
		data = append(data, byte(deleted))
		data = binary.AppendUvarint(data, uint64(row.DeletedAt.UnixMilli()))
		return data, nil
	}

	data := make([]byte, 0, len(row.Value)+45)

	data = append(data, 0)
	if row.RetrieveInfo != nil {
		data[0] |= byte(retrieved)
		data = binary.AppendUvarint(data, uint64(row.RetrieveInfo.retrieved))
		data = binary.AppendUvarint(data, uint64(row.RetrieveInfo.lastRetrievedAt.UnixMilli()))
	}
	if row.DeadLetterQueueInfo != nil {
		data[0] |= byte(dlq)
		data = binary.AppendUvarint(data, uint64(row.DeadLetterQueueInfo.movedAt.UnixMilli()))
		data = append(data, row.DeadLetterQueueInfo.originQueue[:]...)
	}
	data = binary.AppendUvarint(data, uint64(len(row.Value)))
	data = append(data, row.Value...)
	return data, nil
}

func (row *Row) writeTo(dest []byte, offset uint32) (int, []byte, error) {
	valueBytes, err := row.Marshal()
	if err != nil {
		return 0, nil, err
	}

	if len(valueBytes)+indexEntrySize > len(dest) {
		return 0, nil, errors.New("not enough space")
	}

	n := copy(dest, valueBytes)
	index := make([]byte, 0, 20)
	index = append(index, row.Key[:]...)
	index = binary.AppendUvarint(index, uint64(offset))

	return n, index, nil
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
	offset   uint64
	startKey uuid.UUID
	endKey   uuid.UUID
}

type pageHeader struct {
	pageSpan
	rows       uint32
	rowBytes   uint32
	indexBytes uint32
}

func (header *pageHeader) Marshal() ([]byte, error) {
	data := make([]byte, 0, headerSize)

	data = byteOrder.AppendUint64(data, header.offset)
	data = append(data, header.startKey[:]...)
	data = append(data, header.endKey[:]...)
	data = byteOrder.AppendUint32(data, header.rows)
	data = byteOrder.AppendUint32(data, header.rowBytes)
	data = byteOrder.AppendUint32(data, header.indexBytes)

	return data, nil
}

var byteOrder = binary.BigEndian

const (
	pageSize       = 16 * 1024
	headerSize     = 64
	indexEntrySize = 20
)

func writePage(data *[pageSize]byte, iterator Iterator[Row]) pageSpan {
	index := make([]byte, 0, 64*indexEntrySize)
	offset := headerSize
	lastFlagOffset := headerSize

	header := pageHeader{}

	for iterator.HasNext() && offset+len(index) < pageSize {
		value := iterator.Next()

		n, indexAppend, err := value.writeTo(data[offset:pageSize-len(index)], uint32(offset))
		if err != nil {
			break
		}
		lastFlagOffset = offset
		offset += n
		index = append(index, indexAppend...)

		if header.rows == 0 {
			header.startKey = value.Key
		}
		header.endKey = value.Key
		header.rows += 1
	}

	data[lastFlagOffset] |= endOfPage
	copy(data[pageSize-len(index):], index)

	header.rowBytes = uint32(offset - headerSize)
	header.indexBytes = uint32(len(index))

	headerBytes, _ := header.Marshal()

	copy(data[0:headerSize], headerBytes)

	return header.pageSpan
}

func SSTableFromIterator(handler ReadWriteSeekCloser, data Iterator[Row]) (*SSTable, error) {
	pageBytes := [pageSize]byte{}
	writePage(&pageBytes, data)
	_, err := handler.Write(pageBytes[:])
	if err != nil {
		return nil, err
	}

	return &SSTable{}, nil
}
