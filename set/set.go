package set

import (
	"errors"
	"sync/atomic"

	"github.com/danielhookx/xcontainer"
)

// https://en.wikipedia.org/wiki/Set_(abstract_data_type)#Operations
type Set[T comparable] struct {
	cap int
	m   map[T]struct{}
}

// checks whether the value x is in the set S.
func (s *Set[T]) IsElementOf(x T) bool {
	_, ok := s.m[x]
	return ok
}

// returns the number of elements in S.
func (s *Set[T]) Size() int {
	return len(s.m)
}

// returns a function that returns one more value of S at each call, in some arbitrary order.
func (s *Set[T]) Iterate() (xcontainer.IterateHandler[T], xcontainer.CancelHandler) {
	ch := make(chan T)
	var isStoped atomic.Bool
	stopCh := make(chan struct{})

	go func() {
		defer close(ch)
		for elem := range s.m {
			select {
			case <-stopCh:
				return
			case ch <- elem:
			}
		}
	}()
	return func() (T, bool) {
			v, ok := <-ch
			return v, ok
		}, func() {
			if isStoped.CompareAndSwap(false, true) {
				close(stopCh)
			}
		}
}

// returns a list containing the elements of S in some arbitrary order.
func (s *Set[T]) Enumerate() []T {
	list := make([]T, 0, len(s.m))
	for elem := range s.m {
		list = append(list, elem)
	}
	return list
}

// adds the element x to S, if it is not present already.
func (s *Set[T]) Add(x T) error {
	if s.cap > 0 && s.Size()+1 > s.cap {
		return errors.New("out of cap")
	}
	s.m[x] = struct{}{}
	return nil
}

// removes the element x from S, if it is present.
func (s *Set[T]) Remove(x T) {
	delete(s.m, x)
}

// returns the maximum number of values that S can hold.
func (s *Set[T]) Capacity() int {
	return s.cap
}

// returns an arbitrary element of S, deleting it from S
func (s *Set[T]) Pop() (v T, ok bool) {
	for item := range s.m {
		delete(s.m, item)
		return item, true
	}
	return v, false
}

// returns an arbitrary element of S.
// Functionally, the mutator pop can be interpreted as the pair of selectors (pick, rest), where rest returns the set consisting of all elements except for the arbitrary element. Can be interpreted in terms of iterate.
func (s *Set[T]) Pick() {
	//TODO implement me
	panic("implement me")
}

// returns the set of distinct values resulting from applying function F to each element of S.
func (s *Set[T]) Map(f func(T) T) *Set[T] {
	iterate, cancel := s.Iterate()
	defer cancel()
	newS := CreateSet[T]()
	for {
		v, ok := iterate()
		if !ok {
			break
		}
		newS.Add(f(v))
	}
	return newS
}

// returns the subset containing all elements of S that satisfy a given predicate P.
func (s *Set[T]) Filter(p func(T) bool) *Set[T] {
	iterate, cancel := s.Iterate()
	defer cancel()
	newS := CreateSet[T]()
	for {
		v, ok := iterate()
		if !ok {
			break
		}
		if p(v) {
			newS.Add(v)
		}
	}
	return newS
}

// returns the value A|S| after applying Ai+1 := F(Ai, e) for each element e of S, for some binary operation F. F must be associative and commutative for this to be well-defined.
// https://en.wikipedia.org/wiki/Fold_(higher-order_function)
func (s *Set[T]) fold(a0 T, f func(T) T) {
	//TODO implement me
	panic("implement me")
}

// delete all elements of S.
func (s *Set[T]) Clear() {
	// Constructions like this are optimised by compiler, and replaced by
	// mapclear() function, defined in
	// https://github.com/golang/go/blob/29bbca5c2c1ad41b2a9747890d183b6dd3a4ace4/src/runtime/map.go#L993)
	for key := range s.m {
		delete(s.m, key)
	}
}

// checks whether the two given sets are equal (i.e. contain all and only the same elements).
func (s *Set[T]) Equal(s2 *Set[T]) bool {
	if s2 == nil {
		return s.Size() == 0
	}
	if s.Size() != s2.Size() {
		return false
	}
	for item := range s.m {
		if !s2.IsElementOf(item) {
			return false
		}
	}
	return true
}

// returns a hash value for the static set S such that if equal(S1, S2) then hash(S1) = hash(S2)
func (s *Set[T]) Hash() string {
	//TODO implement me
	panic("implement me")
}

// creates a set structure with values x1,x2,...,xn.
func BuildSet[T comparable](xn ...T) *Set[T] {
	s := CreateSet[T]()
	for _, v := range xn {
		s.Add(v)
	}
	return s
}

// creates a new set structure containing all the elements of the given collection or all the elements returned by the given iterator.
func CreateFrom[T comparable](collection xcontainer.Iterator[T]) *Set[T] {
	iterate, cancel := collection.Iterate()
	defer cancel()
	newS := CreateSet[T]()
	for {
		v, ok := iterate()
		if !ok {
			break
		}
		newS.Add(v)
	}
	return newS
}

// creates a new, initially empty set structure.
func CreateSet[T comparable]() *Set[T] {
	return &Set[T]{
		cap: 0,
		m:   map[T]struct{}{},
	}
}

// creates a new set structure, initially empty but capable of holding up to n elements.
func CreateSetWithCapacity[T comparable](n int) *Set[T] {
	return &Set[T]{
		cap: n,
		m:   map[T]struct{}{},
	}
}

// checks whether the set S is empty.
func IsEmptySet[T comparable](s *Set[T]) bool {
	return s == nil || len(s.m) == 0
}
