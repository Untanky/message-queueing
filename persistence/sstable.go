package persistence

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
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

func compareBytes(a, b []byte) byte {
	l := min(len(a), len(b))
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return a[i] - b[i]
		}
	}

	return byte(len(a) - len(b))
}

func (span pageSpan) containsKey(key []byte) bool {
	return compareBytes(key, span.startKey) > 0 && compareBytes(key, span.endKey) < 0
}

type SSTable struct {
	lock   sync.Mutex
	header *tableHeader
	reader io.ReadSeekCloser
}

func SSTableFromIterator(handler ReadWriteSeekCloser, data Iterator[Row]) (*SSTable, error) {
	header := newTableHeader()

	initialOffset, err := handler.Write(make([]byte, pageSize))
	if err != nil {
		return nil, err
	}

	offset := int64(initialOffset)
	page := newDataPage()

	for data.HasNext() {
		current := data.Next()
		ok := page.addRow(current)

		if !ok {
			var n int64
			n, err = page.WriteTo(handler)
			if err != nil {
				return nil, err
			}

			offset += n
			header.addPage(pageSpanWithOffset{offset: uint64(offset), pageSpan: page.getPageSpan()})

			page = newDataPage()
			page.addRow(current)
		}
	}

	n, err := page.WriteTo(handler)
	if err != nil {
		return nil, err
	}

	offset += n
	header.addPage(pageSpanWithOffset{offset: uint64(offset), pageSpan: page.getPageSpan()})

	_, err = handler.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	_, err = header.WriteTo(handler)
	if err != nil {
		return nil, err
	}

	return &SSTable{
		header: header,
		reader: handler,
	}, nil
}

func (table *SSTable) loadPageFromSpan(span pageSpanWithOffset) (*dataPage, error) {
	page := newDataPage()

	table.lock.Lock()
	defer table.lock.Unlock()

	// set reader to offset for page
	_, err := table.reader.Seek(int64(span.offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	// read the page from the given offset off of the reader
	_, err = page.ReadFrom(table.reader)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (table *SSTable) Get(key []byte) (Row, error) {
	span := table.header.get(key)

	// TODO: check if the cached page is already the required page
	// page, ok := table.getPageFromCache(span)
	// if !ok { ... }

	page, err := table.loadPageFromSpan(span)
	if err != nil {
		return Row{}, err
	}

	// TODO: cache the page
	// table.cachePage(page)

	row, ok := page.get(key)
	if !ok {
		return Row{}, errors.New("not found")
	}

	if row.DeletedAt != nil {
		return Row{}, errors.New("row marked deleted")
	}

	return row, nil
}

func (table *SSTable) Close() error {
	table.lock.Lock()
	defer table.lock.Unlock()

	return table.reader.Close()
}
