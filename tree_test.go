package queueing

import (
	"math/rand"
	"testing"
)

func newTestTree() *AVLTree[int, string] {
	return &AVLTree[int, string]{
		root: &treeNode[int, string]{
			key:   7,
			value: "abc",
			left: &treeNode[int, string]{
				key:   3,
				value: "def",
				left: &treeNode[int, string]{
					key:   1,
					value: "ghi",
				},
				right: &treeNode[int, string]{
					key:   5,
					value: "jkl",
				},
			},
			right: &treeNode[int, string]{
				key:   11,
				value: "mno",
				left: &treeNode[int, string]{
					key:   9,
					value: "pqr",
				},
				right: &treeNode[int, string]{
					key:   13,
					value: "stu",
				},
			},
		},
	}
}

func TestAVLTree_Get(t *testing.T) {
	type kv struct {
		key   int
		value string
		ok    bool
	}

	cases := []struct {
		name string
		tree *AVLTree[int, string]
		kvs  []kv
	}{
		{
			name: "WithEmptyTree",
			tree: &AVLTree[int, string]{},
			kvs: []kv{
				{
					key:   -1,
					value: "",
					ok:    false,
				},
				{
					key:   0,
					value: "",
					ok:    false,
				},
				{
					key:   1,
					value: "",
					ok:    false,
				},
			},
		},
		{
			name: "WithBalancedFullTree",
			tree: newTestTree(),
			kvs: []kv{
				{
					key:   0,
					value: "",
					ok:    false,
				},
				{
					key:   1,
					value: "ghi",
					ok:    true,
				},
				{
					key:   2,
					value: "",
					ok:    false,
				},
				{
					key:   3,
					value: "def",
					ok:    true,
				},
				{
					key:   4,
					value: "",
					ok:    false,
				},
				{
					key:   5,
					value: "jkl",
					ok:    true,
				},
				{
					key:   6,
					value: "",
					ok:    false,
				},
				{
					key:   7,
					value: "abc",
					ok:    true,
				},
				{
					key:   8,
					value: "",
					ok:    false,
				},
				{
					key:   9,
					value: "pqr",
					ok:    true,
				},
				{
					key:   10,
					value: "",
					ok:    false,
				},
				{
					key:   11,
					value: "mno",
					ok:    true,
				},
				{
					key:   12,
					value: "",
					ok:    false,
				},
				{
					key:   13,
					value: "stu",
					ok:    true,
				},
				{
					key:   14,
					value: "",
					ok:    false,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(tt *testing.T) {
				for _, kv := range c.kvs {
					value, ok := c.tree.Get(kv.key)
					if value != kv.value {
						tt.Errorf("value: expected %s; got %s", kv.value, value)
					}
					if ok != kv.ok {
						tt.Errorf("ok: expected %t; got %t", kv.ok, ok)
					}
				}
			},
		)
	}
}

func TestAVLTree_Depth(t *testing.T) {
	cases := []struct {
		name        string
		tree        *AVLTree[int, string]
		wantedDepth int
	}{
		{
			name:        "WithEmptyTree",
			tree:        &AVLTree[int, string]{},
			wantedDepth: 0,
		},
		{
			name:        "WithBalancedFullTree",
			tree:        newTestTree(),
			wantedDepth: 3,
		},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(tt *testing.T) {
				depth := c.tree.Depth()
				if c.wantedDepth != depth {
					tt.Errorf("depth: expected %d; got %d", c.wantedDepth, depth)
				}
			},
		)
	}
}

func TestAVLTree_Set(t *testing.T) {
	t.Run(
		"EmptyTree",
		func(tt *testing.T) {
			tree := new(AVLTree[int, string])

			tree.Set(4, "foo")
			if tree.root.value != "foo" {
				t.Errorf("node is not at the expected location")
			}
		},
	)

	t.Run(
		"NoRotation",
		func(tt *testing.T) {
			tree := newTestTree()

			tree.Set(4, "foo")
			if tree.root.left.right.left.value != "foo" {
				t.Errorf("node is not at the expected location")
			}
		},
	)

	t.Run(
		"WithLeftRotation",
		func(tt *testing.T) {
			tree := newTestTree()

			tree.Set(16, "bar")
			if tree.root.right.right.right.value != "bar" {
				t.Errorf("node is not at the expected location")
			}

			tree.Set(17, "baz")
			if tree.root.right.right.value != "bar" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.right.right.left.value != "stu" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.right.right.right.value != "baz" {
				t.Errorf("node is not at the expected location")
			}
		},
	)

	t.Run(
		"WithRightRotation", func(t *testing.T) {
			tree := newTestTree()

			tree.Set(0, "bar")
			if tree.root.left.left.left.value != "bar" {
				t.Errorf("node is not at the expected location")
			}

			tree.Set(-1, "baz")
			if tree.root.left.left.value != "bar" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.left.left.right.value != "ghi" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.left.left.left.value != "baz" {
				t.Errorf("node is not at the expected location")
			}
		},
	)

	t.Run(
		"WithRightLeftRotation", func(tt *testing.T) {
			tree := &AVLTree[int, string]{
				root: &treeNode[int, string]{
					key:   5,
					value: "abc",
					right: &treeNode[int, string]{
						key:   8,
						value: "def",
					},
				},
			}

			tree.Set(7, "bar")
			if tree.root.left.value != "abc" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.right.value != "def" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.value != "bar" {
				t.Errorf("node is not at the expected location")
			}
		},
	)

	t.Run(
		"WithLeftRightRotation", func(tt *testing.T) {
			tree := &AVLTree[int, string]{
				root: &treeNode[int, string]{
					key:   5,
					value: "abc",
					left: &treeNode[int, string]{
						key:   2,
						value: "def",
					},
				},
			}

			tree.Set(3, "bar")
			if tree.root.right.value != "abc" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.left.value != "def" {
				t.Errorf("node is not at the expected location")
			}
			if tree.root.value != "bar" {
				t.Errorf("node is not at the expected location")
			}
		},
	)
}

func TestTreeNode_Balance(t *testing.T) {
	cases := []struct {
		name          string
		root          *treeNode[int, string]
		wantedBalance int
	}{
		{
			name:          "WithEmptyTree",
			root:          &treeNode[int, string]{},
			wantedBalance: 0,
		},
		{
			name: "WithLeftHangingTree",
			root: &treeNode[int, string]{
				left: &treeNode[int, string]{},
			},
			wantedBalance: -1,
		},
		{
			name: "WithRightHangingTree",
			root: &treeNode[int, string]{
				right: &treeNode[int, string]{},
			},
			wantedBalance: 1,
		},
		{
			name: "WithLeftImbalancedTree",
			root: &treeNode[int, string]{
				left: &treeNode[int, string]{
					left: &treeNode[int, string]{},
				},
			},
			wantedBalance: -2,
		},
		{
			name: "WithLeftImbalancedTree",
			root: &treeNode[int, string]{
				left: &treeNode[int, string]{
					right: &treeNode[int, string]{},
				},
			},
			wantedBalance: -2,
		},
		{
			name: "WithLeftHangingTree",
			root: &treeNode[int, string]{
				left: &treeNode[int, string]{
					right: &treeNode[int, string]{},
				},
				right: &treeNode[int, string]{},
			},
			wantedBalance: -1,
		},
		{
			name: "WithRightImbalancedTree",
			root: &treeNode[int, string]{
				right: &treeNode[int, string]{
					right: &treeNode[int, string]{},
				},
			},
			wantedBalance: 2,
		},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(tt *testing.T) {
				depth := c.root.balance()
				if c.wantedBalance != depth {
					tt.Errorf("depth: expected %d; got %d", c.wantedBalance, depth)
				}
			},
		)
	}
}

func BenchmarkAVLTree_Set(b *testing.B) {
	tree := new(AVLTree[int32, string])

	for i := 0; i < b.N; i++ {
		tree.Set(rand.Int31(), "abc")
	}
}
