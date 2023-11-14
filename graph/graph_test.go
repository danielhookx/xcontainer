package types

import (
	"testing"

	"github.com/danielhookx/xcontainer/set"

	"github.com/stretchr/testify/assert"
)

func TestUndirectedGraphBFS(t *testing.T) {
	A := 'A'
	B := 'B'
	C := 'C'
	D := 'D'
	E := 'E'
	F := 'F'

	g := NewUndirectedGraph[rune](6)
	g.AddEdge(A, B)
	g.AddEdge(A, C)
	g.AddEdge(B, A)
	g.AddEdge(B, C)
	g.AddEdge(B, D)
	g.AddEdge(C, A)
	g.AddEdge(C, B)
	g.AddEdge(C, D)
	g.AddEdge(C, E)
	g.AddEdge(D, C)
	g.AddEdge(D, E)
	g.AddEdge(D, F)
	g.AddEdge(E, C)
	g.AddEdge(E, D)
	g.AddEdge(F, D)

	ret := UndirectedGraphBFS(g, A)
	assert.True(t, set.BuildSet[rune]([]rune{A, B, C, D, E, F}...).Equal(set.BuildSet[rune](ret...)))
}

func TestUndirectedGraphDFS(t *testing.T) {
	A := 'A'
	B := 'B'
	C := 'C'
	D := 'D'
	E := 'E'
	F := 'F'

	g := NewUndirectedGraph[rune](6)
	g.AddEdge(A, B)
	g.AddEdge(A, C)
	g.AddEdge(B, A)
	g.AddEdge(B, C)
	g.AddEdge(B, D)
	g.AddEdge(C, A)
	g.AddEdge(C, B)
	g.AddEdge(C, D)
	g.AddEdge(C, E)
	g.AddEdge(D, C)
	g.AddEdge(D, E)
	g.AddEdge(D, F)
	g.AddEdge(E, C)
	g.AddEdge(E, D)
	g.AddEdge(F, D)

	ret := UndirectedGraphDFS(g, A)
	assert.True(t, set.BuildSet[rune]([]rune{A, B, C, D, E, F}...).Equal(set.BuildSet[rune](ret...)))
}
