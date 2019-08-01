package cache

import (
	"fmt"
	"sync"
	"time"
)



type strvalue struct {
	key string
	mu sync.RWMutex
	value interface{}
	start time.Time  // 创建的时间
}



var gocache *cache

type cache struct {
	str map[string]*strvalue
	defaultExpiration time.Duration
}

func Init() {
	gocache = &cache{
		str: make(map[string]*strvalue, 0),
	}
}

func Get(key string) interface{} {
	fmt.Println(time.Since(gocache.str[key].start))
	if time.Since(gocache.str[key].start) <= 0  {
		return gocache.str[key].value
		fmt.Println(time.Since(gocache.str[key].start))
	}
	return nil
}

func Set(key string, value interface{}, d time.Duration) {
	ss := &strvalue{
		key: key,
		value: value,
		mu: sync.RWMutex{},
		start: time.Now().Add(d),
	}
	gocache.str[key] = ss
}

