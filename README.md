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
type Cacher interface {
	//size: 1KB 100KB 1MB  2MB 1GB
	SetMaxMemory(size string) bool

	//写入缓存
	Set(key string, value any, expire time.Duration) bool

	//获取
	Get(key string) (any, bool)

	//删除
	Del(key string) bool

	//判断是否存在
	Exist(key string) bool

	//清空所有key
	Flush() bool

	//获取key的数量
	Count() int64
}
```