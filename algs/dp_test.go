package algs

import "testing"

func TestEditDistanceDP(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"kitten", "sitting", 3},
		{"", "", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"abc", "abc", 0},
		{"abc", "def", 3},
	}

	for _, test := range tests {
		actual := editDistanceDP(test.s1, test.s2)
		if actual != test.expected {
			t.Errorf("editDistanceDP(%q, %q) = %d; expected %d", test.s1, test.s2, actual, test.expected)
		}
	}
}
