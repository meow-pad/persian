package collections

import (
	"sync"
)

type SyncMap[K, V any] struct {
	syncMap sync.Map
}

func (sMap *SyncMap[K, V]) Store(key K, value V) {
	sMap.syncMap.Store(key, value)
}

func (sMap *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	sActual, sLoaded := sMap.syncMap.LoadOrStore(key, value)
	if !sLoaded {
		actual = value
		loaded = false
		return
	}
	if sActual != nil {
		actual = sActual.(V)
	}
	loaded = true
	return
}

func (sMap *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	sValue, sOk := sMap.syncMap.Load(key)
	if !sOk {
		ok = false
		return
	}
	if sValue != nil {
		value = sValue.(V)
	}
	ok = true
	return
}

func (sMap *SyncMap[K, V]) Delete(key K) (value V, loaded bool) {
	sValue, sLoaded := sMap.syncMap.LoadAndDelete(key)
	if !sLoaded {
		loaded = false
		return
	}
	if sValue != nil {
		value = sValue.(V)
	}
	loaded = true
	return
}

// Range
//
//	@Description: 使用函数遍历访问集合元素
//	@receiver sMap *SyncMap[K, V]
//	@param f func(key K, value V) bool 元素处理函数，函数返回值为false时，打断遍历
func (sMap *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	sMap.syncMap.Range(func(sKey, sValue any) bool {
		key := sKey.(K)
		value := sValue.(V)
		return f(key, value)
	})
}

func (sMap *SyncMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	sMap.syncMap.Range(func(sKey, sValue any) bool {
		key := sKey.(K)
		keys = append(keys, key)
		return true
	})
	return keys
}

func (sMap *SyncMap[K, V]) Values() []V {
	values := make([]V, 0)
	sMap.syncMap.Range(func(sKey, sValue any) bool {
		value := sValue.(V)
		values = append(values, value)
		return true
	})
	return values
}
