package xcontainer

type CancelHandler func()

type IterateHandler[T comparable] func() (T, bool)

type Iterator[T comparable] interface {
	Iterate() (IterateHandler[T], CancelHandler)
}
