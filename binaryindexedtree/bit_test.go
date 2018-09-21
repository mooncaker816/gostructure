package bit

import (
	"fmt"
	"testing"
)

func TestNewBIT(t *testing.T) {
	a := []int{1, 7, 3, 0, 5, 8, 3, 2, 6, 2, 1, 1, 4, 5}
	bit := NewBIT(a)
	if fmt.Sprintf("%v", bit) != "[0 1 8 3 11 5 13 3 29 6 8 1 10 4 9]" {
		t.Errorf("not ok")
	}
}

func TestSum(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	var tt = []struct {
		prefix int
		sum    int
	}{
		{0, 1},
		{1, 3},
		{2, 6},
		{3, 10},
		{4, 15},
		{5, 21},
		{6, 28},
	}
	bit := NewBIT(a)
	for _, tc := range tt {
		if sum := bit.Sum(tc.prefix); sum != tc.sum {
			t.Errorf("failed for prefix %d, expected %d, got %d\n", tc.prefix, tc.sum, sum)
		}
	}
}
