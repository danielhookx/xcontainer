package stack

import (
	"testing"
)

func TestNewStackInt(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Push(4)

	for stack.Len() > 0 {
		t.Log(stack.Pop())
	}
	t.Log(stack.Pop())
	t.Log(stack.Pop())
}
