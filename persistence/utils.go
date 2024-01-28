package persistence

func compareBytes(a, b []byte) int {
	l := min(len(a), len(b))
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}

	return len(a) - len(b)
}

type pageSpan struct {
	startKey []byte
	endKey   []byte
}

func (span pageSpan) containsKey(key []byte) bool {
	startKey := compareBytes(key, span.startKey)
	endKey := compareBytes(key, span.endKey)
	return startKey >= 0 && endKey <= 0
}
