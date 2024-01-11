package persistence

import (
	"io"
	"math/rand"
	"time"
)

const SSTableVersion = uint32(1)

type compactionInformation struct {
	table1ID    uint64
	table2ID    uint64
	keysDeleted uint64
	keysUpdated uint64
	keysKept    uint64
}

type pageSpanWithOffset struct {
	offset uint64
	pageSpan
}

type tableHeader struct {
	tableVersion          uint32
	tableID               uint64
	createdAt             time.Time
	compactionInformation *compactionInformation
	spans                 []pageSpanWithOffset
}

// TODO: find a nicer way to do this
var random = rand.New(rand.NewSource(rand.Int63()))
var now = time.Now

func newTableHeader() *tableHeader {
	return &tableHeader{
		tableVersion: SSTableVersion,
		tableID:      random.Uint64(),
		createdAt:    now(),
		spans:        make([]pageSpanWithOffset, 0, 4),
	}
}

func (header *tableHeader) addPage(span pageSpanWithOffset) {
	header.spans = append(header.spans, span)
}

func (header *tableHeader) get(key []byte) pageSpanWithOffset {
	// TODO: implement
	panic("not implemented")
}

func (header *tableHeader) Marshal() ([]byte, error) {
	const pageSpanSize = 40
	size := 24 + len(header.spans)*pageSpanSize
	if header.compactionInformation != nil {
		size += 40
	}

	data := make([]byte, 0, size)

	data = byteOrder.AppendUint32(data, header.tableVersion)
	data = byteOrder.AppendUint64(data, header.tableID)
	data = byteOrder.AppendUint64(data, uint64(header.createdAt.UnixMilli()))

	if header.compactionInformation != nil {
		data = byteOrder.AppendUint64(data, header.compactionInformation.table1ID)
		data = byteOrder.AppendUint64(data, header.compactionInformation.table2ID)
		data = byteOrder.AppendUint64(data, header.compactionInformation.keysDeleted)
		data = byteOrder.AppendUint64(data, header.compactionInformation.keysUpdated)
		data = byteOrder.AppendUint64(data, header.compactionInformation.keysKept)
	}

	for _, span := range header.spans {
		data = append(data, span.startKey...)
		data = append(data, span.endKey...)
		data = byteOrder.AppendUint64(data, span.offset)
	}

	return data, nil
}

func (header *tableHeader) WriteTo(writer io.Writer) (int64, error) {
	headerBytes, _ := header.Marshal()

	n, err := writer.Write(headerBytes)
	return int64(n), err
}

func (header *tableHeader) ReadFrom(reader io.Reader) (int64, error) {
	// TODO: implement
	panic("not implemented")
}
