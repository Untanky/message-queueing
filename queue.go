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
	location MessageLocation
}

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
	val := h.data[l-1]
	h.data = h.data[:l-1]

	fmt.Println(h.data, val)
	return val
}

type heapQueue struct {
	lock sync.Locker

	heap *heapWrapper
}

func NewHeapQueue() *heapQueue {
	return &heapQueue{
		lock: &sync.Mutex{},
		heap: &heapWrapper{
			data: make([]queueTuple, 0, 16),
		},
	}
}

func (queue *heapQueue) Enqueue(timeout time.Time, location MessageLocation) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	heap.Push(queue.heap, queueTuple{timeout: timeout, location: location})
}

func (queue *heapQueue) Dequeue() (MessageLocation, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	return queue.dequeue()
}

func (queue *heapQueue) dequeue() (MessageLocation, error) {
	value := heap.Pop(queue.heap)
	tuple := value.(queueTuple)
	fmt.Println(tuple.timeout, time.Now())
	if tuple.timeout.After(time.Now()) {
		heap.Push(queue.heap, value)
		return 0, errors.New("next message not ready yet")
	}

	return tuple.location, nil
}

func (queue *heapQueue) DequeueMultiple(location []MessageLocation) (int, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	l := len(location)

	for i := 0; i < l; i++ {
		loc, err := queue.dequeue()
		if err != nil {
			return i, fmt.Errorf("maximum number of messages retrieved, because %w", err)
		}
		location[i] = loc
	}

	return l, nil
}
