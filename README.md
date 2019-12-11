# cache
go 缓存， 后期有时间会弄的

# 为什么要写这个
也是对redis的补充，我想通过一个key改变一个值， redis是可以做到的，而且很好， 
缓存的数据很多是关系型数据库的表，   那么多个key对应一个值呢， 比如一张表有多个唯一值， 我想通过人一个来改变值， 这个值是这些key共同的
redis 就无法实现了

多key对应一个值的思路
最难实现应该算是过期时间， 协议， 认证和存储
所以涉及到过期时间还是用redis，
应为多key对应唯一值目前考虑存储思路， 目前是增量存储


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
