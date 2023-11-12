package queue

import (
	"testing"
)

func TestNewQueueInt(t *testing.T) {
	queue := NewQueue[int]()
	queue.EnQueue(1)
	queue.EnQueue(2)
	queue.EnQueue(3)
	queue.EnQueue(4)

	for queue.Len() > 0 {
		t.Log(queue.DeQueue())
	}
	t.Log(queue.DeQueue())
	t.Log(queue.DeQueue())
}
