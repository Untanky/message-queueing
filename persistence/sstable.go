package persistence

import (
	"encoding/binary"
	"errors"
	"fmt"
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

var (
	NotFoundError      = errors.New("not found")
	MarkedDeletedError = fmt.Errorf("%w: marked for deletion", NotFoundError)
	byteOrder          = binary.BigEndian
)

type RetrieveInfo struct {
	retrieved       uint32
	lastRetrievedAt time.Time
}

type DeadLetterQueueInfo struct {
	movedAt     time.Time
	originQueue []byte
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
	data = append(data, row.Value...)
	return data, nil
}

func (row *Row) Unmarshal(data []byte) error {
	flag := data[0]

	if flag&deleted != 0 {
		deletedAt := time.UnixMilli(int64(byteOrder.Uint64(data[1:9])))
		row.DeletedAt = &deletedAt
	}

	offset := 1

	if flag&retrieved != 0 {
		row.RetrieveInfo = &RetrieveInfo{
			retrieved:       byteOrder.Uint32(data[offset : offset+4]),
			lastRetrievedAt: time.UnixMilli(int64(byteOrder.Uint64(data[offset+4 : offset+12]))),
		}
		offset += 12
	}

	if flag&dlq != 0 {
		queueID := make([]byte, 16)
		copy(queueID, data[offset+8:offset+24])

		row.DeadLetterQueueInfo = &DeadLetterQueueInfo{
			movedAt:     time.UnixMilli(int64(byteOrder.Uint64(data[offset : offset+8]))),
			originQueue: queueID,
		}
		offset += 24
	}

	row.Value = make([]byte, len(data)-offset)
	copy(row.Value, data[offset:])

	return nil
}

type pageSpan struct {
	startKey []byte
	endKey   []byte
}

func compareBytes(a, b []byte) int {
	l := min(len(a), len(b))
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}

	return len(a) - len(b)
}

func (span pageSpan) containsKey(key []byte) bool {
	startKey := compareBytes(key, span.startKey)
	endKey := compareBytes(key, span.endKey)
	return startKey >= 0 && endKey <= 0
}

type SSTable struct {
	lock   sync.Mutex
	header *tableHeader
	reader io.ReadSeekCloser
}

func CreateSSTable(id uint64, createdAt time.Time, handler io.WriteSeeker, data Iterator[Row]) error {
	header := newTableHeader(id, createdAt)

	// write empty header for spacing
	initialOffset, err := handler.Write(make([]byte, pageSize))
	if err != nil {
		return err
	}

	table := &temporarySSTable{
		header: header,
		writer: handler,
		offset: int64(initialOffset),
	}

	err = table.writeAllPages(data)
	if err != nil {
		return err
	}

	_, err = handler.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = header.WriteTo(handler)
	if err != nil {
		return err
	}

	return nil
}

type temporarySSTable struct {
	writer io.Writer
	header *tableHeader
	offset int64
}

func (table *temporarySSTable) writeAllPages(data Iterator[Row]) error {
	var err error
	page := newDataPage()

	for data.HasNext() {
		current := data.Next()
		ok := page.addRow(current)

		if !ok {
			err = table.writePage(page)
			if err != nil {
				return err
			}

			page = newDataPage()
			page.addRow(current)
		}
	}

	err = table.writePage(page)
	if err != nil {
		return err
	}

	return nil
}

func (table *temporarySSTable) writePage(page *dataPage) error {
	n, err := page.WriteTo(table.writer)
	if err != nil {
		return err
	}

	table.header.addPage(pageSpanWithOffset{offset: uint64(table.offset), pageSpan: page.getPageSpan()})
	table.offset += n

	return err
}

func NewSSTable(reader io.ReadSeekCloser) (*SSTable, error) {
	// set reader to offset for page
	_, err := reader.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	header := newTableHeader(0, time.Now())
	_, err = header.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	return &SSTable{
		lock:   sync.Mutex{},
		reader: reader,
		header: header,
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
	span, ok := table.header.get(key)
	if !ok {
		return Row{}, NotFoundError
	}

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
		return Row{}, NotFoundError
	}

	if row.DeletedAt != nil {
		return Row{}, MarkedDeletedError
	}

	return row, nil
}

func (table *SSTable) Close() error {
	table.lock.Lock()
	defer table.lock.Unlock()

	return table.reader.Close()
}
