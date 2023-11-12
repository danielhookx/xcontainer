package stack

type Stack[T any] struct {
	stack []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		stack: make([]T, 0),
	}
}

func (s *Stack[T]) Push(item T) {
	s.stack = append(s.stack, item)
}

func (s *Stack[T]) Pop() T {
	if len(s.stack) == 0 {
		return *new(T)
	}
	v := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return v
}

func (s *Stack[T]) Len() int {
	return len(s.stack)
}
