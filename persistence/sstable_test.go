package persistence

import (
	"github.com/google/uuid"
)

type trueIterator struct{}

func (it *trueIterator) Next() Row {
	id := uuid.New()
	return Row{
		Key:   id[:],
		Value: []byte("Hello World! SSTable are amazing and work well for Key-Value-Database"),
	}
}

func (it *trueIterator) HasNext() bool {
	return true
}
