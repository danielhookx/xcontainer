package set

import (
	"errors"
	"iter"
)

// https://en.wikipedia.org/wiki/Set_(abstract_data_type)#Operations
type Set[T comparable] struct {
	cap int
	m   map[T]struct{}
}

// IsElementOf checks whether the value x is in the set S.
func (s *Set[T]) IsElementOf(x T) bool {
	_, ok := s.m[x]
	return ok
}

// Size returns the number of elements in S.
func (s *Set[T]) Size() int {
	return len(s.m)
}

// Iter returns an iterator over the elements of the set.
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := range s.m {
			if !yield(elem) {
				return
			}
		}
	}
}

// Enumerate returns a list containing the elements of S in some arbitrary order.
func (s *Set[T]) Enumerate() []T {
	list := make([]T, 0, len(s.m))
	for elem := range s.m {
		list = append(list, elem)
	}
	return list
}

// Add adds the element x to S, if it is not present already.
func (s *Set[T]) Add(x T) error {
	if s.cap > 0 && s.Size()+1 > s.cap {
		return errors.New("out of cap")
	}
	s.m[x] = struct{}{}
	return nil
}

// Remove removes the element x from S, if it is present.
func (s *Set[T]) Remove(x T) {
	delete(s.m, x)
}

// Capacity returns the maximum number of values that S can hold.
func (s *Set[T]) Capacity() int {
	return s.cap
}

// Pop returns an arbitrary element of S, deleting it from S
func (s *Set[T]) Pop() (v T, ok bool) {
	for item := range s.m {
		delete(s.m, item)
		return item, true
	}
	return v, false
}

// Pick returns an arbitrary element of S.
// Functionally, the mutator pop can be interpreted as the pair of selectors (pick, rest), where rest returns the set consisting of all elements except for the arbitrary element. Can be interpreted in terms of iterate.
func (s *Set[T]) Pick() {
	//TODO implement me
	panic("implement me")
}

// Map returns the set of distinct values resulting from applying function F to each element of S.
func (s *Set[T]) Map(f func(T) T) *Set[T] {
	newS := CreateSet[T]()
	for item := range s.m {
		newS.Add(f(item))
	}
	return newS
}

// Filter returns the subset containing all elements of S that satisfy a given predicate P.
func (s *Set[T]) Filter(p func(T) bool) *Set[T] {
	newS := CreateSet[T]()
	for item := range s.m {
		if p(item) {
			newS.Add(item)
		}
	}
	return newS
}

// fold returns the value A|S| after applying Ai+1 := F(Ai, e) for each element e of S, for some binary operation F. F must be associative and commutative for this to be well-defined.
// https://en.wikipedia.org/wiki/Fold_(higher-order_function)
func (s *Set[T]) fold(a0 T, f func(T) T) {
	//TODO implement me
	panic("implement me")
}

// Clear deletes all elements of S.
func (s *Set[T]) Clear() {
	// Constructions like this are optimised by compiler, and replaced by
	// mapclear() function, defined in
	// https://github.com/golang/go/blob/29bbca5c2c1ad41b2a9747890d183b6dd3a4ace4/src/runtime/map.go#L993)
	for key := range s.m {
		delete(s.m, key)
	}
}

// Equal checks whether the two given sets are equal (i.e. contain all and only the same elements).
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

// Hash returns a hash value for the static set S such that if equal(S1, S2) then hash(S1) = hash(S2)
func (s *Set[T]) Hash() string {
	//TODO implement me
	panic("implement me")
}

// BuildSet creates a set structure with values x1,x2,...,xn.
func BuildSet[T comparable](xn ...T) *Set[T] {
	s := CreateSet[T]()
	for _, v := range xn {
		s.Add(v)
	}
	return s
}

// Collect creates a new set structure containing all the elements of the given collection or all the elements returned by the given iterator.
func Collect[T comparable](collection iter.Seq[T]) *Set[T] {
	newS := CreateSet[T]()
	for item := range collection {
		newS.Add(item)
	}
	return newS
}

// CreateSet creates a new, initially empty set structure.
func CreateSet[T comparable]() *Set[T] {
	return &Set[T]{
		cap: 0,
		m:   map[T]struct{}{},
	}
}

// CreateSetWithCapacity creates a new set structure, initially empty but capable of holding up to n elements.
func CreateSetWithCapacity[T comparable](n int) *Set[T] {
	return &Set[T]{
		cap: n,
		m:   map[T]struct{}{},
	}
}

// IsEmptySet checks whether the set S is empty.
func IsEmptySet[T comparable](s *Set[T]) bool {
	return s == nil || len(s.m) == 0
}
