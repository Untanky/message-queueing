package persistence

import (
	"github.com/google/uuid"
	"io"
)

type dataPageHeader struct {
	pageSpan
	rows       uint32
	rowBytes   uint64
	indexBytes uint64
}

func (header dataPageHeader) Marshal() ([]byte, error) {
	data := make([]byte, 0, headerSize)

	data = append(data, header.startKey[:]...)
	data = append(data, header.endKey[:]...)
	data = byteOrder.AppendUint32(data, header.rows)
	data = byteOrder.AppendUint64(data, header.rowBytes)
	data = byteOrder.AppendUint64(data, header.indexBytes)

	return data, nil
}

type dataPage struct {
	writtenBytes uint64
	rows         []Row
}

func newDataPage() *dataPage {
	return &dataPage{
		rows: make([]Row, 0, 16),
	}
}

func (page *dataPage) addRow(row Row) bool {
	availableBytes := uint64(pageSize-headerSize-indexEntrySize*uint64(len(page.rows))-1) - page.writtenBytes
	rowBytes := row.getByteCount()

	if availableBytes < rowBytes {
		return false
	}

	page.writtenBytes += rowBytes
	page.rows = append(page.rows, row)
	return true
}

func (page *dataPage) get(key uuid.UUID) (Row, bool) {
	if !page.getPageSpan().containsKey(key) {
		return Row{}, false
	}

	rows := page.rows

	for len(rows) > 0 {
		middle := len(rows) / 2
		switch cmp := compareUUID(key, rows[middle].Key); {
		case cmp < 0:
			rows = rows[:middle]
		case cmp == 0:
			return rows[middle], true
		case cmp > 0:
			rows = rows[middle+1:]
		}
	}

	return Row{}, false
}

func (page *dataPage) getPageSpan() pageSpan {
	numberOfRows := len(page.rows)
	if numberOfRows == 0 {
		return pageSpan{}
	}

	return pageSpan{
		startKey: page.rows[0].Key,
		endKey:   page.rows[numberOfRows-1].Key,
	}
}

func (page *dataPage) getPageHeader() dataPageHeader {
	return dataPageHeader{
		pageSpan:   page.getPageSpan(),
		rows:       uint32(len(page.rows)),
		rowBytes:   page.writtenBytes,
		indexBytes: indexEntrySize * uint64(len(page.rows)),
	}
}

func (page *dataPage) WriteTo(writer io.Writer) (int64, error) {
	data := [pageSize]byte{}

	header := page.getPageHeader()
	headerBytes, _ := header.Marshal()
	copy(data[:headerSize], headerBytes)

	rowOffset := headerSize
	indexOffset := uint64(pageSize - header.indexBytes - 1)
	for _, row := range page.rows {
		rowBytes, _ := row.Marshal()
		copy(data[rowOffset:], rowBytes)

		copy(data[indexOffset:], row.Key[:])
		byteOrder.PutUint32(data[indexOffset+16:], uint32(rowOffset))

		rowOffset += uint64(len(rowBytes))
		indexOffset += indexEntrySize
	}

	data[rowOffset] = endOfPage

	n, err := writer.Write(data[:])
	return int64(n), err
}

func (page *dataPage) ReadFrom(reader io.Reader) (int64, error) {
	// TODO: implement
	panic("not implemented")
}
