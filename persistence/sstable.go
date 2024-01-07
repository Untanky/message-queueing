package persistence

import (
	"encoding/binary"
	"github.com/google/uuid"
	"io"
)

type BinaryMarshaler interface {
	Marshal() ([]byte, error)
}

type BinaryUnmarshaler interface {
	Unmarshal([]byte) error
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
	offset   int64
	startKey uuid.UUID
	endKey   uuid.UUID
}

type pageHeader struct {
	pageSpan
	rows       uint32
	rowBytes   uint32
	indexBytes uint32
}

type KeyValue struct {
	Key   uuid.UUID
	Value BinaryMarshaler
}

var byteOrder = binary.BigEndian

const (
	pageVersion    = 1
	pageSize       = 16 * 1024
	headerSize     = 64
	indexEntrySize = 20
)

func writePage(data *[pageSize]byte, iterator Iterator[KeyValue]) pageSpan {
	index := make([]byte, 64*indexEntrySize)
	offset := headerSize

	header := pageHeader{}

	for iterator.HasNext() && offset+len(index) < pageSize {
		value := iterator.Next()

		valueBytes, err := value.Value.Marshal()
		if err != nil {
			// TODO: error handling
			continue
		}

		if len(valueBytes)+offset+len(index)+indexEntrySize > pageSize {
			break
		}

		if header.rows == 0 {
			header.startKey = value.Key
		}
		header.endKey = value.Key
		header.rows += 1

		copy(data[offset:], valueBytes)
		index = append(index, value.Key[:]...)
		index = byteOrder.AppendUint32(index, uint32(offset))
		offset += len(valueBytes)
	}

	copy(data[pageSize-len(index):], index)

	header.rowBytes = uint32(offset - headerSize)
	header.indexBytes = uint32(len(index))

	// TODO: write header

	return header.pageSpan
}

func SSTableFromIterator(handler ReadWriteSeekCloser, data Iterator[KeyValue]) (*SSTable, error) {
	pageBytes := [pageSize]byte{}
	writePage(&pageBytes, data)
	_, err := handler.Write(pageBytes[:])
	if err != nil {
		return nil, err
	}

	return &SSTable{}, nil
}
