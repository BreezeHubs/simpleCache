package cache

import (
	"encoding/json"
)

func getValueSize(value any) MemorySize {
	bytes, _ := json.Marshal(value)
	//fmt.Println(string(bytes), MemorySize(len(bytes)))
	return MemorySize(len(bytes))
}
