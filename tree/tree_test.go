package tree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreOrder(t *testing.T) {
	l1 := []string{"6"}
	l2 := []string{"1", "3"}
	l3 := []string{"9", "4", "2", "7"}
	l4 := []string{"5", "nil", "nil", "nil", "8", "nil", "nil", "nil"}
	root := NewTree[int](l1, l2, l3, l4)
	rlt := make([]int, 0)
	PreOrder[int](root, &rlt)
	assert.EqualValues(t, []int{6, 1, 9, 5, 4, 3, 2, 8, 7}, rlt)
}

func TestInOrder(t *testing.T) {
	l1 := []string{"6"}
	l2 := []string{"1", "3"}
	l3 := []string{"9", "4", "2", "7"}
	l4 := []string{"5", "nil", "nil", "nil", "8", "nil", "nil", "nil"}
	root := NewTree[int](l1, l2, l3, l4)
	rlt := make([]int, 0)
	InOrder[int](root, &rlt)
	assert.EqualValues(t, []int{5, 9, 1, 4, 6, 8, 2, 3, 7}, rlt)
}

func TestPostOrder(t *testing.T) {
	l1 := []string{"6"}
	l2 := []string{"1", "3"}
	l3 := []string{"9", "4", "2", "7"}
	l4 := []string{"5", "nil", "nil", "nil", "8", "nil", "nil", "nil"}
	root := NewTree[int](l1, l2, l3, l4)
	rlt := make([]int, 0)
	PostOrder[int](root, &rlt)
	assert.EqualValues(t, []int{5, 9, 4, 1, 8, 2, 7, 3, 6}, rlt)
}

func TestTreeBFS(t *testing.T) {
	l1 := []string{"6"}
	l2 := []string{"1", "3"}
	l3 := []string{"9", "4", "2", "7"}
	l4 := []string{"5", "nil", "nil", "nil", "8", "nil", "nil", "nil"}
	root := NewTree[int](l1, l2, l3, l4)
	rlt := TreeBFS[int](root)
	assert.EqualValues(t, [][]int{
		{6},
		{1, 3},
		{9, 4, 2, 7},
		{5, 8},
	}, rlt)
}

func TestSearchTree(t *testing.T) {
	st := &SearchTree[int]{}
	src := []int{5, 9, 4, 1, 8, 2, 7, 3, 6}
	for _, v := range src {
		st.Put(v)
	}
	n := st.Find(9)
	assert.EqualValues(t, 9, n.Val())
	inOrderRlt := make([]int, 0)
	InOrder[int](st.root, &inOrderRlt)
	assert.EqualValues(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, inOrderRlt)
}

func TestSearchTreeDel(t *testing.T) {
	st := &SearchTree[int]{}
	src := []int{5, 9, 4, 1, 8, 2, 7, 3, 6, 10}
	for _, v := range src {
		st.Put(v)
	}
	st.Del(6)
	st.Del(9)
	st.Del(4)
	inOrderRlt := make([]int, 0)
	InOrder[int](st.root, &inOrderRlt)
	assert.EqualValues(t, []int{1, 2, 3, 5, 7, 8, 10}, inOrderRlt)
	st.Del(5)
	inOrderRlt = make([]int, 0)
	InOrder[int](st.root, &inOrderRlt)
	assert.EqualValues(t, []int{1, 2, 3, 7, 8, 10}, inOrderRlt)

	st = &SearchTree[int]{}
	src = []int{1, 2, 3}
	for _, v := range src {
		st.Put(v)
	}
	st.Del(1)
	inOrderRlt = make([]int, 0)
	InOrder[int](st.root, &inOrderRlt)
	assert.EqualValues(t, []int{2, 3}, inOrderRlt)
	st.Del(2)
	inOrderRlt = make([]int, 0)
	InOrder[int](st.root, &inOrderRlt)
	assert.EqualValues(t, []int{3}, inOrderRlt)
	st.Del(3)
	inOrderRlt = make([]int, 0)
	InOrder[int](st.root, &inOrderRlt)
	assert.EqualValues(t, []int{}, inOrderRlt)
}

func TestParseFuncInt(t *testing.T) {
	intParser := ParseFunc[int]()
	v, err := intParser("123")
	assert.Nil(t, err)
	assert.Equal(t, 123, v)
}

func TestParseFuncString(t *testing.T) {
	intParser := ParseFunc[string]()
	v, err := intParser("nil")
	assert.Equal(t, errors.New("nil"), err)
	assert.Equal(t, "nil", v)
}
