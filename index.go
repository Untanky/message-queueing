package queueing

import (
	"errors"
	"slices"
	"sync"
)

var NotFoundError = errors.New("not found")

type indexTuple struct {
	id       MessageId
	location MessageLocation
}

func compareTuples(a, b indexTuple) int {
	for i := 0; i < 16; i++ {
		if a.id[i] < b.id[i] {
			return -1
		} else if a.id[i] > b.id[i] {
			return 1
		}
	}

	return 0
}

type naiveIndex struct {
	lock sync.Locker
	data []indexTuple
}

func NewNaiveIndex() Index[MessageId, MessageLocation] {
	return &naiveIndex{
		lock: &sync.Mutex{},
		data: make([]indexTuple, 0, 16),
	}
}

func (index *naiveIndex) Get(id MessageId) (MessageLocation, bool) {
	index.lock.Lock()
	defer index.lock.Unlock()

	target := indexTuple{id: id}
	i, ok := slices.BinarySearchFunc(index.data, target, compareTuples)
	if !ok {
		return 0, false
	}

	return index.data[i].location, true
}

func (index *naiveIndex) Set(id MessageId, location MessageLocation) {
	index.lock.Lock()
	defer index.lock.Unlock()

	tuple := indexTuple{
		id:       id,
		location: location,
	}

	unsorted := append(index.data, tuple)
	slices.SortFunc(unsorted, compareTuples)
	index.data = unsorted

	return
}

func (index *naiveIndex) Delete(id MessageId) (MessageLocation, bool) {
	index.lock.Lock()
	defer index.lock.Unlock()

	target := indexTuple{id: id}
	i, ok := slices.BinarySearchFunc(index.data, target, compareTuples)
	if !ok {
		return 0, false
	}

	loc := index.data[i].location
	index.data = append(index.data[:i], index.data[i+1:]...)

	return loc, true
}
