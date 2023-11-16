package xcontainer

type CancelHandler func()

type IterateHandler[T any] func() (T, bool)

type Iterator[T any] interface {
	Iterate() (IterateHandler[T], CancelHandler)
}
