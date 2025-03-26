package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBasicOperations(t *testing.T) {
	t.Run("Create and check empty set", func(t *testing.T) {
		s := CreateSet[int]()
		assert.True(t, IsEmptySet(s))
		assert.Equal(t, 0, s.Size())
	})

	t.Run("Add and remove elements", func(t *testing.T) {
		s := CreateSet[string]()
		err := s.Add("test")
		assert.NoError(t, err)
		assert.True(t, s.IsElementOf("test"))
		s.Remove("test")
		assert.False(t, s.IsElementOf("test"))
	})

	t.Run("Set capacity", func(t *testing.T) {
		s := CreateSetWithCapacity[int](2)
		assert.Equal(t, 2, s.Capacity())
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		err = s.Add(3)
		assert.Error(t, err)
	})
}

func TestSetOperations(t *testing.T) {
	t.Run("Pop operation", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		val, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 1, val)
		assert.Equal(t, 0, s.Size())

		// Test case: pop from empty set
		val, ok = s.Pop()
		assert.False(t, ok)
		assert.Equal(t, 0, val) // for int type, zero value is 0
		assert.Equal(t, 0, s.Size())
	})

	t.Run("Clear operation", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		s.Clear()
		assert.Equal(t, 0, s.Size())
	})

	t.Run("Equal operation", func(t *testing.T) {
		s1 := CreateSet[int]()
		s2 := CreateSet[int]()
		err := s1.Add(1)
		assert.NoError(t, err)
		err = s2.Add(1)
		assert.NoError(t, err)
		assert.True(t, s1.Equal(s2))
		err = s2.Add(2)
		assert.NoError(t, err)
		assert.False(t, s1.Equal(s2))

		// Test case: sets with same size but different elements
		s3 := CreateSet[int]()
		s4 := CreateSet[int]()
		err = s3.Add(1)
		assert.NoError(t, err)
		err = s3.Add(2)
		assert.NoError(t, err)
		err = s4.Add(2)
		assert.NoError(t, err)
		err = s4.Add(3)
		assert.NoError(t, err)
		assert.Equal(t, s3.Size(), s4.Size())
		assert.False(t, s3.Equal(s4))

		// Test case: comparing with nil set
		s5 := CreateSet[int]()
		assert.True(t, s5.Equal(nil)) // empty set equals nil
		err = s5.Add(1)
		assert.NoError(t, err)
		assert.False(t, s5.Equal(nil)) // non-empty set not equals nil
	})
}

func TestSetCollectionOperations(t *testing.T) {
	t.Run("BuildSet", func(t *testing.T) {
		s := BuildSet(1, 2, 3)
		assert.Equal(t, 3, s.Size())
		assert.True(t, s.IsElementOf(1))
		assert.True(t, s.IsElementOf(2))
		assert.True(t, s.IsElementOf(3))
	})

	t.Run("Collect operation", func(t *testing.T) {
		// 创建一个迭代器函数
		seq := func(yield func(int) bool) {
			values := []int{1, 2, 3, 4, 5}
			for _, v := range values {
				if !yield(v) {
					return
				}
			}
		}

		s := Collect(seq)
		assert.Equal(t, 5, s.Size())
		assert.True(t, s.IsElementOf(1))
		assert.True(t, s.IsElementOf(2))
		assert.True(t, s.IsElementOf(3))
		assert.True(t, s.IsElementOf(4))
		assert.True(t, s.IsElementOf(5))
	})

	t.Run("Map operation", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		mapped := s.Map(func(x int) int { return x * 2 })
		assert.Equal(t, 2, mapped.Size())
		assert.True(t, mapped.IsElementOf(2))
		assert.True(t, mapped.IsElementOf(4))
	})

	t.Run("Filter operation", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		err = s.Add(3)
		assert.NoError(t, err)
		filtered := s.Filter(func(x int) bool { return x%2 == 0 })
		assert.Equal(t, 1, filtered.Size())
		assert.True(t, filtered.IsElementOf(2))
	})
}

func TestSetIteration(t *testing.T) {
	t.Run("Iter operation", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		count := 0
		for x := range s.Iter() {
			assert.True(t, s.IsElementOf(x))
			count++
		}
		assert.Equal(t, 2, count)
	})

	t.Run("Iter operation with early termination", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		err = s.Add(3)
		assert.NoError(t, err)
		err = s.Add(4)
		assert.NoError(t, err)
		err = s.Add(5)
		assert.NoError(t, err)

		count := 0
		for x := range s.Iter() {
			assert.True(t, s.IsElementOf(x))
			count++
			if count == 2 { // terminate early after the second element
				break
			}
		}
		assert.Equal(t, 2, count) // verify only the first two elements are iterated
	})

	t.Run("Enumerate operation", func(t *testing.T) {
		s := CreateSet[int]()
		err := s.Add(1)
		assert.NoError(t, err)
		err = s.Add(2)
		assert.NoError(t, err)
		list := s.Enumerate()
		assert.Equal(t, 2, len(list))
		for _, x := range list {
			assert.True(t, s.IsElementOf(x))
		}
	})
}
