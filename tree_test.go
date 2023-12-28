package queueing

import (
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
