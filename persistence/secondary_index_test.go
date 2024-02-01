package persistence

import (
	"testing"
)

func TestInMemorySecondaryIndex_Add(t *testing.T) {
	type testCase struct {
		name            string
		secondaryIndex  *inMemorySecondaryIndex[int]
		key             int
		bytes           []byte
		expectedEntries []*secondaryIndexEntry[[]byte, int]
	}

	createDefaultIndex := func() *inMemorySecondaryIndex[int] {
		return &inMemorySecondaryIndex[int]{
			entries: []*secondaryIndexEntry[[]byte, int]{
				{1, []byte("abc")},
				{3, []byte("def")},
				{5, []byte("ghi")},
				{7, []byte("jkl")},
				{9, []byte("mno")},
				{11, []byte("pqr")},
				{13, []byte("stu")},
			},
		}
	}

	cases := []testCase{
		{
			name:           "Add entry in the middle",
			secondaryIndex: createDefaultIndex(),
			key:            6,
			bytes:          []byte("xyz"),
			expectedEntries: []*secondaryIndexEntry[[]byte, int]{
				{1, []byte("abc")},
				{3, []byte("def")},
				{5, []byte("ghi")},
				{6, []byte("xyz")},
				{7, []byte("jkl")},
				{9, []byte("mno")},
				{11, []byte("pqr")},
				{13, []byte("stu")},
			},
		},
		{
			name:           "Add entry in the front",
			secondaryIndex: createDefaultIndex(),
			key:            0,
			bytes:          []byte("xyz"),
			expectedEntries: []*secondaryIndexEntry[[]byte, int]{
				{0, []byte("xyz")},
				{1, []byte("abc")},
				{3, []byte("def")},
				{5, []byte("ghi")},
				{7, []byte("jkl")},
				{9, []byte("mno")},
				{11, []byte("pqr")},
				{13, []byte("stu")},
			},
		},
		{
			name:           "Add entry at the end",
			secondaryIndex: createDefaultIndex(),
			key:            15,
			bytes:          []byte("xyz"),
			expectedEntries: []*secondaryIndexEntry[[]byte, int]{
				{1, []byte("abc")},
				{3, []byte("def")},
				{5, []byte("ghi")},
				{7, []byte("jkl")},
				{9, []byte("mno")},
				{11, []byte("pqr")},
				{13, []byte("stu")},
				{15, []byte("xyz")},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			c.secondaryIndex.Add(c.key, c.bytes)

			if len(c.secondaryIndex.entries) != len(c.expectedEntries) {
				tt.Errorf("len(bytes): got %v; expected %v", len(c.secondaryIndex.entries), len(c.expectedEntries))
			}

			for i := 0; i < len(c.secondaryIndex.entries); i++ {
				if c.secondaryIndex.entries[i].key != c.expectedEntries[i].key {
					tt.Errorf("key[%d]: got %v; expected %v", i, c.secondaryIndex.entries[i].key, c.expectedEntries[i].key)
				}
				if compareBytes(c.secondaryIndex.entries[i].value, c.expectedEntries[i].value) != 0 {
					tt.Errorf("bytes[%d]: got %v; expected %v", i, c.secondaryIndex.entries[i].value, c.expectedEntries[i].value)
				}
			}
		})
	}
}

func TestInMemorySecondaryIndex_Get(t *testing.T) {
	type testCase struct {
		name           string
		secondaryIndex SecondaryIndex[[]byte, int]
		key            int
		wantBytes      []byte
		wantOk         bool
	}

	defaultIndex := &inMemorySecondaryIndex[int]{
		entries: []*secondaryIndexEntry[[]byte, int]{
			{1, []byte("abc")},
			{3, []byte("def")},
			{5, []byte("ghi")},
			{7, []byte("jkl")},
			{9, []byte("mno")},
			{11, []byte("pqr")},
			{13, []byte("stu")},
		},
	}

	cases := []testCase{
		{
			name:           "Find element in middle of entries",
			secondaryIndex: defaultIndex,
			key:            7,
			wantBytes:      []byte("jkl"),
			wantOk:         true,
		},
		{
			name:           "Find element at the start of entries",
			secondaryIndex: defaultIndex,
			key:            1,
			wantBytes:      []byte("abc"),
			wantOk:         true,
		},
		{
			name:           "Find element at the end of entries",
			secondaryIndex: defaultIndex,
			key:            13,
			wantBytes:      []byte("stu"),
			wantOk:         true,
		},
		{
			name:           "Do NOT find element in middle of entries",
			secondaryIndex: defaultIndex,
			key:            2,
			wantBytes:      nil,
			wantOk:         false,
		},
		{
			name:           "Do NOT find element before start of entries",
			secondaryIndex: defaultIndex,
			key:            0,
			wantBytes:      nil,
			wantOk:         false,
		},
		{
			name:           "Do NOT find element after end of entries",
			secondaryIndex: defaultIndex,
			key:            15,
			wantBytes:      nil,
			wantOk:         false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			bytes, ok := c.secondaryIndex.Get(c.key)

			if compareBytes(bytes, c.wantBytes) != 0 {
				tt.Errorf("bytes: got %v; expected %v", bytes, c.wantBytes)
			}
			if ok != c.wantOk {
				tt.Errorf("ok: got %v; expected %v", ok, c.wantOk)
			}
		})
	}
}

func TestInMemorySecondaryIndex_Range(t *testing.T) {
	type testCase struct {
		name           string
		secondaryIndex SecondaryIndex[[]byte, int]
		lower          int
		upper          int
		wantBytes      [][]byte
	}

	defaultIndex := &inMemorySecondaryIndex[int]{
		entries: []*secondaryIndexEntry[[]byte, int]{
			{1, []byte("abc")},
			{3, []byte("def")},
			{5, []byte("ghi")},
			{7, []byte("jkl")},
			{9, []byte("mno")},
			{11, []byte("pqr")},
			{13, []byte("stu")},
		},
	}

	cases := []testCase{
		{
			name:           "Find range completely in entries",
			secondaryIndex: defaultIndex,
			lower:          3,
			upper:          11,
			wantBytes: [][]byte{
				[]byte("def"),
				[]byte("ghi"),
				[]byte("jkl"),
				[]byte("mno"),
				[]byte("pqr"),
			},
		},
		{
			name:           "Find range completely in entries",
			secondaryIndex: defaultIndex,
			lower:          2,
			upper:          12,
			wantBytes: [][]byte{
				[]byte("def"),
				[]byte("ghi"),
				[]byte("jkl"),
				[]byte("mno"),
				[]byte("pqr"),
			},
		},
		{
			name:           "Find range completely matching entries",
			secondaryIndex: defaultIndex,
			lower:          1,
			upper:          13,
			wantBytes: [][]byte{
				[]byte("abc"),
				[]byte("def"),
				[]byte("ghi"),
				[]byte("jkl"),
				[]byte("mno"),
				[]byte("pqr"),
				[]byte("stu"),
			},
		},
		{
			name:           "Find range beginning outside entries",
			secondaryIndex: defaultIndex,
			lower:          0,
			upper:          13,
			wantBytes: [][]byte{
				[]byte("abc"),
				[]byte("def"),
				[]byte("ghi"),
				[]byte("jkl"),
				[]byte("mno"),
				[]byte("pqr"),
				[]byte("stu"),
			},
		},
		{
			name:           "Find range ending outside entries",
			secondaryIndex: defaultIndex,
			lower:          1,
			upper:          14,
			wantBytes: [][]byte{
				[]byte("abc"),
				[]byte("def"),
				[]byte("ghi"),
				[]byte("jkl"),
				[]byte("mno"),
				[]byte("pqr"),
				[]byte("stu"),
			},
		},
		{
			name:           "Find range beginning and ending outside entries",
			secondaryIndex: defaultIndex,
			lower:          1,
			upper:          14,
			wantBytes: [][]byte{
				[]byte("abc"),
				[]byte("def"),
				[]byte("ghi"),
				[]byte("jkl"),
				[]byte("mno"),
				[]byte("pqr"),
				[]byte("stu"),
			},
		},
		{
			name:           "Find single entries",
			secondaryIndex: defaultIndex,
			lower:          3,
			upper:          3,
			wantBytes: [][]byte{
				[]byte("def"),
			},
		},
		{
			name:           "Find single entries",
			secondaryIndex: defaultIndex,
			lower:          2,
			upper:          4,
			wantBytes: [][]byte{
				[]byte("def"),
			},
		},
		{
			name:           "Find range beginning and ending smaller than all entries",
			secondaryIndex: defaultIndex,
			lower:          -1,
			upper:          0,
			wantBytes:      [][]byte{},
		},
		{
			name:           "Find range beginning smaller than all entries",
			secondaryIndex: defaultIndex,
			lower:          -1,
			upper:          1,
			wantBytes: [][]byte{
				[]byte("abc"),
			},
		},
		{
			name:           "Find range ending greater than all entries",
			secondaryIndex: defaultIndex,
			lower:          13,
			upper:          16,
			wantBytes: [][]byte{
				[]byte("stu"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			results := c.secondaryIndex.Range(c.lower, c.upper)

			if len(results) != len(c.wantBytes) {
				tt.Errorf("len(results): got %v; expected %v", len(results), len(c.wantBytes))
				return
			}

			for i := 0; i < len(results); i++ {
				if compareBytes(results[i], c.wantBytes[i]) != 0 {
					tt.Errorf("bytes[%d]: got %v; expected %v", i, results[i], c.wantBytes[i])
				}
			}
		})
	}
}

type secondaryIndexTest struct{}

func (secondaryIndexTest) newPage() *inMemorySecondaryIndex[int] {
	entries := make([]*secondaryIndexEntry[[]byte, int], 0, 2000)
	for i := 0; i < cap(entries); i++ {
		entries = append(entries, &secondaryIndexEntry[[]byte, int]{
			key:   i,
			value: []byte("abcdefghabcdefgh"),
		})
	}

	return &inMemorySecondaryIndex[int]{
		entries: entries,
	}
}

func (secondaryIndexTest) setupPage(page *inMemorySecondaryIndex[int]) {
}

func (secondaryIndexTest) compare(t *testing.T, a, b *inMemorySecondaryIndex[int]) {
	if len(a.entries) != len(b.entries) {
		t.Errorf("len(entries): tableA %v; tableB %v", len(a.entries), len(b.entries))
	}
}

type secondaryIndexPageTest struct{}

func (secondaryIndexPageTest) newPage() *secondaryIndexPage[int] {
	entries := make([]*secondaryIndexEntry[[]byte, int], 0, 680)
	for i := 0; i < cap(entries); i++ {
		entries = append(entries, &secondaryIndexEntry[[]byte, int]{
			key:   i,
			value: []byte("abcdefghabcdefgh"),
		})
	}

	return &secondaryIndexPage[int]{
		flag:                 firstPage & lastPage,
		keySize:              8,
		indexPageNumber:      0,
		numberOfIndexPages:   0,
		pageID:               10,
		totalNumberOfEntries: 680,
		entries:              entries,
	}
}

func (secondaryIndexPageTest) setupPage(page *secondaryIndexPage[int]) {
}

func (secondaryIndexPageTest) compare(t *testing.T, a, b *secondaryIndexPage[int]) {
	if a.flag != b.flag {
		t.Errorf("flag: tableA %v; tableB %v", a.flag, b.flag)
	}

	if a.keySize != b.keySize {
		t.Errorf("keySize: tableA %v; tableB %v", a.keySize, b.keySize)
	}

	if a.indexPageNumber != b.indexPageNumber {
		t.Errorf("indexPageNumber: tableA %v; tableB %v", a.indexPageNumber, b.indexPageNumber)
	}

	if a.numberOfIndexPages != b.numberOfIndexPages {
		t.Errorf("numberOfIndexPages: tableA %v; tableB %v", a.numberOfIndexPages, b.numberOfIndexPages)
	}

	if a.pageID != b.pageID {
		t.Errorf("pageID: tableA %v; tableB %v", a.pageID, b.pageID)
	}

	if a.totalNumberOfEntries != b.totalNumberOfEntries {
		t.Errorf("totalNumberOfEntries: tableA %v; tableB %v", a.totalNumberOfEntries, b.totalNumberOfEntries)
	}

	if len(a.entries) != len(b.entries) {
		t.Errorf("len(entries): tableA %v; tableB %v", len(a.entries), len(b.entries))
	}
}
