package queue

type Queue[T any] struct {
	queue []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		queue: make([]T, 0),
	}
}

func (q *Queue[T]) EnQueue(v T) {
	q.queue = append(q.queue, v)
}

func (q *Queue[T]) DeQueue() T {
	if len(q.queue) == 0 {
		return *new(T)
	}
	v := q.queue[0]
	q.queue = q.queue[1:]
	return v
}

func (q *Queue[T]) Len() int {
	return len(q.queue)
}
