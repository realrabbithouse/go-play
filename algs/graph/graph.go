package graph

import (
	"github.com/realrabbithouse/go-play/comparable"
)

type Graph struct {
	v   int
	e   int
	adj [][]int
}

// *************************************************************** //

type Edge struct {
	v      int
	w      int
	weight comparable.Comparable
}

type EdgeWeightedGraph struct {
	v   int
	e   int
	adj [][]Edge
}
