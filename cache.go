package main

import (
	"cacheutil/cache"
	"time"
)

type cacheServer struct {
	cache cache.ICacher
}

// NewCache 实例化
func NewCache() *cacheServer {
	return &cacheServer{
		cache: cache.NewCache(
			cache.SetMaxMemory(200*cache.MBSize),
			cache.SetGCTime(500*time.Millisecond),
			cache.SetGCWorker(5),
		),
	}
}

// Set 写入缓存
func (c *cacheServer) Set(key string, value any, expire ...time.Duration) bool {
	//过期时间
	var tmpEx time.Duration
	if len(expire) > 0 {
		tmpEx = expire[0]
	}
	return c.cache.Set(key, value, tmpEx)
}

// Get 获取
func (c *cacheServer) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

// Del 删除
func (c *cacheServer) Del(key string) bool {
	return c.cache.Del(key)
}

// Exist 判断是否存在
func (c *cacheServer) Exist(key string) bool {
	return c.cache.Exist(key)
}

// Flush 清空所有key
func (c *cacheServer) Flush() {
	c.cache.Flush()
}

// Count 获取key的数量
func (c *cacheServer) Count() int64 {
	return c.cache.Count()
}

// Keys 获取所有key
func (c *cacheServer) Keys() []string {
	return c.cache.Keys()
}

// GetUseMemory 获取当前已使用的内存
func (c *cacheServer) GetUseMemory() int64 {
	return int64(c.cache.GetUseMemory())
}
