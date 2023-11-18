package queueing_test

import (
	"errors"
	"github.com/google/uuid"
	queueing "message-queueing"
	"testing"
)

func TestNaiveIndex_BasicFlow(t *testing.T) {
	index := queueing.NewNaiveIndex()

	testIndexBasicFlow(t, index)
}

func testIndexBasicFlow(t *testing.T, index queueing.Index[uuid.UUID, int]) {
	id := queueing.MessageId(uuid.MustParse("3fafc417-4438-4c5d-91c1-19549128b382"))
	location := queueing.MessageLocation(4)

	err := index.Set(id, location)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}

	val, err := index.Get(id)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}
	if val != location {
		t.Fatalf("val: expected %d, got %d", location, val)
	}

	val, err = index.Delete(id)
	if err != nil {
		t.Fatalf("err: expected nil, got %v", err)
	}
	if val != location {
		t.Fatalf("val: expected %d, got %d", location, val)
	}

	val, err = index.Get(id)
	if !errors.Is(err, queueing.NotFoundError) {
		t.Fatalf("err: expected %v, got %v", queueing.NotFoundError, err)
	}
	if val != 0 {
		t.Fatalf("val: expected: %d, got %d", 0, val)
	}
}
