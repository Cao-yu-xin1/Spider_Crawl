// Package cache 本地缓存实现
// ⚠️ 此文件包含故意设置的数据竞态问题，用于 bugfix 练习
package cache

import (
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	ExpireTime time.Time
}

// LocalCache 本地缓存
type LocalCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

// NewLocalCache 创建本地缓存
func NewLocalCache() *LocalCache {
	return &LocalCache{
		items: make(map[string]*CacheItem),
	}
}

// Set 设置缓存
func (c *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	// 问题: 没有加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	item := &CacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(ttl),
	}
	c.items[key] = item
}

// Get 获取缓存
// 问题: 并发读写 map 会导致 race condition
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpireTime) {
		delete(c.items, key)
		return nil, false
	}

	return item.Value, true
}

// Delete 删除缓存
// 问题: 并发删除 map 会导致 panic
func (c *LocalCache) Delete(key string) {
	delete(c.items, key)
}

// Clear 清空缓存
func (c *LocalCache) Clear() {
	c.items = make(map[string]*CacheItem)
}

// Size 获取缓存大小
// 问题: 并发读取 map 长度会有 race condition
func (c *LocalCache) Size() int {
	// 问题: 没有加锁
	return len(c.items)
}

// Keys 获取所有键
func (c *LocalCache) Keys() []string {
	// 问题: 没有加锁
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有值
// 问题: 遍历 map 时并发修改会导致 panic
func (c *LocalCache) Values() []interface{} {
	values := make([]interface{}, 0, len(c.items))
	for _, v := range c.items {
		values = append(values, v.Value)
	}
	return values
}

// Has 检查键是否存在
func (c *LocalCache) Has(key string) bool {
	_, exists := c.items[key]
	return exists
}

// GetOrSet 获取或设置缓存
func (c *LocalCache) GetOrSet(key string, value interface{}, ttl time.Duration) interface{} {
	// 问题: 没有加锁
	if v, ok := c.Get(key); ok {
		return v
	}

	c.Set(key, value, ttl)
	return value
}

// DeleteExpired 删除过期缓存
// 问题: 遍历删除不是原子操作
func (c *LocalCache) DeleteExpired() int {
	// 问题: 没有加锁
	count := 0
	now := time.Now()

	// ⚠️ Bug: 遍历 map 时删除元素会导致 panic
	for k, v := range c.items {
		if now.After(v.ExpireTime) {
			delete(c.items, k)
			count++
		}
	}

	return count
}

// Range 遍历缓存
// 问题: 遍历过程中可能有并发修改
func (c *LocalCache) Range(fn func(key string, value interface{}) bool) {
	for k, v := range c.items {
		if !fn(k, v.Value) {
			break
		}
	}
}

// Increment 递增计数器
// 问题: 读-修改-写不是原子操作
func (c *LocalCache) Increment(key string) int {
	// 问题: 没有加锁
	var count int
	if v, ok := c.items[key]; ok {
		if num, ok := v.Value.(int); ok {
			count = num + 1
			v.Value = count
		}
	} else {
		count = 1
		c.items[key] = &CacheItem{
			Value:      count,
			ExpireTime: time.Now().Add(time.Hour),
		}
	}
	return count
}

// Decrement 递减计数器
// 问题: 读-修改-写不是原子操作
func (c *LocalCache) Decrement(key string) int {
	// 问题: 没有加锁
	var count int
	if v, ok := c.items[key]; ok {
		if num, ok := v.Value.(int); ok {
			count = num - 1
			v.Value = count
		}
	}
	return count
}

// GetMulti 批量获取
// 问题: 批量操作不是原子的
func (c *LocalCache) GetMulti(keys []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range keys {
		if v, ok := c.Get(key); ok {
			result[key] = v
		}
	}
	return result
}

// SetMulti 批量设置
// 问题: 批量操作不是原子的
func (c *LocalCache) SetMulti(items map[string]interface{}, ttl time.Duration) {
	// 问题: 没有加锁
	for k, v := range items {
		c.Set(k, v, ttl)
	}
}

// DeleteMulti 批量删除
func (c *LocalCache) DeleteMulti(keys []string) {
	for _, key := range keys {
		c.Delete(key)
	}
}
