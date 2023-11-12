package heap

import (
	"math/rand"
	"testing"
)

func TestHeapSort(t *testing.T) {
	const count = 20
	var src = make([]int, count)
	for i := 0; i < count; i++ {
		src[i] = rand.Intn(count)
	}
	t.Log(src)
	h := HeadBuildHeap[int](src)
	t.Log(h.nodes)
	preV := h.Pop()
	for h.last > 0 {
		v := h.Pop()
		if preV < v {
			t.Error("not maxHeap")
		}
		preV = v
	}
}

func TestHeadBuildHeap(t *testing.T) {
	h := HeadBuildHeap[int]([]int{7, 2, 1, 4, 5, 6, 3})
	t.Log(h.nodes)
}

func TestTailBuildHeap(t *testing.T) {
	h := TailBuildHeap[int]([]int{7, 2, 1, 4, 5, 6, 3})
	t.Log(h.nodes)
	preV := h.Pop()
	for i := 0; i < len(h.nodes); i++ {
		v := h.Pop()
		if preV < v {
			t.Error("not maxHeap")
		}
		preV = v
	}
	t.Log(h.nodes)
}
