package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var _testCompare1 = func(v1, v2 int) int {
	if v1 > v2 {
		return 1
	} else if v1 < v2 {
		return -1
	} else {
		return 0
	}
}

var _testCompare2 = func(v1, v2 int) int {
	return (v1 / 10) - (v2 / 10)
}

var _testEqual = func(v1, v2 int) bool {
	return v1 == v2
}

func TestSortedSliceSearch1(t *testing.T) {
	should := require.New(t)
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	index, _ := SortedSliceSearch[int](array, 5, _testCompare1, _testEqual)
	should.Equal(4, index)
	index, _ = SortedSliceSearch[int](array, 13, _testCompare1, _testEqual)
	should.Equal(12, index)
	index, _ = SortedSliceSearch[int](array, 1, _testCompare1, _testEqual)
	should.Equal(0, index)
	index, _ = SortedSliceSearch[int](array, 15, _testCompare1, _testEqual)
	should.Equal(14, index)
	index, _ = SortedSliceSearch[int](array, 30, _testCompare1, _testEqual)
	should.Equal(-1, index)
	index, _ = SortedSliceSearch[int](nil, 30, _testCompare1, _testEqual)
	should.Equal(-1, index)
}

func TestSortedSliceSearch2(t *testing.T) {
	should := require.New(t)
	array := []int{22, 31, 300, 401, 403, 404, 402, 405, 513, 614, 715}
	index, insertIndex := SortedSliceSearch[int](array, 407, _testCompare2, _testEqual)
	should.Equal(-1, index)
	t.Logf("insertIndex: %v", insertIndex)
	index, _ = SortedSliceSearch[int](array, 404, _testCompare2, _testEqual)
	should.Equal(5, index)
	index, _ = SortedSliceSearch[int](array, 401, _testCompare2, _testEqual)
	should.Equal(3, index)
	index, insertIndex = SortedSliceSearch[int](array, 524, _testCompare2, _testEqual)
	should.Equal(-1, index)
	should.Equal(9, insertIndex)
	index, insertIndex = SortedSliceSearch[int](array, 11, _testCompare2, _testEqual)
	should.Equal(-1, index)
	should.Equal(0, insertIndex)
}
