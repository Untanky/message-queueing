package queueing_test

import (
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

	index.Set(id, location)

	val, ok := index.Get(id)
	if !ok {
		t.Fatalf("err: expected true, got %v", ok)
	}
	if val != location {
		t.Fatalf("val: expected %d, got %d", location, val)
	}

	val, ok = index.Delete(id)
	if !ok {
		t.Fatalf("err: expected true, got %v", ok)
	}
	if val != location {
		t.Fatalf("val: expected %d, got %d", location, val)
	}

	val, ok = index.Get(id)
	if !ok {
		t.Fatalf("err: expected true, got %v", ok)
	}
	if val != 0 {
		t.Fatalf("val: expected: %d, got %d", 0, val)
	}
}
