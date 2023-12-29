package queueing

import (
	"cmp"
)

type treeNode[Key cmp.Ordered, Value any] struct {
	depthCache int
	left       *treeNode[Key, Value]
	right      *treeNode[Key, Value]
	key        Key
	value      Value
}

func (node *treeNode[Key, Value]) depth() int {
	return node.depthCache
}

func (node *treeNode[Key, Value]) calculateDepth() int {
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
	node := &treeNode[Key, Value]{
		key:        key,
		value:      value,
		depthCache: 1,
	}

	if tree.root == nil {
		tree.root = node
		return
	}

	tree.root = insert(tree.root, node)
}

func insert[Key cmp.Ordered, Value any](root *treeNode[Key, Value], node *treeNode[Key, Value]) *treeNode[Key, Value] {
	current := root

	stack := make([]*treeNode[Key, Value], 0, root.depth())
	for current != nil {
		stack = append(stack, current)

		if node.key < current.key {
			if current.left == nil {
				current.left = node
				return rebalance(stack)
			}
			current = current.left
		} else if current.key == node.key {
			current.value = node.value
			return root
		} else {
			if current.right == nil {
				current.right = node
				return rebalance(stack)
			}
			current = current.right
		}
	}

	return root
}

func rebalance[Key cmp.Ordered, Value any](stack []*treeNode[Key, Value]) *treeNode[Key, Value] {
	for i := len(stack) - 1; i >= 0; i-- {
		node := stack[i]
		var updateParent *treeNode[Key, Value]
		node.depthCache = node.calculateDepth()
		balance := node.balance()
		if balance <= -2 {
			if node.left.balance() <= 0 {
				updateParent = rotateRight(node)
			} else {
				updateParent = rotateLeftRight(node)
			}
		} else if balance >= 2 {
			if node.right.balance() >= 0 {
				updateParent = rotateLeft(node)
			} else {
				updateParent = rotateRightLeft(node)
			}
		}

		if updateParent != nil {
			if i == 0 {
				return updateParent
			}
			if stack[i-1].left == node {
				stack[i-1].left = updateParent
			} else {
				stack[i-1].right = updateParent
			}
		}
	}

	return stack[0]
}

func rotateLeft[Key cmp.Ordered, Value any](node *treeNode[Key, Value]) *treeNode[Key, Value] {
	root := node.right
	node.right = root.left
	root.left = node
	return root
}

func rotateRight[Key cmp.Ordered, Value any](node *treeNode[Key, Value]) *treeNode[Key, Value] {
	root := node.left
	node.left = root.right
	root.right = node
	return root
}

func rotateRightLeft[Key cmp.Ordered, Value any](node *treeNode[Key, Value]) *treeNode[Key, Value] {
	node.right = rotateRight(node.right)
	return rotateLeft(node)
}

func rotateLeftRight[Key cmp.Ordered, Value any](node *treeNode[Key, Value]) *treeNode[Key, Value] {
	node.left = rotateLeft(node.left)
	return rotateRight(node)
}
