package xmap

import (
	"iter"

	"github.com/danielhookx/xcontainer/list"
)

type MapNode[K comparable, V any] struct {
	K K
	V V
}

type OrderedMap[K comparable, V any] struct {
	kv map[K]*list.Element[*MapNode[K, V]]
	l  *list.List[*MapNode[K, V]]
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*list.Element[*MapNode[K, V]]),
		l:  list.New[*MapNode[K, V]](),
	}
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	e, ok := m.kv[key]
	if ok {
		return e.Value.V, true
	}
	var v V
	return v, false
}

func (m *OrderedMap[K, V]) Set(key K, value V) bool {
	_, ok := m.kv[key]
	if ok {
		m.kv[key].Value.V = value
		return false
	}

	e := m.l.PushBack(&MapNode[K, V]{K: key, V: value})
	m.kv[key] = e
	return true
}

func (m *OrderedMap[K, V]) Delete(key K) bool {
	e, ok := m.kv[key]
	if ok {
		m.l.Remove(e)
		delete(m.kv, key)
	}
	return ok
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.kv)
}

func (m *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m2 := NewOrderedMap[K, V]()

	for e := m.l.Front(); e != nil; e = e.Next() {
		m2.Set(e.Value.K, e.Value.V)
	}

	return m2
}

func (m *OrderedMap[K, V]) ToArray() []V {
	array := make([]V, 0, m.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		array = append(array, e.Value.V)
	}
	return array
}

func (m *OrderedMap[K, V]) Iter() iter.Seq2[K, V] {
	// 在迭代开始时创建一个快照
	snapshot := make([]*MapNode[K, V], 0, m.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		snapshot = append(snapshot, e.Value)
	}

	return func(yield func(K, V) bool) {
		// 使用快照进行迭代
		for _, node := range snapshot {
			if !yield(node.K, node.V) {
				return
			}
		}
	}
}
