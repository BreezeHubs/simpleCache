# 内存缓存工具

## 涉及知识
+ 集合
+ iota
+ 位运算
+ 读写锁
+ 协程
+ 定时器
+ select case

## 系统需求
+ 支持设置过期时间，精度到秒
+ 支持设置最大内存，当内存溢出时做出合适的处理
+ 并发安全

## 接口
```go
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
```