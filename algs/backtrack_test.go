package algs

import "testing"

func TestNQueens_Zero(t *testing.T) {
	result := nQueens(0)
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %v", result)
	}
}

func TestNQueens(t *testing.T) {
	tests := []struct {
		n      int
		result [][]int
	}{
		{
			n: 1,
			result: [][]int{
				{0},
			},
		},
		{
			n:      2,
			result: [][]int{},
		},
		{
			n:      3,
			result: [][]int{},
		},
		{
			n: 4,
			result: [][]int{
				{1, 3, 0, 2},
				{2, 0, 3, 1},
			},
		},
		{
			n: 5,
			result: [][]int{
				{0, 2, 4, 1, 3},
				{0, 3, 1, 4, 2},
				{1, 3, 0, 2, 4},
				{1, 4, 2, 0, 3},
				{2, 0, 3, 1, 4},
				{2, 4, 1, 3, 0},
				{3, 0, 2, 4, 1},
				{3, 1, 4, 2, 0},
				{4, 1, 3, 0, 2},
				{4, 2, 0, 3, 1},
			},
		},
	}

	for _, test := range tests {
		actual := nQueens(test.n)
		if len(actual) != len(test.result) {
			t.Errorf("Expected %v, got %v", test.result, actual)
		} else {
			for i := range actual {
				if len(actual[i]) != len(test.result[i]) {
					t.Errorf("Expected %v, got %v", test.result, actual)
				} else {
					for j := range actual[i] {
						if actual[i][j] != test.result[i][j] {
							t.Errorf("Expected %v, got %v", test.result, actual)
						}
					}
				}
			}
		}
	}
}

func TestNQueens_Large(t *testing.T) {
	n := 10
	result := nQueens(n)
	if len(result) != 724 {
		t.Errorf("Expected 724 solutions for n=10, got %d", len(result))
	}
}
