package main

import (
	"fmt"
	"github.com/google/uuid"
	queueing "message-queueing"
)

func main() {
	index := queueing.NewPrimaryIndex()

	a := queueing.MessageId(uuid.New())
	b := queueing.MessageId(uuid.New())
	c := queueing.MessageId(uuid.New())

	index.Set(a, 0)
	index.Set(b, 1234)
	index.Set(c, 2468)

	fmt.Println(index.Get(a))
	fmt.Println(index.Get(b))
	fmt.Println(index.Get(c))
	fmt.Println(index.Get(queueing.MessageId(uuid.New())))

	fmt.Println(index)

	index.Delete(a)

	fmt.Println(index)
}
