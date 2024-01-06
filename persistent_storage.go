package queueing

import (
	"cmp"
	"context"
	"errors"
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
	tables []SSTable[Key, Value]
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

type SSTable[Key cmp.Ordered, Value any] interface {
	Get(context.Context, Key) (Value, error)
	Compact(context.Context, SSTable[Key, Value]) (SSTable[Key, Value], error)
}

func SSTableFrom[Key cmp.Ordered, Value any](iterator Iterator[Value]) SSTable[Key, Value] {
	// TODO: implement
	panic("not implemented")
}

type StorageEngine[Key cmp.Ordered, Value any] struct {
	memtable     MemTable[Key, Value]
	tablehandler TableHandler[Key, Value]
}

func (engine *StorageEngine[Key, Value]) Get(ctx context.Context, key Key) (Value, error) {
	value, ok := engine.memtable.Get(key)
	if ok {
		return value, nil
	}

	return engine.tablehandler.Get(ctx, key)
}

func (engine *StorageEngine[Key, Value]) Set(ctx context.Context, key Key, value Value) error {
	engine.memtable.Set(key, value)
	return nil
}

func (engine *StorageEngine[Key, Value]) Delete(ctx context.Context, key Key) (Value, error) {
	value, ok := engine.memtable.Delete(key)
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
