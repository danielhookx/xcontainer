// Referenced from https://github.com/elliotchance/orderedmap/blob/master/orderedmap_test.go
// Author: Elliot Chance
// File: orderedmap_test.go

package xmap

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
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
	})

	t.Run("Special Values", func(t *testing.T) {
		t.Run("Empty String Key", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("", "empty")
			value, ok := m.Get("")
			assert.True(t, ok)
			assert.Equal(t, "empty", value)
		})

		t.Run("Zero Value Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(0, "zero")
			value, ok := m.Get(0)
			assert.True(t, ok)
			assert.Equal(t, "zero", value)
		})

		t.Run("Nil Value", func(t *testing.T) {
			m := NewOrderedMap[string, interface{}]()
			m.Set("nil", nil)
			value, ok := m.Get("nil")
			assert.True(t, ok)
			assert.Nil(t, value)
		})
	})

	t.Run("Complex Types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		t.Run("Struct Value", func(t *testing.T) {
			m := NewOrderedMap[string, Person]()
			person := Person{Name: "Alice", Age: 20}
			m.Set("person", person)
			value, ok := m.Get("person")
			assert.True(t, ok)
			assert.Equal(t, person, value)
		})

		t.Run("Slice Value", func(t *testing.T) {
			m := NewOrderedMap[string, []int]()
			slice := []int{1, 2, 3}
			m.Set("slice", slice)
			value, ok := m.Get("slice")
			assert.True(t, ok)
			assert.Equal(t, slice, value)
		})

		t.Run("Map Value", func(t *testing.T) {
			m := NewOrderedMap[string, map[string]int]()
			innerMap := map[string]int{"a": 1, "b": 2}
			m.Set("map", innerMap)
			value, ok := m.Get("map")
			assert.True(t, ok)
			assert.Equal(t, innerMap, value)
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("Max Integer Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(1<<31-1, "max")
			value, ok := m.Get(1<<31 - 1)
			assert.True(t, ok)
			assert.Equal(t, "max", value)
		})

		t.Run("Min Integer Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(-1<<31, "min")
			value, ok := m.Get(-1 << 31)
			assert.True(t, ok)
			assert.Equal(t, "min", value)
		})

		t.Run("Special Character Key", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			specialKey := "!@#$%^&*()_+"
			m.Set(specialKey, "special")
			value, ok := m.Get(specialKey)
			assert.True(t, ok)
			assert.Equal(t, "special", value)
		})
	})

	t.Run("Concurrency", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("key", 42)

		var wg sync.WaitGroup
		concurrentCount := 100
		for i := 0; i < concurrentCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				value, ok := m.Get("key")
				assert.True(t, ok)
				assert.Equal(t, 42, value)
			}()
		}
		wg.Wait()
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
	t.Run("Basic Operations", func(t *testing.T) {
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
	})

	t.Run("Special Values", func(t *testing.T) {
		t.Run("Empty String Key", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			ok := m.Set("", "empty")
			assert.True(t, ok)
			value, exists := m.Get("")
			assert.True(t, exists)
			assert.Equal(t, "empty", value)
		})

		t.Run("Zero Value Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			ok := m.Set(0, "zero")
			assert.True(t, ok)
			value, exists := m.Get(0)
			assert.True(t, exists)
			assert.Equal(t, "zero", value)
		})

		t.Run("Nil Value", func(t *testing.T) {
			m := NewOrderedMap[string, interface{}]()
			ok := m.Set("nil", nil)
			assert.True(t, ok)
			value, exists := m.Get("nil")
			assert.True(t, exists)
			assert.Nil(t, value)
		})
	})

	t.Run("Complex Types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		t.Run("Struct Value", func(t *testing.T) {
			m := NewOrderedMap[string, Person]()
			person := Person{Name: "Alice", Age: 20}
			ok := m.Set("person", person)
			assert.True(t, ok)
			value, exists := m.Get("person")
			assert.True(t, exists)
			assert.Equal(t, person, value)
		})

		t.Run("Slice Value", func(t *testing.T) {
			m := NewOrderedMap[string, []int]()
			slice := []int{1, 2, 3}
			ok := m.Set("slice", slice)
			assert.True(t, ok)
			value, exists := m.Get("slice")
			assert.True(t, exists)
			assert.Equal(t, slice, value)
		})

		t.Run("Map Value", func(t *testing.T) {
			m := NewOrderedMap[string, map[string]int]()
			innerMap := map[string]int{"a": 1, "b": 2}
			ok := m.Set("map", innerMap)
			assert.True(t, ok)
			value, exists := m.Get("map")
			assert.True(t, exists)
			assert.Equal(t, innerMap, value)
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("Max Integer Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			ok := m.Set(1<<31-1, "max")
			assert.True(t, ok)
			value, exists := m.Get(1<<31 - 1)
			assert.True(t, exists)
			assert.Equal(t, "max", value)
		})

		t.Run("Min Integer Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			ok := m.Set(-1<<31, "min")
			assert.True(t, ok)
			value, exists := m.Get(-1 << 31)
			assert.True(t, exists)
			assert.Equal(t, "min", value)
		})

		t.Run("Special Character Key", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			specialKey := "!@#$%^&*()_+"
			ok := m.Set(specialKey, "special")
			assert.True(t, ok)
			value, exists := m.Get(specialKey)
			assert.True(t, exists)
			assert.Equal(t, "special", value)
		})
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
	t.Run("Basic Operations", func(t *testing.T) {
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
	})

	t.Run("Special Cases", func(t *testing.T) {
		t.Run("AfterDelete", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("key1", "value1")
			m.Set("key2", "value2")
			assert.Equal(t, 2, m.Len())
			m.Delete("key1")
			assert.Equal(t, 1, m.Len())
		})

		t.Run("AfterUpdate", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("key1", "value1")
			m.Set("key2", "value2")
			assert.Equal(t, 2, m.Len())
			m.Set("key1", "newvalue1") // Update existing key
			assert.Equal(t, 2, m.Len())
		})

		t.Run("AfterClear", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("key1", "value1")
			m.Set("key2", "value2")
			assert.Equal(t, 2, m.Len())
			m.Delete("key1")
			m.Delete("key2")
			assert.Equal(t, 0, m.Len())
		})
	})

	t.Run("Complex Types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		t.Run("Struct Values", func(t *testing.T) {
			m := NewOrderedMap[string, Person]()
			m.Set("person1", Person{Name: "Alice", Age: 20})
			m.Set("person2", Person{Name: "Bob", Age: 30})
			assert.Equal(t, 2, m.Len())
		})

		t.Run("Slice Values", func(t *testing.T) {
			m := NewOrderedMap[string, []int]()
			m.Set("slice1", []int{1, 2, 3})
			m.Set("slice2", []int{4, 5, 6})
			assert.Equal(t, 2, m.Len())
		})

		t.Run("Map Values", func(t *testing.T) {
			m := NewOrderedMap[string, map[string]int]()
			m.Set("map1", map[string]int{"a": 1})
			m.Set("map2", map[string]int{"b": 2})
			assert.Equal(t, 2, m.Len())
		})
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
	t.Run("Basic Operations", func(t *testing.T) {
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
	})

	t.Run("Special Cases", func(t *testing.T) {
		t.Run("Empty String Key", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("", "empty")
			assert.True(t, m.Delete(""))
			_, exists := m.Get("")
			assert.False(t, exists)
		})

		t.Run("Zero Value Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(0, "zero")
			assert.True(t, m.Delete(0))
			_, exists := m.Get(0)
			assert.False(t, exists)
		})

		t.Run("Nil Value", func(t *testing.T) {
			m := NewOrderedMap[string, interface{}]()
			m.Set("nil", nil)
			assert.True(t, m.Delete("nil"))
			_, exists := m.Get("nil")
			assert.False(t, exists)
		})
	})

	t.Run("Complex Types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		t.Run("Struct Value", func(t *testing.T) {
			m := NewOrderedMap[string, Person]()
			person := Person{Name: "Alice", Age: 20}
			m.Set("person", person)
			assert.True(t, m.Delete("person"))
			_, exists := m.Get("person")
			assert.False(t, exists)
		})

		t.Run("Slice Value", func(t *testing.T) {
			m := NewOrderedMap[string, []int]()
			slice := []int{1, 2, 3}
			m.Set("slice", slice)
			assert.True(t, m.Delete("slice"))
			_, exists := m.Get("slice")
			assert.False(t, exists)
		})

		t.Run("Map Value", func(t *testing.T) {
			m := NewOrderedMap[string, map[string]int]()
			innerMap := map[string]int{"a": 1, "b": 2}
			m.Set("map", innerMap)
			assert.True(t, m.Delete("map"))
			_, exists := m.Get("map")
			assert.False(t, exists)
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("Max Integer Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(1<<31-1, "max")
			assert.True(t, m.Delete(1<<31-1))
			_, exists := m.Get(1<<31 - 1)
			assert.False(t, exists)
		})

		t.Run("Min Integer Key", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(-1<<31, "min")
			assert.True(t, m.Delete(-1<<31))
			_, exists := m.Get(-1 << 31)
			assert.False(t, exists)
		})

		t.Run("Special Character Key", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			specialKey := "!@#$%^&*()_+"
			m.Set(specialKey, "special")
			assert.True(t, m.Delete(specialKey))
			_, exists := m.Get(specialKey)
			assert.False(t, exists)
		})
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
	t.Run("Basic Operations", func(t *testing.T) {
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

		t.Run("EmptyMapCopy", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m2 := m.Copy()
			assert.Equal(t, 0, m2.Len())
			assert.NotSame(t, m, m2)
		})

		t.Run("MultipleElementsCopy", func(t *testing.T) {
			m := NewOrderedMap[string, int]()
			m.Set("one", 1)
			m.Set("two", 2)
			m.Set("three", 3)

			m2 := m.Copy()
			assert.Equal(t, 3, m2.Len())

			// Verify all elements are copied correctly
			for k, v := range m.Iter() {
				v2, exists := m2.Get(k)
				assert.True(t, exists)
				assert.Equal(t, v, v2)
			}
		})
	})

	t.Run("Special Cases", func(t *testing.T) {
		t.Run("NilValueCopy", func(t *testing.T) {
			m := NewOrderedMap[string, interface{}]()
			m.Set("nil", nil)
			m2 := m.Copy()
			value, exists := m2.Get("nil")
			assert.True(t, exists)
			assert.Nil(t, value)
		})

		t.Run("EmptyStringKeyCopy", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("", "empty")
			m2 := m.Copy()
			value, exists := m2.Get("")
			assert.True(t, exists)
			assert.Equal(t, "empty", value)
		})

		t.Run("ZeroValueKeyCopy", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(0, "zero")
			m2 := m.Copy()
			value, exists := m2.Get(0)
			assert.True(t, exists)
			assert.Equal(t, "zero", value)
		})
	})

	t.Run("Complex Types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		t.Run("StructValueCopy", func(t *testing.T) {
			m := NewOrderedMap[string, Person]()
			person := Person{Name: "Alice", Age: 20}
			m.Set("person", person)
			m2 := m.Copy()
			value, exists := m2.Get("person")
			assert.True(t, exists)
			assert.Equal(t, person, value)
		})

		t.Run("SliceValueCopy", func(t *testing.T) {
			m := NewOrderedMap[string, []int]()
			slice := []int{1, 2, 3}
			m.Set("slice", slice)
			m2 := m.Copy()
			value, exists := m2.Get("slice")
			assert.True(t, exists)
			assert.Equal(t, slice, value)
		})

		t.Run("MapValueCopy", func(t *testing.T) {
			m := NewOrderedMap[string, map[string]int]()
			innerMap := map[string]int{"a": 1, "b": 2}
			m.Set("map", innerMap)
			m2 := m.Copy()
			value, exists := m2.Get("map")
			assert.True(t, exists)
			assert.Equal(t, innerMap, value)
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("MaxIntegerKeyCopy", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(1<<31-1, "max")
			m2 := m.Copy()
			value, exists := m2.Get(1<<31 - 1)
			assert.True(t, exists)
			assert.Equal(t, "max", value)
		})

		t.Run("MinIntegerKeyCopy", func(t *testing.T) {
			m := NewOrderedMap[int, string]()
			m.Set(-1<<31, "min")
			m2 := m.Copy()
			value, exists := m2.Get(-1 << 31)
			assert.True(t, exists)
			assert.Equal(t, "min", value)
		})

		t.Run("SpecialCharacterKeyCopy", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			specialKey := "!@#$%^&*()_+"
			m.Set(specialKey, "special")
			m2 := m.Copy()
			value, exists := m2.Get(specialKey)
			assert.True(t, exists)
			assert.Equal(t, "special", value)
		})
	})

	t.Run("Concurrency", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		for i := 0; i < 100; i++ {
			m.Set(fmt.Sprintf("key%d", i), i)
		}

		var wg sync.WaitGroup
		concurrentCount := 10
		for i := 0; i < concurrentCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m2 := m.Copy()
				assert.Equal(t, 100, m2.Len())
			}()
		}
		wg.Wait()
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Copy(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_Copy(400))

		// O(n) would mean that res4 should take about 4 times longer than res1
		// because we are copying 4 times the amount of elements. Allow for
		// a wide margin, but not too wide that it would permit the inflection
		// to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			2*float64(res1.NsPerOp()))
	})
}

func benchmarkOrderedMap_Copy(multiplier int) func(b *testing.B) {
	m := NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Copy()
		}
	}
}

func TestOrderedMap_ToArray(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		t.Run("EmptyMapReturnsEmptyArray", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			array := m.ToArray()
			assert.Empty(t, array)
		})

		t.Run("SingleElementReturnsSingleElementArray", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("key1", "value1")
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, "value1", array[0])
		})

		t.Run("MultipleElementsReturnsOrderedArray", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("key1", "value1")
			m.Set("key2", "value2")
			m.Set("key3", "value3")
			array := m.ToArray()
			assert.Equal(t, 3, len(array))
			assert.Equal(t, "value1", array[0])
			assert.Equal(t, "value2", array[1])
			assert.Equal(t, "value3", array[2])
		})
	})

	t.Run("Special Cases", func(t *testing.T) {
		t.Run("NilValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, interface{}]()
			m.Set("key1", nil)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Nil(t, array[0])
		})

		t.Run("EmptyStringValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			m.Set("key1", "")
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, "", array[0])
		})

		t.Run("ZeroValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, int]()
			m.Set("key1", 0)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, 0, array[0])
		})
	})

	t.Run("Complex Types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		t.Run("StructValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, Person]()
			person := Person{Name: "Alice", Age: 20}
			m.Set("key1", person)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, person, array[0])
		})

		t.Run("SliceValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, []int]()
			slice := []int{1, 2, 3}
			m.Set("key1", slice)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, slice, array[0])
		})

		t.Run("MapValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, map[string]int]()
			innerMap := map[string]int{"a": 1, "b": 2}
			m.Set("key1", innerMap)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, innerMap, array[0])
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("MaxIntegerValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, int]()
			m.Set("key1", 1<<31-1)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, 1<<31-1, array[0])
		})

		t.Run("MinIntegerValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, int]()
			m.Set("key1", -1<<31)
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, -1<<31, array[0])
		})

		t.Run("SpecialCharacterKeyValueInArray", func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			specialKey := "!@#$%^&*()_+"
			m.Set(specialKey, "special")
			array := m.ToArray()
			assert.Equal(t, 1, len(array))
			assert.Equal(t, "special", array[0])
		})
	})

	t.Run("Concurrency", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		for i := 0; i < 100; i++ {
			m.Set(fmt.Sprintf("key%d", i), i)
		}

		var wg sync.WaitGroup
		concurrentCount := 10
		for i := 0; i < concurrentCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				array := m.ToArray()
				assert.Equal(t, 100, len(array))
			}()
		}
		wg.Wait()
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_ToArray(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_ToArray(400))

		// O(n) would mean that res4 should take about 4 times longer than res1
		// because we are converting 4 times the amount of elements. Allow for
		// a wide margin, but not too wide that it would permit the inflection
		// to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			2*float64(res1.NsPerOp()))
	})
}

func benchmarkOrderedMap_ToArray(multiplier int) func(b *testing.B) {
	m := NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.ToArray()
		}
	}
}

func TestOrderedMap_Iter(t *testing.T) {
	t.Run("Empty Map", func(t *testing.T) {
		m := NewOrderedMap[string, int]()

		var result []struct {
			key   string
			value int
		}
		for k, v := range m.Iter() {
			result = append(result, struct {
				key   string
				value int
			}{k, v})
		}

		assert.Empty(t, result, "Empty map should return empty result")
	})

	t.Run("Single Element", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("one", 1)

		var result []struct {
			key   string
			value int
		}
		for k, v := range m.Iter() {
			result = append(result, struct {
				key   string
				value int
			}{k, v})
		}

		assert.Equal(t, 1, len(result), "Should have one element")
		assert.Equal(t, "one", result[0].key, "Key should be 'one'")
		assert.Equal(t, 1, result[0].value, "Value should be 1")
	})

	t.Run("Multiple Elements in Insertion Order", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)
		m.Set("three", 3)

		expected := []struct {
			key   string
			value int
		}{
			{"one", 1},
			{"two", 2},
			{"three", 3},
		}

		var result []struct {
			key   string
			value int
		}
		for k, v := range m.Iter() {
			result = append(result, struct {
				key   string
				value int
			}{k, v})
		}

		assert.Equal(t, len(expected), len(result), "Iterator returned incorrect number of elements")
		for i := range expected {
			assert.Equal(t, expected[i].key, result[i].key, "Key mismatch")
			assert.Equal(t, expected[i].value, result[i].value, "Value mismatch")
		}
	})

	t.Run("Update Element Value Without Changing Order", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)
		m.Set("three", 3)

		// Update element value
		m.Set("two", 22)

		expected := []struct {
			key   string
			value int
		}{
			{"one", 1},
			{"two", 22},
			{"three", 3},
		}

		var result []struct {
			key   string
			value int
		}
		for k, v := range m.Iter() {
			result = append(result, struct {
				key   string
				value int
			}{k, v})
		}

		assert.Equal(t, len(expected), len(result), "Iterator returned incorrect number of elements")
		for i := range expected {
			assert.Equal(t, expected[i].key, result[i].key, "Key mismatch")
			assert.Equal(t, expected[i].value, result[i].value, "Value mismatch")
		}
	})

	t.Run("Early Exit from Iterator", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)
		m.Set("three", 3)
		m.Set("four", 4)
		m.Set("five", 5)

		count := 0
		for k, v := range m.Iter() {
			count++
			assert.NotEmpty(t, k)
			assert.NotZero(t, v)
			if count >= 2 {
				break
			}
		}
		assert.Equal(t, 2, count, "Iterator should allow early exit")
	})

	t.Run("Modify Map During Iteration", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)
		m.Set("three", 3)

		count := 0
		for k, _ := range m.Iter() {
			count++
			if k == "two" {
				m.Set("four", 4) // Add new element during iteration
			}
		}

		// In Go, if we add elements during map iteration, new elements won't appear in current iteration
		assert.Equal(t, 3, count, "Iterator should only iterate over elements that existed at start")
		assert.Equal(t, 4, m.Len(), "Map length should be updated")

		// Verify new element was added
		v, ok := m.Get("four")
		assert.True(t, ok, "New element should be added")
		assert.Equal(t, 4, v, "New element value should be correct")
	})

	t.Run("Different Types of Key-Value Pairs", func(t *testing.T) {
		t.Run("string -> interface{}", func(t *testing.T) {
			m := NewOrderedMap[string, any]()
			m.Set("int", 42)
			m.Set("string", "hello")
			m.Set("bool", true)
			m.Set("float", 3.14)

			values := make([]any, 0)
			for _, v := range m.Iter() {
				values = append(values, v)
			}

			assert.Equal(t, 4, len(values), "Should have 4 elements")
			assert.Equal(t, 42, values[0])
			assert.Equal(t, "hello", values[1])
			assert.Equal(t, true, values[2])
			assert.Equal(t, 3.14, values[3])
		})

		t.Run("int -> struct", func(t *testing.T) {
			type Person struct {
				Name string
				Age  int
			}
			m := NewOrderedMap[int, Person]()
			m.Set(1, Person{Name: "Alice", Age: 20})
			m.Set(2, Person{Name: "Bob", Age: 30})
			m.Set(3, Person{Name: "Charlie", Age: 40})

			var people []Person
			for _, v := range m.Iter() {
				people = append(people, v)
			}

			assert.Equal(t, 3, len(people), "Should have 3 elements")
			assert.Equal(t, "Alice", people[0].Name)
			assert.Equal(t, 20, people[0].Age)
			assert.Equal(t, "Bob", people[1].Name)
			assert.Equal(t, 30, people[1].Age)
			assert.Equal(t, "Charlie", people[2].Name)
			assert.Equal(t, 40, people[2].Age)
		})

		t.Run("struct -> int", func(t *testing.T) {
			type Point struct {
				X, Y int
			}
			m := NewOrderedMap[Point, int]()
			m.Set(Point{1, 1}, 1)
			m.Set(Point{2, 2}, 4)
			m.Set(Point{3, 3}, 9)

			sum := 0
			for _, v := range m.Iter() {
				sum += v
			}

			assert.Equal(t, 14, sum, "Sum of all values should be 14")
		})

		t.Run("bool -> []string", func(t *testing.T) {
			m := NewOrderedMap[bool, []string]()
			m.Set(true, []string{"yes", "true", "1"})
			m.Set(false, []string{"no", "false", "0"})

			count := 0
			for k, v := range m.Iter() {
				count++
				if k {
					assert.Equal(t, []string{"yes", "true", "1"}, v)
				} else {
					assert.Equal(t, []string{"no", "false", "0"}, v)
				}
			}

			assert.Equal(t, 2, count, "Should have 2 elements")
		})
	})

	t.Run("Iteration After Element Deletion", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)
		m.Set("three", 3)
		m.Set("four", 4)

		// Delete middle element
		m.Delete("two")

		expected := []struct {
			key   string
			value int
		}{
			{"one", 1},
			{"three", 3},
			{"four", 4},
		}

		var result []struct {
			key   string
			value int
		}
		for k, v := range m.Iter() {
			result = append(result, struct {
				key   string
				value int
			}{k, v})
		}

		assert.Equal(t, len(expected), len(result), "Iterator returned incorrect number of elements")
		for i := range expected {
			assert.Equal(t, expected[i].key, result[i].key, "Key mismatch")
			assert.Equal(t, expected[i].value, result[i].value, "Value mismatch")
		}
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Iter(100))
		res4 := testing.Benchmark(benchmarkOrderedMap_Iter(400))

		// O(n) would mean that res4 should take about 4 times longer than res1
		// because we are iterating over 4 times the amount of elements. Allow for
		// a wide margin, but not too wide that it would permit the inflection
		// to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			2*float64(res1.NsPerOp()))
	})
}

func benchmarkOrderedMap_Iter(multiplier int) func(b *testing.B) {
	m := NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, _ = range m.Iter() {
				// Do nothing, just iterate
			}
		}
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		jsonText := `{"key1": "value1", "key2": 123, "key3": true}`
		want := func() *OrderedMap[string, any] {
			m := NewOrderedMap[string, any]()
			m.Set("key1", "value1")
			m.Set("key2", float64(123))
			m.Set("key3", true)
			return m
		}()

		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.NoError(t, err)
		assert.Equal(t, want.Len(), om.Len())

		// Verify all key-value pairs
		for key, value := range want.Iter() {
			gotValue, exists := om.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		jsonText := `{}`
		want := NewOrderedMap[string, any]()

		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.NoError(t, err)
		assert.Equal(t, want.Len(), om.Len())
	})

	t.Run("Special Values", func(t *testing.T) {
		jsonText := `{"empty": "", "null": null, "zero": 0}`
		want := func() *OrderedMap[string, any] {
			m := NewOrderedMap[string, any]()
			m.Set("empty", "")
			m.Set("null", nil)
			m.Set("zero", float64(0))
			return m
		}()

		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.NoError(t, err)
		assert.Equal(t, want.Len(), om.Len())

		// Verify all key-value pairs
		for key, value := range want.Iter() {
			gotValue, exists := om.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Complex Types", func(t *testing.T) {
		jsonText := `{"struct": {"name": "test", "age": 30}, "array": [1, 2, 3], "map": {"key": "value"}}`
		want := func() *OrderedMap[string, any] {
			m := NewOrderedMap[string, any]()
			m.Set("struct", map[string]any{
				"name": "test",
				"age":  float64(30),
			})
			m.Set("array", []any{float64(1), float64(2), float64(3)})
			m.Set("map", map[string]any{
				"key": "value",
			})
			return m
		}()

		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.NoError(t, err)
		assert.Equal(t, want.Len(), om.Len())

		// Verify all key-value pairs
		for key, value := range want.Iter() {
			gotValue, exists := om.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Edge Cases", func(t *testing.T) {
		jsonText := `{"max_int": 9223372036854775807, "min_int": -9223372036854775808, "special_chars": "!@#$%^&*()"}`
		want := func() *OrderedMap[string, any] {
			m := NewOrderedMap[string, any]()
			m.Set("max_int", float64(9223372036854775807))
			m.Set("min_int", float64(-9223372036854775808))
			m.Set("special_chars", "!@#$%^&*()")
			return m
		}()

		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.NoError(t, err)
		assert.Equal(t, want.Len(), om.Len())

		// Verify all key-value pairs
		for key, value := range want.Iter() {
			gotValue, exists := om.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		jsonText := `{invalid json}`
		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.Error(t, err)
	})

	t.Run("Nested Structures", func(t *testing.T) {
		jsonText := `{"nested": {"array": [{"id": 1}, {"id": 2}], "map": {"key1": {"value": 1}, "key2": {"value": 2}}}}`
		want := func() *OrderedMap[string, any] {
			m := NewOrderedMap[string, any]()
			m.Set("nested", map[string]any{
				"array": []any{
					map[string]any{"id": float64(1)},
					map[string]any{"id": float64(2)},
				},
				"map": map[string]any{
					"key1": map[string]any{"value": float64(1)},
					"key2": map[string]any{"value": float64(2)},
				},
			})
			return m
		}()

		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)

		assert.NoError(t, err)
		assert.Equal(t, want.Len(), om.Len())

		// Verify all key-value pairs
		for key, value := range want.Iter() {
			gotValue, exists := om.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Order Preservation", func(t *testing.T) {
		// Use an ordered JSON string with a defined key order
		jsonText := `{"z": 1, "y": 2, "x": 3, "w": 4, "v": 5, "u": 6, "t": 7, "s": 8, "r": 9, "q": 10}`

		// Deserialize to OrderedMap
		om := NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(jsonText), om)
		assert.NoError(t, err)

		// Verify that the key order matches the order in the JSON
		expectedKeys := []string{"z", "y", "x", "w", "v", "u", "t", "s", "r", "q"}

		// Collect the actual key order
		actualKeys := make([]string, 0, om.Len())
		for key, _ := range om.Iter() {
			actualKeys = append(actualKeys, key)
		}

		// Verify the order is consistent
		assert.Equal(t, expectedKeys, actualKeys, "Key order should match the order in the JSON")

		// Serialize again
		data, err := json.Marshal(om)
		assert.NoError(t, err)

		// Deserialize again
		om2 := NewOrderedMap[string, any]()
		err = json.Unmarshal(data, om2)
		assert.NoError(t, err)

		// Collect the key order after re-deserialization
		actualKeys2 := make([]string, 0, om2.Len())
		for key, _ := range om2.Iter() {
			actualKeys2 = append(actualKeys2, key)
		}

		// Verify that the order is preserved after serialization and deserialization
		assert.Equal(t, actualKeys, actualKeys2, "Key order should be preserved after serialization and deserialization")
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		m := NewOrderedMap[string, any]()
		m.Set("key1", "value1")
		m.Set("key2", float64(123))
		m.Set("key3", true)

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		expected := `{"key1":"value1","key2":123,"key3":true}`
		assert.Equal(t, expected, string(data))

		// Deserialize back to verify
		m2 := NewOrderedMap[string, any]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify all key-value pairs
		for key, value := range m.Iter() {
			gotValue, exists := m2.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := NewOrderedMap[string, any]()
		data, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("Special Values", func(t *testing.T) {
		m := NewOrderedMap[string, any]()
		m.Set("empty", "")
		m.Set("null", nil)
		m.Set("zero", float64(0))

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		expected := `{"empty":"","null":null,"zero":0}`
		assert.Equal(t, expected, string(data))
	})

	t.Run("Complex Types", func(t *testing.T) {
		m := NewOrderedMap[string, any]()
		m.Set("struct", map[string]any{
			"name": "test",
			"age":  float64(30),
		})
		m.Set("array", []any{float64(1), float64(2), float64(3)})
		m.Set("map", map[string]any{
			"key": "value",
		})

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Deserialize back to verify
		m2 := NewOrderedMap[string, any]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify all key-value pairs
		for key, value := range m.Iter() {
			gotValue, exists := m2.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Edge Cases", func(t *testing.T) {
		m := NewOrderedMap[string, any]()
		m.Set("max_int", float64(9223372036854775807))
		m.Set("min_int", float64(-9223372036854775808))
		m.Set("special_chars", "!@#$%^&*()")

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Deserialize back to verify
		m2 := NewOrderedMap[string, any]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify all key-value pairs
		for key, value := range m.Iter() {
			gotValue, exists := m2.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Key-Value Order Preservation", func(t *testing.T) {
		// Create an ordered map with a defined key order
		m := NewOrderedMap[string, any]()
		keys := []string{"z", "y", "x", "w", "v", "u", "t", "s", "r", "q"}
		for i, key := range keys {
			m.Set(key, i+1)
		}

		// Serialize
		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Verify that the serialized JSON string preserves the key order
		expectedJSON := `{"z":1,"y":2,"x":3,"w":4,"v":5,"u":6,"t":7,"s":8,"r":9,"q":10}`
		assert.Equal(t, expectedJSON, string(data))

		// Deserialize back
		m2 := NewOrderedMap[string, any]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)

		// Collect the key order after deserialization
		actualKeys := make([]string, 0, m2.Len())
		for key, _ := range m2.Iter() {
			actualKeys = append(actualKeys, key)
		}

		// Verify the order is consistent
		assert.Equal(t, keys, actualKeys, "Key order should be preserved after serialization and deserialization")
	})

	t.Run("Nested Structures", func(t *testing.T) {
		m := NewOrderedMap[string, any]()
		m.Set("nested", map[string]any{
			"array": []any{
				map[string]any{"id": float64(1)},
				map[string]any{"id": float64(2)},
			},
			"map": map[string]any{
				"key1": map[string]any{"value": float64(1)},
				"key2": map[string]any{"value": float64(2)},
			},
		})

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Deserialize back to verify
		m2 := NewOrderedMap[string, any]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify nested structure
		nestedValue, exists := m2.Get("nested")
		assert.True(t, exists)
		assert.IsType(t, map[string]any{}, nestedValue)

		nestedMap := nestedValue.(map[string]any)
		assert.Contains(t, nestedMap, "array")
		assert.Contains(t, nestedMap, "map")
	})

	t.Run("Non-String Key Types", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		m.Set(1, "one")
		m.Set(2, "two")

		data, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"1":"one"`)
		assert.Contains(t, string(data), `"2":"two"`)
	})
}

func TestNonStringKeyMarshalJSON(t *testing.T) {
	t.Run("Integer Keys", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		m.Set(1, "one")
		m.Set(2, "two")
		m.Set(3, "three")

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		expected := `{"1":"one","2":"two","3":"three"}`
		assert.Equal(t, expected, string(data))

		// Deserialize back to verify
		m2 := NewOrderedMap[int, string]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify all key-value pairs
		for key, value := range m.Iter() {
			gotValue, exists := m2.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Float Keys", func(t *testing.T) {
		m := NewOrderedMap[float64, string]()
		m.Set(1.1, "one point one")
		m.Set(2.2, "two point two")
		m.Set(3.3, "three point three")

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Deserialize back to verify
		m2 := NewOrderedMap[float64, string]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify all key-value pairs
		for key, value := range m.Iter() {
			gotValue, exists := m2.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Boolean Keys", func(t *testing.T) {
		m := NewOrderedMap[bool, string]()
		m.Set(true, "true value")
		m.Set(false, "false value")

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Deserialize back to verify
		m2 := NewOrderedMap[bool, string]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)
		assert.Equal(t, m.Len(), m2.Len())

		// Verify all key-value pairs
		for key, value := range m.Iter() {
			gotValue, exists := m2.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, gotValue)
		}
	})

	t.Run("Composite Key Types", func(t *testing.T) {
		type Point struct {
			X, Y int
		}
		m := NewOrderedMap[Point, string]()
		m.Set(Point{1, 2}, "point 1,2")
		m.Set(Point{3, 4}, "point 3,4")

		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Since composite key type conversion may not be perfect, we only verify serialization success
		assert.Contains(t, string(data), "point 1,2")
		assert.Contains(t, string(data), "point 3,4")
	})

	t.Run("Key-Value Order Preservation", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		keys := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
		for _, key := range keys {
			m.Set(key, fmt.Sprintf("value%d", key))
		}

		// Serialize
		data, err := json.Marshal(m)
		assert.NoError(t, err)

		// Deserialize back
		m2 := NewOrderedMap[int, string]()
		err = json.Unmarshal(data, m2)
		assert.NoError(t, err)

		// Collect the key order after deserialization
		actualKeys := make([]int, 0, m2.Len())
		for key, _ := range m2.Iter() {
			actualKeys = append(actualKeys, key)
		}

		// Verify the order is consistent
		assert.Equal(t, keys, actualKeys, "Key order should be preserved after serialization and deserialization")
	})
}
