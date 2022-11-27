package cache

import (
	"sync"
	"time"
)

type Cache struct {
	currentMemorySize MemorySize //当前使用内存大小
	maxMemorySize     MemorySize //最大内存

	values map[string]*cacheValue //数据
	vLock  sync.RWMutex           //数据读写锁

	garbageCollectionTime   time.Duration //回收间隔时间
	garbageCollectionWorker int           //回收消费者数
}

// NewCache 实例化
func NewCache(opts ...option) *Cache {
	c := &Cache{
		maxMemorySize:           200 * MBSize,                    //默认200MB
		values:                  make(map[string]*cacheValue, 0), //数据初始化
		garbageCollectionTime:   500 * time.Millisecond,          //默认500ms回收一次
		garbageCollectionWorker: 5,                               //默认每次启动5个消费者
	}

	//设置参数
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(c)
		}
	}

	go c.clearExpiredValues() //回收
	return c
}

// Set 写入缓存
func (c *Cache) Set(key string, value any, expire time.Duration) bool {
	c.vLock.Lock()
	defer c.vLock.Unlock()

	kSize := getValueSize(key)   //获取当前键大小
	vSize := getValueSize(value) //获取当前值大小

	var oldSize MemorySize
	oldValue, ok := c.get(key) //判断是否存在旧值，存在则获取旧值的size
	if ok {
		oldSize = oldValue.size
	}

	//判断是否超出最大内存限制，超出则set失败 （当前占用内存 - 旧值的size + 新值的size）
	if c.currentMemorySize-oldSize+kSize+vSize > c.maxMemorySize {
		return false
	}

	data := &cacheValue{
		value:     value,
		createdAt: time.Now(),
		expire:    expire,
		size:      kSize + vSize, //键名+键值
	}

	//删除旧值，增加新值
	c.del(key)
	c.add(key, data)

	return true
}

// ExistOrStore 若key不存在则写入
func (c *Cache) ExistOrStore(key string, value any, expire time.Duration) bool {
	c.vLock.Lock()
	defer c.vLock.Unlock()

	_, ok := c.get(key) //判断是否存在旧值
	if ok {
		return false
	}

	kSize := getValueSize(key)   //获取当前键大小
	vSize := getValueSize(value) //获取当前值大小

	//判断是否超出最大内存限制，超出则set失败 （当前占用内存 + 新值的size）
	if c.currentMemorySize+kSize+vSize > c.maxMemorySize {
		return false
	}

	data := &cacheValue{
		value:     value,
		createdAt: time.Now(),
		expire:    expire,
		size:      kSize + vSize, //键名+键值
	}

	//删除旧值，增加新值
	c.add(key, data)

	return true
}

// Get 获取
func (c *Cache) Get(key string) (any, bool) {
	c.vLock.RLock()
	defer c.vLock.RUnlock()

	data, ok := c.get(key)
	if ok {
		//永不超时 || 未超时，成功返回
		if data.expire == 0 || time.Now().Before(data.createdAt.Add(data.expire)) {
			return data.value, ok
		}
		c.del(key) //超时删除
	}
	return nil, false
}

// Del 删除
func (c *Cache) Del(key string) bool {
	c.vLock.Lock()
	defer c.vLock.Unlock()

	return c.del(key)
}

// Exist 判断是否存在
func (c *Cache) Exist(key string) bool {
	c.vLock.RLock()
	defer c.vLock.RUnlock()

	_, ok := c.values[key]
	return ok
}

// Flush 清空所有key
func (c *Cache) Flush() {
	c.vLock.Lock()
	defer c.vLock.Unlock()

	c.values = make(map[string]*cacheValue, 0)
	c.currentMemorySize = 0
}

// Count 获取key的数量
func (c *Cache) Count() int64 {
	c.vLock.RLock()
	defer c.vLock.RUnlock()

	return int64(len(c.values))
}

// Keys 获取所有key
func (c *Cache) Keys() []string {
	c.vLock.RLock()
	defer c.vLock.RUnlock()

	var keys []string
	for s := range c.values {
		keys = append(keys, s)
	}
	return keys
}

// GetUseMemory 获取当前已使用的内存
func (c *Cache) GetUseMemory() MemorySize {
	return c.currentMemorySize
}

func (c *Cache) get(key string) (*cacheValue, bool) {
	value, ok := c.values[key]
	return value, ok
}

func (c *Cache) add(key string, value *cacheValue) {
	c.values[key] = value
	c.currentMemorySize += value.size
}

func (c *Cache) del(key string) bool {
	tmp, ok := c.get(key)
	if ok {
		delete(c.values, key)           //从map删除
		c.currentMemorySize -= tmp.size //减去累计size
	}
	return ok //获取成功则确认删除成功，不存在则为失败
}

func (c *Cache) clearExpiredValues() {
	//定时器
	ticker := time.NewTicker(c.garbageCollectionTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//消费者
			for i := 0; i < c.garbageCollectionWorker; i++ {
				go c.gcDo()
			}
		}
	}
}

func (c *Cache) gcDo() {
	//遍历全部值
	for key, value := range c.values {
		c.vLock.Lock()
		//淘汰过期值
		if value.expire != 0 && time.Now().After(value.createdAt.Add(value.expire)) {
			c.del(key)
		}
		c.vLock.Unlock()
	}
}

type option func(*Cache)

// SetMaxMemory 设置最大内存 opt模式
func SetMaxMemory(size MemorySize) option {
	return func(c *Cache) {
		c.maxMemorySize = size
	}
}

// SetGCTime 设置回收间隔时间
func SetGCTime(t time.Duration) option {
	return func(c *Cache) {
		c.garbageCollectionTime = t
	}
}

// SetGCWorker 设置回收消费者数
func SetGCWorker(i int) option {
	return func(c *Cache) {
		c.garbageCollectionWorker = i
	}
}
