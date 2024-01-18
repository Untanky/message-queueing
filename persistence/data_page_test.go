package persistence

import (
	"github.com/google/uuid"
	"testing"
)

type dataPageTest struct{}

func (dataPageTest) newPage() *dataPage {
	return newDataPage()
}

func (dataPageTest) setupPage(page *dataPage) {
	ok := true
	for ok {
		id := uuid.New()
		ok = page.addRow(Row{
			Key:   id[:],
			Value: []byte("Hello World! SSTable are amazing and work well for Key-Value-Database"),
		})
	}
}

func (dataPageTest) compare(t *testing.T, a, b *dataPage) {
	//if len(a.rows) != len(b.rows) {
	//	t.Errorf("len(rows): tableA %v; tableB: %v", len(a.rows), len(b.rows))
	//}
}
