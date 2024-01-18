package persistence

import (
	"github.com/google/uuid"
	"testing"
)

type tableHeaderTest struct{}

func (tableHeaderTest) newPage() *tableHeader {
	return newTableHeader()
}

func (tableHeaderTest) setupPage(header *tableHeader) {
	for i := uint64(1); i < 17; i++ {
		id1 := uuid.New()
		id2 := uuid.New()

		header.addPage(pageSpanWithOffset{
			pageSpan: pageSpan{
				startKey: id1[:],
				endKey:   id2[:],
			},
			offset: i * pageSize,
		})
	}
}

func (tableHeaderTest) compare(t *testing.T, headerA, headerB *tableHeader) {
	if headerA.tableVersion != headerB.tableVersion {
		t.Errorf("tableVersion: tableA %v; tableB: %v", headerA.tableVersion, headerB.tableVersion)
	}

	if headerA.tableID != headerB.tableID {
		t.Errorf("tableID: tableA %v; tableB: %v", headerA.tableID, headerB.tableID)
	}

	if headerA.createdAt.Sub(headerB.createdAt).Milliseconds() > 1 {
		t.Errorf("createdAt: tableA %v; tableB: %v", headerA.createdAt, headerB.createdAt)
	}

	if headerA.compactionInformation != headerB.compactionInformation {
		t.Errorf("compactionInformation: tableA %v; tableB: %v", headerA.compactionInformation, headerB.compactionInformation)
	}

	if len(headerA.spans) != len(headerB.spans) {
		t.Errorf("len(spans): tableA %v; tableB: %v", len(headerA.spans), len(headerB.spans))
	}
}
