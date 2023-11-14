package heap

import "github.com/danielhookx/xcontainer"

type MaxHeap[T xcontainer.Orderliness] struct {
	nodes    []T
	last     int
	maxCount int
}

func HeadBuildHeap[T xcontainer.Orderliness](src []T) *MaxHeap[T] {
	h := MaxHeap[T]{
		nodes:    src,
		last:     len(src) - 1,
		maxCount: len(src),
	}
	for i := len(src) / 2; i >= 0; i-- {
		h.down(i)
	}
	return &h
}

func TailBuildHeap[T xcontainer.Orderliness](src []T) *MaxHeap[T] {
	h := MaxHeap[T]{
		nodes:    make([]T, 0),
		last:     -1,
		maxCount: -1,
	}
	for _, i := range src {
		h.Add(i)
	}
	return &h
}

func (h *MaxHeap[T]) Add(item T) {
	if cap(h.nodes) <= h.last+1 {
		//grow
		h.grow(cap(h.nodes))
	}
	h.last++
	h.nodes[h.last] = item
	h.up(h.last)
}

func (h *MaxHeap[T]) Pop() T {
	if h.last <= 0 {
		return *new(T)
	}
	h.swap(0, h.last)
	h.last--
	h.down(0)
	return h.nodes[h.last+1]
}

func (h *MaxHeap[T]) up(j int) {
	for {
		i := (j - 1) / 2 //parent
		if i >= j || h.nodes[j] < h.nodes[i] {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *MaxHeap[T]) down(i int) {
	for {
		//left
		j1 := 2*i + 1
		if j1 > h.last || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 <= h.last && h.nodes[j1] <= h.nodes[j2] {
			j = j2
		}
		if h.nodes[j] <= h.nodes[i] {
			break
		}
		h.swap(i, j)
		i = j
	}
}

func (h *MaxHeap[T]) swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h *MaxHeap[T]) grow(size int) {
	if size == 0 {
		size = 100
	}
	new := make([]T, cap(h.nodes)+size)
	copy(new, h.nodes)
	h.nodes = new
}
