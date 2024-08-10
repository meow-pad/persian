package collections

import "github.com/meow-pad/persian/utils/rand"

func IsInSlice[T comparable](slice []T, value T) bool {
	for _, sv := range slice {
		if sv == value {
			return true
		}
	}
	return false
}

func IsSameSlice[T comparable](slice1 []T, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, sv := range slice1 {
		if sv != slice2[i] {
			return false
		}
	}
	return true
}

func SliceSearchIndex[T comparable](slice []T, value T) int {
	for index, sv := range slice {
		if sv == value {
			return index
		}
	}
	return -1
}

// SelectFromSlice
//
//	@Description: 从切片中随机选择指定数量的元素
//	@param slice
//	@param num
//	@return []T
func SelectFromSlice[T any](slice []T, num int) []T {
	if num <= 0 {
		return nil
	}
	if num >= len(slice) {
		return slice
	}
	var lenSlice = int32(len(slice))
	var selectedIndexes []int
	var selectedValues []T
	for i := 0; i < num; i++ {
		index := rand.Int32n(lenSlice)
		for j := int32(0); j < lenSlice; j++ {
			if IsInSlice[int](selectedIndexes, int(index)) {
				index += 1
			} else {
				selectedIndexes = append(selectedIndexes, int(index))
				break
			}
		}
	}
	for _, index := range selectedIndexes {
		selectedValues = append(selectedValues, slice[index])
	}
	return selectedValues
}

// SortedSliceSearchInsert
//
//	@Description: 查找在有序数组中的插入位置
//	@param array
//	@param value 插入值
//	@param compare 比较函数，插入值与其他值的比较，大于为1，小于为-1，相等为0
//	@return int 插入的位置
func SortedSliceSearchInsert[T any](array []T, value T, compare func(v1, v2 T) int) int {
	var mid int
	low, height := 0, len(array)-1
	if height < 0 {
		return 0
	}
	if compare(value, array[height]) >= 0 {
		return height + 1
	}
	for low <= height {
		mid = (low + height) / 2
		i := compare(value, array[mid])
		if i > 0 {
			low = mid + 1
		} else if i < 0 {
			height = mid - 1
		} else {
			return mid
		}
	}
	return low
}
