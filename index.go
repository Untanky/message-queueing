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

type sliceIndex struct {
	data []indexTuple
}

func NewPrimaryIndex() Index {
	return &sliceIndex{
		data: make([]indexTuple, 0, 16),
	}
}

func (index *sliceIndex) Get(id MessageId) (MessageLocation, error) {
	target := indexTuple{id: id}
	k, ok := slices.BinarySearchFunc(index.data, target, compareTuples)
	if !ok {
		return 0, errors.New("not found")
	}

	return index.data[k].location, nil
}

func (index *sliceIndex) Set(id MessageId, location MessageLocation) error {
	tuple := indexTuple{
		id:       id,
		location: location,
	}

	unsorted := append(index.data, tuple)
	slices.SortFunc(unsorted, compareTuples)
	index.data = unsorted

	return nil
}

func (index *sliceIndex) Delete(id MessageId) (MessageLocation, error) {
	target := indexTuple{id: id}
	k, ok := slices.BinarySearchFunc(index.data, target, compareTuples)
	if !ok {
		return 0, errors.New("not found")
	}

	loc := index.data[k].location
	index.data = append(index.data[:k], index.data[k+1:]...)

	return loc, nil
}
