# cache
go 缓存
# 注意点
必须要有服务才能跑起来
# 例子 来自doc/main.go
```
package main

import (
	"cache"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cache.Init()
	cache.Set("aaa", "bbb", time.Second * 3)
	cache.Set("ccc", "bbb", time.Second * 5)
	fmt.Println("key:", cache.Get("aaa"))
	//var t *time.Timer
	time.Sleep(2 * time.Second)
	fmt.Println("key: aaa, value: ", cache.Get("aaa"))
	time.Sleep(2 * time.Second)
	fmt.Println("key: aaa, value: ", cache.Get("aaa"))
	fmt.Println("key: ccc, value: ", cache.Get("ccc"))

	if err := http.ListenAndServe(":7070", nil);err != nil {
		log.Fatal(err)
	}

}

```
输出
```
key: bbb
key: aaa, value:  bbb
key: aaa, value:  <nil>
key: ccc, value:  bbb
```
