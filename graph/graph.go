package types

import (
	"xcontainer"
	"xcontainer/list"
	"xcontainer/queue"
	"xcontainer/stack"
)

// UndirectedGraph stored by adjacency list
type UndirectedGraph[T comparable] struct {
	v   int
	adj map[T]*list.SingleListSet[T]
}

func NewUndirectedGraph[T comparable](v int) *UndirectedGraph[T] {
	g := &UndirectedGraph[T]{
		v:   v,
		adj: make(map[T]*list.SingleListSet[T]),
	}
	return g
}

func (g *UndirectedGraph[T]) AddEdge(s, t T) {
	if _, ok := g.adj[s]; !ok {
		g.adj[s] = list.NewSingleListSet[T]()
	}
	if _, ok := g.adj[t]; !ok {
		g.adj[t] = list.NewSingleListSet[T]()
	}
	g.adj[s].Add(t)
	g.adj[t].Add(s)
}

func UndirectedGraphBFS[T comparable](g *UndirectedGraph[T], start T) []T {
	ret := make([]T, 0)
	q := queue.NewQueue[T]()
	visited := make(map[T]bool)
	q.EnQueue(start)
	visited[start] = true

	for q.Len() > 0 {
		vertex := q.DeQueue()
		l := g.adj[vertex]
		for node := l.Head(); !xcontainer.IsNil[*list.SingleListNode[T]](node); node = node.Next {
			if v, ok := visited[node.Val]; v && ok {
				continue
			}
			q.EnQueue(node.Val)
			visited[node.Val] = true
		}
		ret = append(ret, vertex)
	}
	return ret
}

func UndirectedGraphDFS[T comparable](g *UndirectedGraph[T], start T) []T {
	ret := make([]T, 0)
	s := stack.NewStack[T]()
	visited := make(map[T]bool)
	s.Push(start)
	visited[start] = true

	for s.Len() > 0 {
		vertex := s.Pop()
		l := g.adj[vertex]
		for node := l.Head(); !xcontainer.IsNil[*list.SingleListNode[T]](node); node = node.Next {
			if v, ok := visited[node.Val]; v && ok {
				continue
			}
			s.Push(node.Val)
			visited[node.Val] = true
		}
		ret = append(ret, vertex)
	}
	return ret
}
