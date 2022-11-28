package main

import (
	"cacheutil/lrucache"
	"time"
)

type lruCacheServer struct {
	cache lrucache.ICacher
}

// NewCache 实例化
func NewLruCache() *lruCacheServer {
	return &lruCacheServer{
		cache: lrucache.NewLruCache(
			lrucache.SetMaxMemory(200*lrucache.MBSize),
			lrucache.SetGCTime(500*time.Millisecond),
			lrucache.SetGCWorker(5),
		),
	}
}

// Set 写入缓存
func (c *lruCacheServer) Set(key string, value any, expire ...time.Duration) bool {
	//过期时间
	var tmpEx time.Duration
	if len(expire) > 0 {
		tmpEx = expire[0]
	}
	return c.cache.Set(key, value, tmpEx)
}

// ExistOrStore 若key不存在则写入
func (c *lruCacheServer) ExistOrStore(key string, value any, expire ...time.Duration) bool {
	//过期时间
	var tmpEx time.Duration
	if len(expire) > 0 {
		tmpEx = expire[0]
	}
	return c.cache.ExistOrStore(key, value, tmpEx)
}

// Get 获取
func (c *lruCacheServer) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

// Del 删除
func (c *lruCacheServer) Del(key string) bool {
	return c.cache.Del(key)
}

// Exist 判断是否存在
func (c *lruCacheServer) Exist(key string) bool {
	return c.cache.Exist(key)
}

// Flush 清空所有key
func (c *lruCacheServer) Flush() {
	c.cache.Flush()
}

// Count 获取key的数量
func (c *lruCacheServer) Count() int64 {
	return c.cache.Count()
}

// Keys 获取所有key
func (c *lruCacheServer) Keys() []string {
	return c.cache.Keys()
}

// GetUseMemory 获取当前已使用的内存
func (c *lruCacheServer) GetUseMemory() int64 {
	return int64(c.cache.GetUseMemory())
}
