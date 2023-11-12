package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndCopy(t *testing.T) {
	a := []int{1, 2, 2, 3, 4, 5}
	c1 := BuildSet[int](a...)

	c2 := CreateSet[int]()
	iterate, cancel := c1.Iterate()
	defer cancel()
	for item, ok := iterate(); ok; item, ok = iterate() {
		c2.Add(item)
	}
	assert.True(t, c1.Equal(c2))

	c3 := CreateFrom[int](c1)
	assert.True(t, c1.Equal(c3))

	for _, item := range c1.Enumerate() {
		assert.True(t, c2.IsElementOf(item))
	}
	c1.Clear()
	assert.True(t, IsEmptySet[int](c1))
}

func TestFilter(t *testing.T) {
	a := []int{1, 2, 2, 3, 4, 5}
	c1 := BuildSet[int](a...)

	c2 := c1.Filter(func(t int) bool {
		return t > 3
	})
	assert.True(t, BuildSet[int]([]int{4, 5}...).Equal(c2))
}

func TestMap(t *testing.T) {
	a := []int{1, 2, 2, 3, 4, 5}
	c1 := BuildSet[int](a...)

	c2 := c1.Map(func(t int) int {
		return t * 2
	})
	assert.True(t, BuildSet[int]([]int{2, 4, 6, 8, 10}...).Equal(c2))
}

func TestCreateSetWithCapacity(t *testing.T) {
	c1 := CreateSetWithCapacity[int](0)
	err := c1.Add(1)
	assert.Nil(t, err)
	err = c1.Add(2)
	assert.Nil(t, err)
	assert.Equal(t, 2, c1.Size())

	c2 := CreateSetWithCapacity[int](1)
	err = c2.Add(1)
	assert.Nil(t, err)
	err = c2.Add(2)
	assert.NotNil(t, err)
	assert.Equal(t, 1, c2.Size())
}
