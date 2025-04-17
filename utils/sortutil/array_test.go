package sortutil

import "testing"

func TestMergeSortedArrays(t *testing.T) {
	arrays := [][]int{
		{1, 3, 5},
		{2, 4, 6},
		{0, 7, 8},
	}
	res := MergeSortedArrays(arrays, func(o1, o2 int) bool {
		return o1 < o2
	})
	t.Log(res)
}
