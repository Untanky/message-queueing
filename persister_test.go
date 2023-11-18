package queueing_test

import (
	"errors"
	"io"
	queueing "message-queueing"
	"slices"
	"testing"
)

type sliceReadWriteSeeker struct {
	loc  int64
	data []byte
}

func (s *sliceReadWriteSeeker) Read(p []byte) (n int, err error) {
	copy(p, s.data[s.loc:])
	s.loc += int64(len(p))
	return len(p), nil
}

func (s *sliceReadWriteSeeker) Write(p []byte) (n int, err error) {
	s.data = append(s.data[:s.loc], p...)
	s.loc += int64(len(p))
	return len(p), nil
}

func (s *sliceReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var newLoc int64
	if whence == io.SeekStart {
		newLoc = offset
	} else if whence == io.SeekCurrent {
		newLoc = s.loc + offset
	} else if whence == io.SeekEnd {
		newLoc = int64(len(s.data)) + offset
	}

	if newLoc < 0 || newLoc > int64(len(s.data)) {
		return -1, errors.New("out of range")
	}

	s.loc = newLoc

	return newLoc, nil
}

func TestPersister_APIContract(t *testing.T) {
	h := new(sliceReadWriteSeeker)
	persister := queueing.NewPersister(h)

	testPersisterAPI(t, persister)
}

func testPersisterAPI(t *testing.T, persister queueing.Persister) {

	a := []byte("Hello World")
	b := []byte("Goodbye World")

	locA, err := persister.Write(a)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}
	if locA != 0 {
		t.Fatalf("locA: expected %d, got %d", 0, locA)
	}

	locB, err := persister.Write(b)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}
	if locB != 19 {
		t.Fatalf("locA: expected %d, got %d", 19, locB)
	}

	dataA, err := persister.Read(locA)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}
	if !slices.Equal(dataA, a) {
		t.Fatalf("dataA: expected %v, got %v", a, dataA)
	}

	dataB, err := persister.Read(locB)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}
	if !slices.Equal(dataB, b) {
		t.Fatalf("dataB: expected %v, got %v", b, dataB)
	}
}
