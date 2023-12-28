package queueing

import "cmp"

type treeNode[Key cmp.Ordered, Value any] struct {
	key   Key
	value Value
	left  *treeNode[Key, Value]
	right *treeNode[Key, Value]
}

func (node *treeNode[Key, Value]) depth() int {
	depth := 1

	if node.left != nil {
		depth = max(depth, node.left.depth()+1)
	}

	if node.right != nil {
		depth = max(depth, node.right.depth()+1)
	}

	return depth
}

func (node *treeNode[Key, Value]) balance() int {
	leftDepth, rightDepth := 0, 0

	if node.left != nil {
		leftDepth = node.left.depth()
	}
	if node.right != nil {
		rightDepth = node.right.depth()
	}

	return rightDepth - leftDepth
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

func (tree *AVLTree[Key, Value]) Depth() int {
	if tree.root == nil {
		return 0
	}

	return tree.root.depth()
}

func (tree *AVLTree[Key, Value]) Set(key Key, value Value) {
	panic("not implemented")
}
