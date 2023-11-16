package types

import (
	"github.com/danielhookx/xcontainer/list"
	"github.com/danielhookx/xcontainer/queue"
	"github.com/danielhookx/xcontainer/stack"
)

// UndirectedGraph stored by adjacency list
type UndirectedGraph[T comparable] struct {
	v   int
	adj map[T]*list.SetList[T]
}

func NewUndirectedGraph[T comparable](v int) *UndirectedGraph[T] {
	g := &UndirectedGraph[T]{
		v:   v,
		adj: make(map[T]*list.SetList[T]),
	}
	return g
}

func (g *UndirectedGraph[T]) AddEdge(s, t T) {
	if _, ok := g.adj[s]; !ok {
		g.adj[s] = list.NewSetList[T]()
	}
	if _, ok := g.adj[t]; !ok {
		g.adj[t] = list.NewSetList[T]()
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
		iterate, cancel := l.Iterate()
		defer cancel()
		for item, ok := iterate(); ok; item, ok = iterate() {
			if v, ok := visited[item]; v && ok {
				continue
			}
			q.EnQueue(item)
			visited[item] = true
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
		iterate, cancel := l.Iterate()
		defer cancel()
		for item, ok := iterate(); ok; item, ok = iterate() {
			if v, ok := visited[item]; v && ok {
				continue
			}
			s.Push(item)
			visited[item] = true
		}
		ret = append(ret, vertex)
	}
	return ret
}
