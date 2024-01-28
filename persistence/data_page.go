package persistence

import (
	"io"
)

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
	availableBytes := pageSize - headerSize - indexEntrySize*uint64(len(page.rows)) - 1 - page.writtenBytes
	rowBytes := row.getByteCount()

	if availableBytes < rowBytes {
		return false
	}

	page.writtenBytes += rowBytes
	page.rows = append(page.rows, row)
	return true
}

func (page *dataPage) get(key []byte) (Row, bool) {
	if !page.getPageSpan().containsKey(key) {
		return Row{}, false
	}

	rows := page.rows

	for len(rows) > 0 {
		middle := len(rows) / 2
		switch cmp := compareBytes(key, rows[middle].Key); {
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

func (page *dataPage) WriteTo(writer io.Writer) (int64, error) {
	data := [pageSize]byte{}

	span := page.getPageSpan()
	rows := uint32(len(page.rows))
	indexBytes := uint64(rows) * indexEntrySize
	copy(data[:16], span.startKey)
	copy(data[16:32], span.endKey)
	byteOrder.PutUint32(data[32:36], rows)
	byteOrder.PutUint64(data[36:44], page.writtenBytes)
	byteOrder.PutUint64(data[44:52], indexBytes)

	rowOffset := headerSize
	indexOffset := pageSize - indexBytes - 1
	for _, row := range page.rows {
		rowBytes, _ := row.Marshal()
		byteOrder.PutUint32(data[rowOffset:], uint32(len(rowBytes)))
		copy(data[rowOffset+4:], rowBytes)

		copy(data[indexOffset:], row.Key[:])
		byteOrder.PutUint32(data[indexOffset+16:], uint32(rowOffset))

		rowOffset += uint64(len(rowBytes)) + 4
		indexOffset += indexEntrySize
	}

	data[rowOffset] = endOfPage

	n, err := writer.Write(data[:])
	return int64(n), err
}

func (page *dataPage) ReadFrom(reader io.Reader) (int64, error) {
	data := [pageSize]byte{}

	n, err := reader.Read(data[:])
	if err != nil {
		return int64(n), err
	}

	rowCount := byteOrder.Uint32(data[32:36])
	indexBytes := byteOrder.Uint64(data[44:52])

	rows := make([]Row, 0, int(rowCount))
	indexOffset := pageSize - indexBytes - 1
	for i := uint32(0); i < rowCount; i++ {
		key := make([]byte, 16)
		copy(key, data[indexOffset:indexOffset+16])
		rowOffset := byteOrder.Uint32(data[indexOffset+16 : indexOffset+20])

		length := byteOrder.Uint32(data[rowOffset:])
		row := Row{
			Key: key,
		}
		_ = row.Unmarshal(data[rowOffset+4 : rowOffset+4+length])

		rows = append(rows, row)
		indexOffset += indexEntrySize
	}

	page.rows = rows

	return int64(n), nil
}
