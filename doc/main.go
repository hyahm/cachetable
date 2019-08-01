package main

import (
	"cache"
	"fmt"
	"time"
)

func main() {
	cache.Init()
	cache.Set("aaa", "bbb", time.Second * 10)
	fmt.Println(cache.Get("aaa"))

}
