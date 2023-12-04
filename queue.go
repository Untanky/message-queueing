package queueing

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"
)

type queueTuple struct {
	timeout  time.Time
	location MessageId
}

var NextMessageUnavailableError = errors.New("message unavailable")

type messageUnavailableError string

func (err messageUnavailableError) Error() string {
	return string(err)
}

func (err messageUnavailableError) Unwrap() error {
	return NextMessageUnavailableError
}

const (
	NoMoreMessages      = messageUnavailableError("no more messages")
	NextMessageNotReady = messageUnavailableError("next message not ready yet")
)

type heapWrapper struct {
	data []queueTuple
}

func (h *heapWrapper) Len() int {
	return len(h.data)
}

func (h *heapWrapper) Less(i, j int) bool {
	return h.data[i].timeout.Before(h.data[j].timeout)
}

func (h *heapWrapper) Swap(i, j int) {
	if i < 0 || j < 0 {
		return
	}
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *heapWrapper) Push(x any) {
	tuple, ok := x.(queueTuple)
	if !ok {
		panic("cannot add non queueTuple to heapWrapper")
	}

	h.data = append(h.data, tuple)

}

func (h *heapWrapper) Pop() any {
	l := len(h.data)
	if l == 0 {
		return nil
	}

	val := h.data[l-1]
	h.data = h.data[:l-1]
	return val
}

type timeoutQueue struct {
	lock sync.Locker

	heap *heapWrapper
}

func NewTimeoutQueue() *timeoutQueue {
	return &timeoutQueue{
		lock: &sync.Mutex{},
		heap: &heapWrapper{
			data: make([]queueTuple, 0, 16),
		},
	}
}

func (queue *timeoutQueue) Enqueue(timeout time.Time, location MessageId) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	heap.Push(queue.heap, queueTuple{timeout: timeout, location: location})
}

func (queue *timeoutQueue) Dequeue(before time.Time) (MessageId, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	return queue.dequeue(before)
}

func (queue *timeoutQueue) dequeue(before time.Time) (MessageId, error) {
	value := heap.Pop(queue.heap)
	if value == nil {
		return MessageId{}, NoMoreMessages
	}

	tuple := value.(queueTuple)
	if tuple.timeout.After(before) {
		heap.Push(queue.heap, value)
		return MessageId{}, NextMessageNotReady
	}

	return tuple.location, nil
}

func (queue *timeoutQueue) DequeueMultiple(location []MessageId, before time.Time) (int, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	l := len(location)

	for i := 0; i < l; i++ {
		loc, err := queue.dequeue(before)
		if err != nil {
			return i, fmt.Errorf("maximum number of messages retrieved, because %w", err)
		}
		location[i] = loc
	}

	return l, nil
}
