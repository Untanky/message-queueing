package queueing

import "cmp"

type treeNode[Key cmp.Ordered, Value any] struct {
	key   Key
	value Value
	left  *treeNode[Key, Value]
	right *treeNode[Key, Value]
}

type AVLTree[Key cmp.Ordered, Value any] struct {
	root *treeNode[Key, Value]
}

func (tree *AVLTree[Key, Value]) Get(key Key) (Value, bool) {
	current := tree.root

	for current != nil {
		if key < current.key {
			current = current.left
		} else if current.key == key {
			return current.value, true
		} else {
			current = current.right
		}
	}

	var noop Value
	return noop, false
}
