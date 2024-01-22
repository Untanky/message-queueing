package testutils

import (
	"errors"
	"io"
)

type SliceReadWriteSeeker struct {
	loc  int64
	Data []byte
}

func (s *SliceReadWriteSeeker) Read(p []byte) (n int, err error) {
	copy(p, s.Data[s.loc:])
	s.loc += int64(len(p))
	return len(p), nil
}

func (s *SliceReadWriteSeeker) Write(p []byte) (n int, err error) {
	offset := max(s.loc+int64(len(p)), int64(len(s.Data)))
	s.Data = append(s.Data[:s.loc], p...)[:offset]
	s.loc += int64(len(p))
	return len(p), nil
}

func (s *SliceReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var newLoc int64
	if whence == io.SeekStart {
		newLoc = offset
	} else if whence == io.SeekCurrent {
		newLoc = s.loc + offset
	} else if whence == io.SeekEnd {
		newLoc = int64(len(s.Data)) + offset
	}

	if newLoc < 0 || newLoc > int64(len(s.Data)) {
		return -1, errors.New("out of range")
	}

	s.loc = newLoc

	return newLoc, nil
}

func (s *SliceReadWriteSeeker) Close() error {
	return nil
}
