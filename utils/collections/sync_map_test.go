package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSyncMap_Base(t *testing.T) {
	should := require.New(t)
	sMap := &SyncMap[int, string]{}
	// 存储
	sMap.Store(1, "1")
	sMap.Store(2, "2")
	sMap.Store(3, "3")
	sMap.Store(4, "3")
	sMap.Store(5, "5")
	sMap.Store(5, "6")
	// 读取
	value2, ok2 := sMap.Load(2)
	should.True(value2 == "2" && ok2)
	// 删除
	sMap.Delete(3)
	// 读取不存在的值
	value3, ok3 := sMap.Load(3)
	t.Logf("key=3,value=%s, ok=%v", value3, ok3)
	value6, ok6 := sMap.Load(6)
	t.Logf("key=6,value=%s, ok=%v", value6, ok6)
	// 遍历
	sMap.Range(func(key int, value string) bool {
		t.Logf("--key=%v, value=%v", key, value)
		return true
	})
}

func TestSyncMap_LoadAndStore(t *testing.T) {
	should := require.New(t)
	sMap := &SyncMap[int, string]{}
	var actual string
	var loaded bool
	actual, loaded = sMap.LoadOrStore(1, "1")
	should.True(actual == "" && loaded == false)
	actual, loaded = sMap.LoadOrStore(1, "1")
	should.True(actual == "1" && loaded == true)
}
