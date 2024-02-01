package persistence

import (
	"io"
	"sort"
)

type ordered interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint |
		~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type SecondaryIndex[Value any, IndexKey ordered] interface {
	Add(key IndexKey, value Value)
	Get(key IndexKey) (Value, bool)
	Range(lower, upper IndexKey) []Value
}

type secondaryIndexEntry[Value any, IndexKey ordered] struct {
	key   IndexKey
	value Value
}

type inMemorySecondaryIndex[IndexKey ordered] struct {
	entries []*secondaryIndexEntry[[]byte, IndexKey]
}

func (secondaryIndex *inMemorySecondaryIndex[IndexKey]) Add(key IndexKey, value []byte) {
	index := sort.Search(len(secondaryIndex.entries), func(i int) bool {
		return secondaryIndex.entries[i].key >= key
	})

	entry := &secondaryIndexEntry[[]byte, IndexKey]{
		key:   key,
		value: value,
	}

	secondaryIndex.entries = append(secondaryIndex.entries, nil)

	copy(secondaryIndex.entries[index+1:], secondaryIndex.entries[index:])
	secondaryIndex.entries[index] = entry
}

func (index *inMemorySecondaryIndex[IndexKey]) Get(key IndexKey) ([]byte, bool) {
	entries := index.entries

	i, found := sort.Find(len(entries), func(i int) int {
		return int(key - index.entries[i].key)
	})
	if !found {
		return nil, false
	}

	return index.entries[i].value, true
}

func (index *inMemorySecondaryIndex[IndexKey]) Range(lower, upper IndexKey) [][]byte {
	lowerIndex := sort.Search(len(index.entries), func(i int) bool {
		return index.entries[i].key >= lower
	})

	upperIndex := sort.Search(len(index.entries), func(i int) bool {
		return index.entries[i].key >= upper
	})
	if upperIndex == len(index.entries) || index.entries[upperIndex].key > upper {
		upperIndex -= 1
	}

	result := make([][]byte, 0, max(1, upperIndex-lowerIndex))
	for i := lowerIndex; i <= upperIndex; i++ {
		result = append(result, index.entries[i].value)
	}

	return result
}

func (index *inMemorySecondaryIndex[IndexKey]) WriteTo(writer io.Writer) (int64, error) {
	entriesPerPage := (pageSize - headerSize) / 24
	numberOfPages := uint8((uint64(len(index.entries)) / entriesPerPage) + 1)
	offset := int64(0)

	for i := uint8(0); i < numberOfPages; i++ {
		flag := indexFlag(0)
		if i == 0 {
			flag &= firstPage
		} else if i == numberOfPages-1 {
			flag &= lastPage
		}

		maxIndex := min(uint64(i+1)*entriesPerPage, uint64(len(index.entries)))

		page := &secondaryIndexPage[IndexKey]{
			flag:                 flag,
			keySize:              8,
			indexPageNumber:      i,
			numberOfIndexPages:   numberOfPages,
			totalNumberOfEntries: uint32(len(index.entries)),
			entries:              index.entries[uint64(i)*entriesPerPage : maxIndex],
		}
		n, err := page.WriteTo(writer)
		if err != nil {
			return offset + n, err
		}
		offset += n
	}

	return offset, nil
}

func (index *inMemorySecondaryIndex[IndexKey]) ReadFrom(reader io.Reader) (int64, error) {
	firstPage := &secondaryIndexPage[IndexKey]{
		entries: make([]*secondaryIndexEntry[[]byte, IndexKey], 680),
	}

	offset, err := firstPage.ReadFrom(reader)
	if err != nil {
		return offset, err
	}

	index.entries = make([]*secondaryIndexEntry[[]byte, IndexKey], firstPage.totalNumberOfEntries)
	copy(index.entries, firstPage.entries)

	for i := uint64(1); i < uint64(firstPage.numberOfIndexPages); i++ {
		maxIndex := min((i+1)*680, uint64(len(index.entries)))
		entries := index.entries[i*680 : maxIndex]
		page := &secondaryIndexPage[IndexKey]{
			entries: entries,
		}

		n, err := page.ReadFrom(reader)
		if err != nil {
			return offset + n, err
		}
		offset += n
	}

	return offset, err
}

type indexFlag uint8

const (
	firstPage = indexFlag(0x01)
	lastPage  = indexFlag(0x02)
)

type secondaryIndexPage[IndexKey ordered] struct {
	flag                 indexFlag
	keySize              uint8
	indexPageNumber      uint8
	numberOfIndexPages   uint8
	pageID               uint64
	totalNumberOfEntries uint32

	entries []*secondaryIndexEntry[[]byte, IndexKey]
}

func (page *secondaryIndexPage[IndexKey]) WriteTo(writer io.Writer) (int64, error) {
	data := make([]byte, pageSize)

	data[0] = byte(page.flag)
	data[1] = page.keySize
	data[2] = page.indexPageNumber
	data[3] = page.numberOfIndexPages
	byteOrder.PutUint64(data[4:12], page.pageID)
	byteOrder.PutUint32(data[12:16], page.totalNumberOfEntries)

	byteOrder.PutUint64(data[16:24], uint64(page.entries[0].key))
	byteOrder.PutUint64(data[24:32], uint64(page.entries[len(page.entries)-1].key))
	byteOrder.PutUint32(data[32:36], uint32(len(page.entries)))

	offset := 64
	for _, entry := range page.entries {
		byteOrder.PutUint64(data[offset:], uint64(entry.key))
		copy(data[offset+8:], entry.value)
		offset += len(entry.value) + 8
	}

	n, err := writer.Write(data)
	return int64(n), err
}

func (page *secondaryIndexPage[IndexKey]) ReadFrom(reader io.Reader) (int64, error) {
	data := make([]byte, pageSize)
	n, err := reader.Read(data)

	page.flag = indexFlag(data[0])
	page.keySize = data[1]
	page.indexPageNumber = data[2]
	page.numberOfIndexPages = data[3]
	page.pageID = byteOrder.Uint64(data[4:12])
	page.totalNumberOfEntries = byteOrder.Uint32(data[12:16])

	entryCount := byteOrder.Uint32(data[32:36])

	offset := 64
	page.entries = page.entries[:entryCount]
	for i, _ := range page.entries {
		key := IndexKey(byteOrder.Uint64(data[offset:]))
		value := make([]byte, 16)
		copy(value, data[offset+8:])
		page.entries[i] = &secondaryIndexEntry[[]byte, IndexKey]{
			key:   key,
			value: value,
		}
		offset += 24
	}

	return int64(n), err
}
