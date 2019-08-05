# cache
go 缓存
# 注意点
必须要有服务才能跑起来
# 优点
实现过程很简单, 再set的时候 通过goroutine来删除key, 虽然删除key 方便,  但是goroutine 也会浪费丁点的资源
# demo 
from doc/main.go
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
	cache.Set("aaa", "bbb", time.Second * 3)  // 设置key和value ,  过期时间,  过期时间 <= 0 就是永不过期
	cache.Set("ccc", "bbb", time.Second * 5)
	fmt.Println("key:", cache.Get("aaa"))   // 获取值, 存在就返回值, 否则返回nil
	//var t *time.Timer
	time.Sleep(2 * time.Second)
	fmt.Println("key: aaa, value: ", cache.Get("aaa"))
	time.Sleep(2 * time.Second)
	fmt.Println(cache.TTL("ccc"))    // 获取key 过期时间, 不存在key 返回-1, 存在过期时间就返回, 否则返回0
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
0.999666731
key: aaa, value:  <nil>
key: ccc, value:  bbb
```
