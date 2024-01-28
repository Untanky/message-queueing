package persistence

import "testing"

func TestCompareBytes(t *testing.T) {
	type testCase struct {
		name   string
		bytesA []byte
		bytesB []byte
		want   int
	}

	cases := []testCase{
		{
			name:   "Equal byte slices (length: 1)",
			bytesA: []byte{'a'},
			bytesB: []byte{'a'},
			want:   0,
		},
		{
			name:   "Equal byte slices (length: 10)",
			bytesA: []byte{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5'},
			bytesB: []byte{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5'},
			want:   0,
		},
		{
			name:   "bytesA less than bytesB (length: 1)",
			bytesA: []byte{'a'},
			bytesB: []byte{'b'},
			want:   -1,
		},
		{
			name:   "bytesA less than bytesB (in the middle)",
			bytesA: []byte{'a', 1, 'b'},
			bytesB: []byte{'a', 2, 'b'},
			want:   -1,
		},
		{
			name:   "bytesA less than bytesB (at the end)",
			bytesA: []byte{'a', 2, 'a'},
			bytesB: []byte{'a', 2, 'b'},
			want:   -1,
		},
		{
			name:   "bytesA greater than bytesB (length: 1)",
			bytesA: []byte{'b'},
			bytesB: []byte{'a'},
			want:   1,
		},
		{
			name:   "bytesA greater than bytesB (in the middle)",
			bytesA: []byte{'a', 2, 'b'},
			bytesB: []byte{'a', 1, 'b'},
			want:   1,
		},
		{
			name:   "bytesA greater than bytesB (at the end)",
			bytesA: []byte{'a', 2, 'b'},
			bytesB: []byte{'a', 2, 'a'},
			want:   1,
		},
		{
			name:   "two nils",
			bytesA: nil,
			bytesB: nil,
			want:   0,
		},
		{
			name:   "nil less than bytesB",
			bytesA: nil,
			bytesB: []byte{'a'},
			want:   -1,
		},
		{
			name:   "bytesA greater than nil",
			bytesA: []byte{'a'},
			bytesB: nil,
			want:   1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := compareBytes(c.bytesA, c.bytesB)

			if result != c.want {
				t.Errorf("result: want %v, got %v", c.want, result)
			}
		})
	}
}
