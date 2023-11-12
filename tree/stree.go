package tree

import "xcontainer"

type SearchTree[T xcontainer.Orderliness] struct {
	root *TreeNode[T]
}

func (r *SearchTree[T]) Find(data T) *TreeNode[T] {
	_, n := r.find(nil, r.root, data)
	return n
}

func (r *SearchTree[T]) find(pre, n *TreeNode[T], data T) (*TreeNode[T], *TreeNode[T]) {
	if xcontainer.IsNil[*TreeNode[T]](n) {
		return nil, nil
	}
	if n.Val() == data {
		return pre, n
	}
	if n.Val() > data {
		return r.find(n, n.left, data)
	} else {
		return r.find(n, n.right, data)
	}
}

func (r *SearchTree[T]) Put(data T) {
	if xcontainer.IsNil[*TreeNode[T]](r.root) {
		r.root = &TreeNode[T]{val: data}
		return
	}
	r.sortPut(r.root, data)
}

func (r *SearchTree[T]) Del(data T) {
	pre, n := r.find(nil, r.root, data)
	if xcontainer.IsNil[*TreeNode[T]](n) {
		return
	}
	if xcontainer.IsNil[*TreeNode[T]](pre) {
		pre = &TreeNode[T]{left: n}
		r.del(pre, n)
		r.root = pre.left
		return
	}
	r.del(pre, n)
}

func (r *SearchTree[T]) del(pre, n *TreeNode[T]) {
	// node have no child
	if xcontainer.IsNil[*TreeNode[T]](n.left) && xcontainer.IsNil[*TreeNode[T]](n.right) {
		if pre.left == n {
			pre.left = nil
		}
		if pre.right == n {
			pre.right = nil
		}
		return
	}
	// node have only one child
	var next *TreeNode[T]
	if xcontainer.IsNil[*TreeNode[T]](n.left) && !xcontainer.IsNil[*TreeNode[T]](n.right) {
		next = n.right
	}
	if xcontainer.IsNil[*TreeNode[T]](n.right) && !xcontainer.IsNil[*TreeNode[T]](n.left) {
		next = n.left
	}
	if !xcontainer.IsNil[*TreeNode[T]](next) {
		if pre.left == n {
			pre.left = next
		}
		if pre.right == n {
			pre.right = next
		}
		return
	}
	// have both child
	preMin, min := r.min(n, n.right)
	n.val = min.val
	r.del(preMin, min)
}

func (r *SearchTree[T]) min(pre, n *TreeNode[T]) (*TreeNode[T], *TreeNode[T]) {
	if !xcontainer.IsNil[*TreeNode[T]](n.left) {
		return r.min(n, n.left)
	}
	return pre, n
}

func (r *SearchTree[T]) sortPut(n *TreeNode[T], data T) {
	if n.val == data {
		return
	}
	if n.val > data {
		if xcontainer.IsNil[*TreeNode[T]](n.left) {
			n.left = &TreeNode[T]{val: data}
			return
		}
		r.sortPut(n.left, data)
	}
	if n.val < data {
		if xcontainer.IsNil[*TreeNode[T]](n.right) {
			n.right = &TreeNode[T]{val: data}
			return
		}
		r.sortPut(n.right, data)
	}
}
