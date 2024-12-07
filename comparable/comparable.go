package comparable

import "strings"

type Comparable interface {
	CompareTo(other Comparable) int
}

type Int int

func (i Int) CompareTo(other Comparable) int {
	delta := i - other.(Int)
	if delta < 0 {
		return -1
	} else if delta > 0 {
		return 1
	} else {
		return 0
	}
}

type Float64 float64

func (f Float64) CompareTo(other Comparable) int {
	delta := f - other.(Float64)
	if delta < 0 {
		return -1
	} else if delta > 0 {
		return 1
	} else {
		return 0
	}
}

type String string

func (s String) CompareTo(other Comparable) int {
	return strings.Compare(string(s), string(other.(String)))
}
