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

const pageSpanSize = 40

func (header *tableHeader) Marshal() ([]byte, error) {
	size := 28 + len(header.spans)*pageSpanSize
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

	data = byteOrder.AppendUint64(data, uint64(len(header.spans)))

	for _, span := range header.spans {
		data = append(data, span.startKey...)
		data = append(data, span.endKey...)
		data = byteOrder.AppendUint64(data, span.offset)
	}

	return data, nil
}

func (header *tableHeader) Unmarshal(data []byte) error {
	header.tableVersion = byteOrder.Uint32(data[:4])
	header.tableID = byteOrder.Uint64(data[4:12])
	header.createdAt = time.UnixMilli(int64(byteOrder.Uint64(data[12:20])))

	// TODO: implement reading compaction information
	l := byteOrder.Uint64(data[20:28])

	header.spans = make([]pageSpanWithOffset, 0, l)

	byteOffset := 28

	for i := uint64(0); i < l; i++ {
		startKey, endKey := make([]byte, 16), make([]byte, 16)
		copy(startKey, data[byteOffset:byteOffset+16])
		copy(endKey, data[byteOffset+16:byteOffset+32])
		offset := byteOrder.Uint64(data[byteOffset+32 : byteOffset+40])
		header.spans = append(header.spans, pageSpanWithOffset{
			offset: offset,
			pageSpan: pageSpan{
				startKey: startKey,
				endKey:   endKey,
			},
		})
		byteOffset += pageSpanSize
	}

	return nil
}

func (header *tableHeader) WriteTo(writer io.Writer) (int64, error) {
	headerBytes, _ := header.Marshal()

	n, err := writer.Write(headerBytes)
	return int64(n), err
}

func (header *tableHeader) ReadFrom(reader io.Reader) (int64, error) {
	data := make([]byte, pageSize)
	n, err := reader.Read(data)
	if err != nil {
		return int64(n), err
	}

	err = header.Unmarshal(data)
	return int64(n), err
}
