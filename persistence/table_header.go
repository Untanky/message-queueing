package persistence

import (
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
	pages                 uint32
	spans                 []pageSpanWithOffset
}

func (header *tableHeader) addPage(span pageSpanWithOffset) {
	header.pages++
	header.spans = append(header.spans, span)
}
