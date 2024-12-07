package algs

import (
	"errors"
	"fmt"
)

var ErrIndexOutOfRange = errors.New("index out of range")

type UF struct {
	parent []int
	rank   []uint8 // rank[i] = rank of subtree rooted at i (never more than 31)
	count  int     // number of unconnected components
}

func NewUF(n int) (*UF, error) {
	if n < 0 {
		return nil, errors.New("n must be positive")
	}
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return &UF{
		parent: parent,
		rank:   make([]uint8, n),
		count:  n,
	}, nil
}

func (u *UF) validate(p int) error {
	n := len(u.parent)
	if p < 0 || p >= n {
		return fmt.Errorf("index %d is not between 0 and %d: %w", p, n-1, ErrIndexOutOfRange)
	}
	return nil
}

// Count returns the number of connected components.
func (u *UF) Count() int {
	return u.count
}

func (u *UF) Connected(p, q int) (bool, error) {
	rootP, err := u.find(p)
	if err != nil {
		return false, err
	}
	rootQ, err := u.find(q)
	if err != nil {
		return false, err
	}

	return rootP == rootQ, nil
}

func (u *UF) Union(p, q int) error {
	rootP, err := u.find(p)
	if err != nil {
		return err
	}
	rootQ, err := u.find(q)
	if err != nil {
		return err
	}

	if rootP == rootQ {
		// already connected
		return nil
	}

	if u.rank[rootP] < u.rank[rootQ] {
		u.parent[rootP] = rootQ
	} else if u.rank[rootP] > u.rank[rootQ] {
		u.parent[rootQ] = rootP
	} else {
		u.parent[rootQ] = rootP
		u.rank[rootP]++
	}
	u.count--

	return nil
}

func (u *UF) find(p int) (int, error) {
	if err := u.validate(p); err != nil {
		return 0, err
	}
	for p != u.parent[p] {
		u.parent[p] = u.parent[u.parent[p]] // path compression by halving
		p = u.parent[p]
	}
	return p, nil
}
