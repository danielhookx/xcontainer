// Referenced from https://github.com/elliotchance/orderedmap/blob/master/orderedmap_test.go
// Author: Elliot Chance
// File: orderedmap_test.go

package xmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("ReturnsNotOKIfStringKeyDoesntExist", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		_, ok := m.Get("foo")
		assert.False(t, ok)
	})

	t.Run("ReturnsNotOKIfNonIntKeyDoesntExist", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		_, ok := m.Get(123)
		assert.False(t, ok)
	})

	t.Run("ReturnsOKIfKeyExists", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		_, ok := m.Get("foo")
		assert.True(t, ok)
	})

	t.Run("ReturnsValueForKey", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		value, _ := m.Get("foo")
		assert.Equal(t, "bar", value)
	})

	t.Run("KeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		m.Set("foo", "baz")
		_, ok := m.Get("bar")
		assert.False(t, ok)
	})

	t.Run("ValueForKeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		m.Set("foo", "baz")
		value, _ := m.Get("bar")
		assert.Empty(t, value)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Get(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_Get(400))

		// O(1) would mean that res4 should take about the same time as res1,
		// because we are accessing the same amount of elements, just on
		// different sized maps.

		assert.InDelta(t,
			res1.NsPerOp(), res4.NsPerOp(),
			0.5*float64(res1.NsPerOp()))
	})
}

func benchmarkOrderedMap_Get(multiplier int) func(b *testing.B) {
	m := NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Get(i % 1000 * multiplier)
		}
	}
}

func TestSet(t *testing.T) {
	t.Run("ReturnsTrueIfStringKeyIsNew", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		ok := m.Set("foo", "bar")
		assert.True(t, ok)
	})

	t.Run("ReturnsTrueIfNonStringKeyIsNew", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		ok := m.Set(123, "bar")
		assert.True(t, ok)
	})

	t.Run("ValueCanBeNonString", func(t *testing.T) {
		m := NewOrderedMap[int, bool]()
		ok := m.Set(123, true)
		assert.True(t, ok)
	})

	t.Run("ReturnsFalseIfKeyIsNotNew", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		ok := m.Set("foo", "bar")
		assert.False(t, ok)
	})

	t.Run("SetThreeDifferentKeys", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		m.Set("baz", "qux")
		ok := m.Set("quux", "corge")
		assert.True(t, ok)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Set(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_Set(400))

		// O(1) would mean that res4 should take about 4 times longer than res1
		// because we are doing 4 times the amount of Set operations. Allow for
		// a wide margin, but not too wide that it would permit the inflection
		// to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			2*float64(res1.NsPerOp()))
	})
}

func benchmarkOrderedMap_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := NewOrderedMap[int, bool]()
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}
	}
}

func TestLen(t *testing.T) {
	t.Run("EmptyMapIsZeroLen", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		assert.Equal(t, 0, m.Len())
	})

	t.Run("SingleElementIsLenOne", func(t *testing.T) {
		m := NewOrderedMap[int, bool]()
		m.Set(123, true)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("ThreeElements", func(t *testing.T) {
		m := NewOrderedMap[int, bool]()
		m.Set(1, true)
		m.Set(2, true)
		m.Set(3, true)
		assert.Equal(t, 3, m.Len())
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Len(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_Len(400))

		// O(1) would mean that res4 should take about the same time as res1,
		// because we are accessing the same amount of elements, just on
		// different sized maps.

		assert.InDelta(t,
			res1.NsPerOp(), res4.NsPerOp(),
			0.5*float64(res1.NsPerOp()))
	})
}

var tempInt int

func benchmarkOrderedMap_Len(multiplier int) func(b *testing.B) {
	m := NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		var temp int
		for i := 0; i < b.N; i++ {
			temp = m.Len()
		}

		// prevent compiler from optimising Len away.
		tempInt = temp
	}
}

func TestDelete(t *testing.T) {
	t.Run("KeyDoesntExistReturnsFalse", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		assert.False(t, m.Delete("foo"))
	})

	t.Run("KeyDoesExist", func(t *testing.T) {
		m := NewOrderedMap[string, interface{}]()
		m.Set("foo", nil)
		assert.True(t, m.Delete("foo"))
	})

	t.Run("KeyNoLongerExists", func(t *testing.T) {
		m := NewOrderedMap[string, interface{}]()
		m.Set("foo", nil)
		m.Delete("foo")
		_, exists := m.Get("foo")
		assert.False(t, exists)
	})

	t.Run("KeyDeleteIsIsolated", func(t *testing.T) {
		m := NewOrderedMap[string, interface{}]()
		m.Set("foo", nil)
		m.Set("bar", nil)
		m.Delete("foo")
		_, exists := m.Get("bar")
		assert.True(t, exists)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Delete(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_Delete(400))

		// O(1) would mean that res4 should take about 4 times longer than res1
		// because we are doing 4 times the amount of Set/Delete operations.
		// Allow for a wide margin, but not too wide that it would permit the
		// inflection to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			float64(res4.NsPerOp()))
	})
}

func benchmarkOrderedMap_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := NewOrderedMap[int, bool]()
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}

		for i := 0; i < b.N; i++ {
			m.Delete(i)
		}
	}
}

func TestOrderedMap_Copy(t *testing.T) {
	t.Run("ReturnsEqualButNotSame", func(t *testing.T) {
		key, value := 1, "a value"
		m := NewOrderedMap[int, string]()
		m.Set(key, value)

		m2 := m.Copy()
		m2.Set(key, "a different value")

		assert.Equal(t, m.Len(), m2.Len(), "not all elements are copied")
		newV, _ := m.Get(key)
		assert.Equal(t, value, newV)
	})
}

func TestOrderedMap_ToArray(t *testing.T) {
	t.Run("ReturnsEqualButNotSame", func(t *testing.T) {
		key1, value1 := 1, "a value 1"
		key2, value2 := 2, "a value 2"
		key3, value3 := 3, "a value 3"
		m := NewOrderedMap[int, string]()
		m.Set(key1, value1)
		m.Set(key2, value2)
		m.Set(key3, value3)

		array := m.ToArray()
		assert.Equal(t, m.Len(), len(array))
		assert.Equal(t, value1, array[0])
		assert.Equal(t, value2, array[1])
		assert.Equal(t, value3, array[2])
	})
}
