package cache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetValueSize(t *testing.T) {
	c := 111
	getValueSize("aaaaaa")
	getValueSize(1231414523)
	getValueSize(map[string]any{
		"A": &c,
		"B": 222,
	})

	s := "bbbbbb"
	getValueSize(&s)
}

func TestGetDataCount(t *testing.T) {
	c := NewCache(SetMaxMemory(100 * GBSize))

	c1 := 111
	c.Set("a", "aaaaaa", 1*time.Hour)
	c.Set("b", 1231414523, 1*time.Hour)
	c.Set("aaa", map[string]any{
		"A": &c1,
		"B": 222,
	}, 1*time.Hour)
	s := "bbbbbb"
	c.Set("c", &s, 1*time.Hour)
	//for i := 0; i < 150000; i++ {
	//	c.Set(strconv.Itoa(i), i, 1*time.Hour)
	//}

	fmt.Println(c.Count())
	fmt.Println(c.Keys())
	fmt.Println(c.GetUseMemory(), "KB")
}

func TestClear(t *testing.T) {
	c := NewCache(SetMaxMemory(100 * GBSize))

	c.Set("a", "aaaaaa", 1*time.Second)
	c.Set("b", 1231414523, 1*time.Hour)

	fmt.Println(c.Count())
	fmt.Println(c.Keys())
	fmt.Println(c.GetUseMemory(), "KB")

	time.Sleep(1 * time.Second)

	fmt.Println(c.Count())
	fmt.Println(c.Keys())
	fmt.Println(c.GetUseMemory(), "KB")
}

func TestCache(t *testing.T) {
	testData := []struct {
		key    string
		value  any
		expire time.Duration
		want   any
	}{
		{"a", "aaaaaa", 10 * time.Second, "aaaaaa"},
		{"b", 2, 11 * time.Second, 2},
		{"c", false, 12 * time.Second, false},
		{"d", 4, 13 * time.Second, 4},
		{"e", 5, 14 * time.Second, 5},
		{"f", map[string]any{"A": "aaa", "B": 222}, 15 * time.Second, map[string]any{"A": "aaa", "B": 222}},
		{"g", 7, 16 * time.Second, 7},
	}

	c := NewCache()
	for _, v := range testData {
		c.Set(v.key, v.value, v.expire)

		val, ok := c.Get(v.key)
		//t.Log("val: ", val, v.value)
		if !ok {
			t.Errorf("%s 缓存失败\n", v.key)
		}

		if !reflect.DeepEqual(val, v.want) {
			t.Errorf("%s 缓存值异常\n", v.key)
		}
	}
}

func TestDelete(t *testing.T) {
	c := NewCache()

	c.Set("a", "aaaaaa", 1*time.Second)
	c.Set("b", 1231414523, 1*time.Hour)

	fmt.Println(c.Exist("b"))
	c.Del("b")
	fmt.Println(c.Exist("b"))
	c.Flush()

	fmt.Println(c.Count())
	fmt.Println(c.Keys())
}
