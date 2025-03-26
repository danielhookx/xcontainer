package list

import (
	"iter"

	"github.com/danielhookx/xcontainer/set"
)

// SetList 节点值不重复的单链表封装
type SetList[T comparable] struct {
	l *List[T]
	s *set.Set[T]
}

func NewSetList[T comparable]() *SetList[T] {
	sl := &SetList[T]{
		l: New[T](),
		s: set.CreateSet[T](),
	}
	return sl
}

func (sl *SetList[T]) Add(val T) {
	if sl.s.IsElementOf(val) {
		return
	}
	sl.l.PushFront(val)
	sl.s.Add(val)
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (sl *SetList[T]) Len() int { return sl.l.len }

func (sl *SetList[T]) Iter() iter.Seq[T] {
	return sl.l.Iter()
}
