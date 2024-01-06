package queueing

import (
	"cmp"
	"context"
	"errors"
	"io"
)

type Store[Key cmp.Ordered, Value any] interface {
	Get(Key) (Value, bool)
	Set(Key, Value)
	Delete(Key) (Value, bool)
}

type Iterator[Value any] interface {
	Next() Value
	HasNext()
}

type MemTable[Key cmp.Ordered, Value any] interface {
	Store[Key, Value]
	Iterator() Iterator[Value]
	ClosedNotifier() <-chan bool
}

type TableHandler[Key cmp.Ordered, Value any] struct {
	tables []*SSTable[Key, Value]
}

func (handler *TableHandler[Key, Value]) CreateTable(ctx context.Context, iterator Iterator[Value]) error {
	newTable := SSTableFrom[Key, Value](iterator)
	handler.tables = append(handler.tables, newTable)
	return nil
}

func (handler *TableHandler[Key, Value]) Get(ctx context.Context, key Key) (Value, error) {
	if len(handler.tables) == 0 {
		var noop Value
		return noop, errors.New("value not found as tables does not exists")
	}

	for _, table := range handler.tables {
		value, err := table.Get(ctx, key)
		if err == nil {
			return value, nil
		}
	}

	var noop Value
	return noop, errors.New("value not found in tables")
}

func (handler *TableHandler[Key, Value]) Compact(ctx context.Context) {
	// TODO: implement
	panic("not implemented")
}

const ssTablePageSize = 64 * 1024

type pageSpan[Key cmp.Ordered] struct {
	startKey Key
	endKey   Key
	offset   int64
}

type SSTable[Key cmp.Ordered, Value any] struct {
	reader io.ReadSeekCloser
	spans  []pageSpan[Key]
}

func SSTableFrom[Key cmp.Ordered, Value any](iterator Iterator[Value]) *SSTable[Key, Value] {
	// TODO: implement
	panic("not implemented")
}

func (table *SSTable[Key, Value]) findPageSpan(key Key) (pageSpan[Key], bool) {
	spanSlice := table.spans

	for len(spanSlice) > 1 {
		i := len(spanSlice) / 2
		if spanSlice[i].startKey > key {
			spanSlice = spanSlice[:i]
		} else if spanSlice[i].endKey < key {
			spanSlice = spanSlice[i+1:]
		} else {
			return spanSlice[i], true
		}
	}

	if spanSlice[0].startKey <= key && spanSlice[0].endKey <= key {
		return spanSlice[0], true
	}

	return pageSpan[Key]{}, false
}

func (table *SSTable[Key, Value]) findInPage(pageSpan pageSpan[Key]) (Value, error) {
	// TODO: implement
	panic("not implemented")
}

func (table *SSTable[Key, Value]) Get(ctx context.Context, key Key) (Value, error) {
	span, ok := table.findPageSpan(key)
	if !ok {
		var noop Value
		return noop, errors.New("no span contains key")
	}

	return table.findInPage(span)
}

func (table *SSTable[Key, Value]) Compact(context.Context, SSTable[Key, Value]) (SSTable[Key, Value], error) {
	// TODO: implement
	panic("not implemented")
}

type StorageEngine[Key cmp.Ordered, Value any] struct {
	memTable     MemTable[Key, Value]
	tableHandler TableHandler[Key, Value]
}

func (engine *StorageEngine[Key, Value]) Get(ctx context.Context, key Key) (Value, error) {
	value, ok := engine.memTable.Get(key)
	if ok {
		return value, nil
	}

	return engine.tableHandler.Get(ctx, key)
}

func (engine *StorageEngine[Key, Value]) Set(ctx context.Context, key Key, value Value) error {
	engine.memTable.Set(key, value)
	return nil
}

func (engine *StorageEngine[Key, Value]) Delete(ctx context.Context, key Key) (Value, error) {
	value, ok := engine.memTable.Delete(key)
	if ok {
		return value, nil
	}

	return value, errors.New("cannot delete as entry does not exists")
}

func (engine *StorageEngine[Key, Value]) RunBackgroundTask() {
	// TODO: implement
	panic("not implemented")
}

func (engine *StorageEngine[Key, Value]) Close() error {
	// TODO: implement
	panic("not implemented")
}
