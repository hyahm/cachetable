package cache

import (
	"sync"
	"time"
)

type strvalue struct {
	key string
	mu sync.RWMutex
	value interface{}
	end time.Time  // 过期的时间
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
	if _, ok := gocache.str[key]; !ok {
		return nil
	}
	if time.Since(gocache.str[key].end) < 0  {
		return gocache.str[key].value
	}
	return nil
}

func Set(key string, value interface{}, d time.Duration) {

	ss := &strvalue{
		key: key,
		value: value,
		mu: sync.RWMutex{},
		end: time.Now().Add(d),
	}
	gocache.str[key] = ss

	if d > 0 {
		go expire(key, d)

	}
}

func TTL(key string) time.Duration {
	return time.Since(gocache.str[key].end)
}

func expire(key string, d time.Duration) {
	select {
	case <- time.After(d):
		delete(gocache.str,key)
	}
}
