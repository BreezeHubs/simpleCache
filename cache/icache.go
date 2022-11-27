package cache

import "time"

type ICacher interface {
	// Set 写入缓存
	Set(key string, value any, expire time.Duration) bool

	// Get 获取
	Get(key string) (any, bool)

	// Del 删除
	Del(key string) bool

	// Exist 判断是否存在
	Exist(key string) bool

	// Flush 清空所有key
	Flush()

	// Count 获取key的数量
	Count() int64

	Keys() []string

	GetUseMemory() MemorySize
}

type MemorySize int64

const (
	KBSize MemorySize = 1 << (10 * iota) //1^0
	MBSize                               //1^10
	GBSize                               //1^20
	TBSize                               //1^30
	PBSize                               //1^40
)

type cacheValue struct {
	value     any           //数据 键名+键值
	size      MemorySize    //数据大小
	createdAt time.Time     //创建时间
	expire    time.Duration //有效时长
}
