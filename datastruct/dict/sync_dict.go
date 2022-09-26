package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {

	val, exists = dict.m.Load(key)
	return val, exists
}

func (dict *SyncDict) Len() int {
	result := 0

	dict.m.Range(func(key, value interface{}) bool {

		result++
		return true

	})

	return result
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {

	val, exists := dict.m.Load(key)
	dict.m.Store(key, val)
	if exists {
		// key存在只是更新
		return 0
	}

	// 新插入
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	val, exists := dict.m.Load(key)
	if !exists {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	val, exists := dict.m.Load(key)

	if exists {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exists := dict.m.Load(key)
	if exists {
		dict.m.Delete(key)
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	keys := make([]string, dict.Len())
	index := 0

	dict.m.Range(func(key, value interface{}) bool {

		keys[index] = key.(string)
		index++
		return true

	})
	return keys
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	keys := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value interface{}) bool {
			keys[i] = key.(string)
			return false
		})

	}
	return keys
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	keys := make([]string, limit)
	index := 0
	dict.m.Range(func(key, value interface{}) bool {
		keys[index] = key.(string)
		index++
		if index == limit {
			return false
		}
		return true
	})
	return keys
}

func (dict *SyncDict) Clear() {
	*dict = *MakeSyncDict()
}
