package sortutil

import "container/heap"

type _heapElement[T any] struct {
	value    T   // 当前元素值
	arrayIdx int // 所属数组的索引
	elemIdx  int // 在数组中的索引
}

// _minHeap 定义最小堆类型
type _minHeap[T any] struct {
	elems []_heapElement[T]
	less  func(o1, o2 T) bool
}

// Len 实现 heap.Interface 接口
func (h *_minHeap[T]) Len() int { return len(h.elems) }
func (h *_minHeap[T]) Less(i, j int) bool {
	return h.less(h.elems[i].value, h.elems[j].value)
}                                    // 使用比较函数
func (h *_minHeap[T]) Swap(i, j int) { h.elems[i], h.elems[j] = h.elems[j], h.elems[i] }

func (h *_minHeap[T]) Push(x interface{}) {
	h.elems = append(h.elems, x.(_heapElement[T]))
}

func (h *_minHeap[T]) Pop() interface{} {
	old := h.elems
	n := len(old)
	x := old[n-1]
	h.elems = old[:n-1]
	return x
}

// MergeSortedArrays
//
//	@Description: 将多个有序数组归并成一个有序数组
//	@param arrays
//	@param less
//	@return []T
func MergeSortedArrays[T any](arrays [][]T, less func(o1, o2 T) bool) []T {
	if len(arrays) <= 0 {
		return nil
	}
	if len(arrays) == 1 {
		return arrays[0]
	}
	// 初始化最小堆
	h := &_minHeap[T]{less: less}
	heap.Init(h)

	// 将每个数组的第一个元素加入堆
	for i, arr := range arrays {
		if len(arr) > 0 {
			heap.Push(h, _heapElement[T]{
				value:    arr[0],
				arrayIdx: i,
				elemIdx:  0,
			})
		} // end of if
	} // end of for

	// 结果数组
	var result []T

	// 归并过程
	for h.Len() > 0 {
		// 取出堆顶元素（当前最小值）
		minElem := heap.Pop(h).(_heapElement[T])
		result = append(result, minElem.value)

		// 如果该元素所属的数组还有剩余元素，则将下一个元素加入堆
		if minElem.elemIdx+1 < len(arrays[minElem.arrayIdx]) {
			heap.Push(h, _heapElement[T]{
				value:    arrays[minElem.arrayIdx][minElem.elemIdx+1],
				arrayIdx: minElem.arrayIdx,
				elemIdx:  minElem.elemIdx + 1,
			})
		} // end of if
	} // end of for

	return result
}
