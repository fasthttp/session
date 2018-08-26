package session

import (
	"sync"
)

var (
	// default slice number 32
	defaultSliceNumber = 32
)

// CCMap Concurrent Map
type CCMap struct {
	sliceNumber int
	values      []*MapSlice
}

// MapSlice Map slice
type MapSlice struct {
	lock  sync.RWMutex
	items map[string]interface{}
}

// NewDefaultCCMap New default slice number ccMap
func NewDefaultCCMap() *CCMap {
	return NewCCMap(defaultSliceNumber)
}

// NewCCMap New slice number ccMap
func NewCCMap(sliceNumber int) *CCMap {
	ccMap := &CCMap{
		sliceNumber: sliceNumber,
		values:      make([]*MapSlice, sliceNumber),
	}
	for i := 0; i < sliceNumber; i++ {
		ccMap.values[i] = &MapSlice{
			items: make(map[string]interface{}),
		}
	}
	return ccMap
}

func fnvHash(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// GetSliceMap get slice key
func (c *CCMap) GetSliceMap(key string) *MapSlice {
	return c.values[uint(fnvHash(key))%uint(c.sliceNumber)]
}

// IsExist key is exist
func (c CCMap) IsExist(key string) bool {
	sliceMap := c.GetSliceMap(key)

	sliceMap.lock.RLock()
	_, ok := sliceMap.items[key]
	sliceMap.lock.RUnlock()

	return ok
}

// Set set key value
func (c *CCMap) Set(key string, value interface{}) {
	sliceMap := c.GetSliceMap(key)

	sliceMap.lock.Lock()
	sliceMap.items[key] = value
	sliceMap.lock.Unlock()
}

// Get get by key
func (c *CCMap) Get(key string) interface{} {
	sliceMap := c.GetSliceMap(key)

	sliceMap.lock.RLock()

	value, ok := sliceMap.items[key]
	if ok {
		sliceMap.lock.RUnlock()
		return value
	}
	sliceMap.lock.RUnlock()
	return nil
}

// Delete delete by key
func (c *CCMap) Delete(key string) {
	sliceMap := c.GetSliceMap(key)

	sliceMap.lock.Lock()
	delete(sliceMap.items, key)
	sliceMap.lock.Unlock()
}

// Update update by key
// if key exist, update value
func (c *CCMap) Update(key string, value interface{}) {
	sliceMap := c.GetSliceMap(key)

	sliceMap.lock.RLock()
	value, ok := sliceMap.items[key]
	if ok {
		sliceMap.lock.RUnlock()
		sliceMap.lock.Lock()
		sliceMap.items[key] = value
		sliceMap.lock.Unlock()
		return
	}

	sliceMap.lock.RUnlock()
}

// Replace replace
// if key exist, update value.
// if key not exist, insert value.
func (c *CCMap) Replace(key string, value interface{}) {
	sliceMap := c.GetSliceMap(key)

	sliceMap.lock.Lock()
	sliceMap.items[key] = value
	sliceMap.lock.Unlock()
}

// MSet multiple set
func (c *CCMap) MSet(data map[string]interface{}) {
	for key, value := range data {
		c.Set(key, value)
	}
}

// MGet multiple get by keys
func (c *CCMap) MGet(keys ...string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, key := range keys {
		val := c.Get(key)
		data[key] = val
	}
	return data
}

// GetOnce get value by key and delete key
func (c *CCMap) GetOnce(key string) interface{} {
	val := c.Get(key)
	c.Delete(key)
	return val
}

// GetAll get all values
func (c *CCMap) GetAll() map[string]interface{} {
	data := make(map[string]interface{})

	for i := 0; i < c.sliceNumber; i++ {
		sliceMap := c.values[i]
		sliceMap.lock.RLock()
		for key, value := range sliceMap.items {
			data[key] = value
		}
		sliceMap.lock.RUnlock()
	}
	return data
}

// Clear clear all values
func (c *CCMap) Clear() {
	for i := 0; i < c.sliceNumber; i++ {
		sliceMap := c.values[i]
		sliceMap.lock.Lock()
		c.values[i].items = make(map[string]interface{})
		sliceMap.lock.Unlock()
	}
}

// Keys get all keys
func (c *CCMap) Keys() []string {
	data := []string{}
	for i := 0; i < c.sliceNumber; i++ {
		sliceMap := c.values[i]
		sliceMap.lock.RLock()
		for key := range sliceMap.items {
			data = append(data, key)
		}
		sliceMap.lock.RUnlock()
	}
	return data
}

// Count values count
func (c *CCMap) Count() int {
	count := 0
	for i := 0; i < c.sliceNumber; i++ {
		sliceMap := c.values[i]
		sliceMap.lock.RLock()
		count += len(sliceMap.items)
		sliceMap.lock.RUnlock()
	}
	return count
}
