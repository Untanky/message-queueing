package queueing

import (
	"errors"
	"slices"
)

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
	data []indexTuple
}

func NewNaiveIndex() Index {
	return &naiveIndex{
		data: make([]indexTuple, 0, 16),
	}
}

func (index *naiveIndex) Get(id MessageId) (MessageLocation, error) {
	target := indexTuple{id: id}
	i, ok := slices.BinarySearchFunc(index.data, target, compareTuples)
	if !ok {
		return 0, errors.New("not found")
	}

	return index.data[i].location, nil
}

func (index *naiveIndex) Set(id MessageId, location MessageLocation) error {
	tuple := indexTuple{
		id:       id,
		location: location,
	}

	unsorted := append(index.data, tuple)
	slices.SortFunc(unsorted, compareTuples)
	index.data = unsorted

	return nil
}

func (index *naiveIndex) Delete(id MessageId) (MessageLocation, error) {
	target := indexTuple{id: id}
	i, ok := slices.BinarySearchFunc(index.data, target, compareTuples)
	if !ok {
		return 0, errors.New("not found")
	}

	loc := index.data[i].location
	index.data = append(index.data[:i], index.data[i+1:]...)

	return loc, nil
}