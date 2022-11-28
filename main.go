package main

import (
	"bytes"
	"cacheutil/cache"
	"fmt"
	"time"
)

func main() {
	c := NewLruCache()

	var s bytes.Buffer
	for i := 0; i < 10200; i++ {
		s.WriteString("aaaaaaaaaaaaaaaaaaaa")
	}
	c.Set("a", s.String(), 1*time.Second)

	c.Set("b", 1231414523)
	ok := c.ExistOrStore("b", 1231414523)

	c.Set("c", "123141452312314145231231414523123141452312314145231231414523123141452312314145231231414523123141452312314145231231414523")
	fmt.Println("ExistOrStore: ", ok)

	//c.Get("a")
	c.Get("c")
	fmt.Println("count: ", c.Count())
	fmt.Println("keys: ", c.Keys())
	//fmt.Println(cache.MemorySize(c.GetUseMemory())/cache.GBSize, "GB")
	fmt.Println(cache.MemorySize(c.GetUseMemory())/cache.MBSize, "MB")

	time.Sleep(2 * time.Second)

	fmt.Println("count: ", c.Count())
	fmt.Println("keys: ", c.Keys())
	fmt.Println(c.GetUseMemory(), "KB")

	for {
	}
}
