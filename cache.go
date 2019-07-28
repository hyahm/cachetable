package cache

import (
	"fmt"
	"sync"
	"time"
)

type auto uint64

func (a auto) Init() {

	go func() {
		a.reduce()
	}()
	select {
	case a == 0:
		
	}
}

func (a auto) reduce() {
	t := time.Tick(time.Second)
	for _ = range t {
		a = a + 1
	}
}

type strvalue struct {
	key string
	mu sync.RWMutex
	value interface{}
	expiration time.Duration // 过期时间， 单位： 秒
	start time.Time  // 创建的时间
}

type element struct {
	key string
	ty string  // # key 的类型， 暂时只有str
}

var el []*element // 保存key, 按照过期时间排序

var gocache *cache

type cache struct {
	str map[string]*strvalue
	defaultExpiration time.Duration
}

func Init() {
	var x auto
	x.add()
	el = make([]*element, 0)
	gocache = &cache{
		str: make(map[string]*strvalue, 0),
	}
}

func Get(key string) interface{} {
	return gocache.str[key]
}

func Set(key string, value interface{}, d time.Duration) {
	ss := &strvalue{
		key: key,
		value: value,
		mu: sync.RWMutex{},
		start: time.Now(),
	}
	if d > 0 {
		ss.expiration = d
		gocache.str[key] = ss
	}
	gocache.str[key] = ss
}

// 按快到期的秒数来排序
func sort(key string, d uint64) {
	for _,e := range el {
		switch  e.ty {
		case "str":
			//ep := time.Since(gocache.str[e.key].start.Add(gocache.str[e.key].expiration)).Seconds()
			print(gocache.str[e.key].expiration)
		}
	}
}
