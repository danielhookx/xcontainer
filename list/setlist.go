package list

import "github.com/danielhookx/xcontainer/set"

// SingleListSet 节点值不重复的单链表封装
type SingleListSet[T comparable] struct {
	len        int
	head, tail *SingleListNode[T]
	set        *set.Set[T]
}

func NewSingleListSet[T comparable]() *SingleListSet[T] {
	l := &SingleListSet[T]{
		len:  0,
		head: &SingleListNode[T]{Val: *new(T), Next: nil},
		tail: nil,
		set:  set.CreateSet[T](),
	}
	l.tail = l.head
	return l
}

func (l *SingleListSet[T]) Add(val T) {
	if l.set.IsElementOf(val) {
		return
	}
	l.tail.Next = &SingleListNode[T]{
		Val:  val,
		Next: nil,
	}
	l.tail = l.tail.Next
	l.set.Add(val)
	l.len++
}

func (l *SingleListSet[T]) Head() *SingleListNode[T] {
	return l.head.Next
}

func (l *SingleListSet[T]) Tail() *SingleListNode[T] {
	return l.tail
}
