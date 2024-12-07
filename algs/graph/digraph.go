package graph

import (
	"errors"
	"fmt"

	"github.com/realrabbithouse/go-play/comparable"
)

var ErrVertexOutOfRange = errors.New("vertex out of range")

type Digraph struct {
	v        int
	e        int
	adj      [][]int
	indegree []int // indegree[v] = number of edges pointing to v
}

func NewDigraph(v int) (*Digraph, error) {
	if v < 0 {
		return nil, fmt.Errorf("number of vertices %d is negative", v)
	}
	return &Digraph{
		v:        v,
		adj:      make([][]int, v),
		indegree: make([]int, v),
	}, nil
}

func (g *Digraph) validateVertex(v int) error {
	if v < 0 || v >= g.v {
		return fmt.Errorf("vertex %d is not between 0 and %d: %w", v, g.v-1, ErrVertexOutOfRange)
	}
	return nil
}

func (g *Digraph) AddEdge(v, w int) error {
	if err := g.validateVertex(v); err != nil {
		return err
	}
	if err := g.validateVertex(w); err != nil {
		return err
	}
	g.adj[v] = append(g.adj[v], w)
	g.indegree[w]++
	g.e++
	return nil
}

func (g *Digraph) Adj(v int) ([]int, error) {
	if err := g.validateVertex(v); err != nil {
		return nil, err
	}
	return g.adj[v], nil
}

func (g *Digraph) V() int {
	return g.v
}

func (g *Digraph) E() int {
	return g.e
}

// *************************************************************** //

type DirectedEdge struct {
	from   int
	to     int
	weight comparable.Comparable
}

type EdgeWeightedDigraph struct {
	v        int
	e        int
	adj      [][]DirectedEdge
	indegree []int // indegree[v] = number of edges pointing to v
}
