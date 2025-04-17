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

func AppendIfNotInSlice[T comparable](slice []T, value T) []T {
	for _, sv := range slice {
		if sv == value {
			return slice
		}
	}
	return append(slice, value)
}

func RemoveFromSlice[T comparable](slice []T, value T) []T {
	for i, sv := range slice {
		if sv == value {
			return append(slice[0:i], slice[i+1:]...)
		}
	}
	return slice
}

func HasSameElements[T comparable](slice1 []T, slice2 []T) bool {
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
//	@param compare 比较函数，插入值与其他值的比较，大于为1（正数），小于为-1（负数），相等为0
//	@return int 插入的位置
func SortedSliceSearchInsert[T any](array []T, value T, compare func(v1, v2 T) int) int {
	var mid int
	low, high := 0, len(array)-1
	if high < 0 {
		return 0
	}
	if compare(value, array[high]) >= 0 {
		return high + 1
	}
	for low <= high {
		mid = (low + high) / 2
		i := compare(value, array[mid])
		if i > 0 {
			low = mid + 1
		} else if i < 0 {
			high = mid - 1
		} else {
			return mid
		}
	}
	return low
}

// SortedSliceSearch
//
//	@Description: 二分查找元素在数组中的位置
//	@param array
//	@param value
//	@param compare 比较函数，插入值与其他值的比较，大于为1（正数），小于为-1（负数），相等为0
//	@param equal 判断是否相等的函数
//	@return int 相等元素所在位置
//	@return int 可插入位置
func SortedSliceSearch[T any](array []T, value T, compare func(v1, v2 T) int, equal func(v1, v2 T) bool) (int, int) {
	var mid int
	low, high := 0, len(array)-1
	if high < 0 {
		return -1, 0
	}
	for low <= high {
		mid = (low + high) / 2
		i := compare(value, array[mid])
		if i > 0 {
			low = mid + 1
		} else if i < 0 {
			high = mid - 1
		} else {
			if equal(value, array[mid]) {
				return mid, mid
			}
			// 前向查找是否有相等的值
			for j := mid - 1; j >= 0; j-- {
				if compare(value, array[j]) != 0 {
					break
				}
				if equal(value, array[j]) {
					return j, j
				}
			}
			// 后向查找是否有相等的值
			for j := mid + 1; j < len(array); j++ {
				if compare(value, array[j]) != 0 {
					break
				}
				if equal(value, array[j]) {
					return j, j
				}
			}
			return -1, mid
		}
	}
	return -1, low
}
