package main

import (
	"bytes"
	"cacheutil/cache"
	"fmt"
	"time"
)

func main() {
	c := NewCache()

	var s bytes.Buffer
	for i := 0; i < 1000000; i++ {
		s.WriteString("aaaaaaaaaaaaaaaaaaaa")
	}
	c.Set("a", s.String(), 1*time.Nanosecond)

	c.Set("b", 1231414523)
	ok := c.ExistOrStore("b", 1231414523)
	fmt.Println("ExistOrStore: ", ok)

	fmt.Println(c.Count())
	fmt.Println(c.Keys())
	fmt.Println(cache.MemorySize(c.GetUseMemory())/cache.GBSize, "GB")

	time.Sleep(1 * time.Second)

	fmt.Println(c.Count())
	fmt.Println(c.Keys())
	fmt.Println(c.GetUseMemory(), "KB")
}
