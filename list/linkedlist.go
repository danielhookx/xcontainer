package list

type SingleListNode[T any] struct {
	Val  T
	Next *SingleListNode[T]
}

func NewSingleListFromArray[T any](src []T) *SingleListNode[T] {
	dummy := &SingleListNode[T]{}
	p := dummy
	for _, v := range src {
		p.Next = &SingleListNode[T]{Val: v}
		p = p.Next
	}
	return dummy.Next
}
