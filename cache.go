package cache

import (
	"sync"
	"time"
)

type strvalue struct {
	key string
	value interface{}
	end time.Time  // 过期的时间
}

var gocache *cache

type cache struct {
	mu sync.RWMutex
	str map[string]*strvalue
	defaultExpiration time.Duration
}

var stop chan string

func Init() {
	stop = make(chan string)
	gocache = &cache{
		str: make(map[string]*strvalue, 0),
		mu: sync.RWMutex{},
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
	gocache.mu.Lock()
	ss := &strvalue{
		key: key,
		value: value,
		end: time.Now().Add(d),
	}
	gocache.str[key] = ss
	gocache.mu.Unlock()
	if d > 0 {
		go expire(key, d)
	}
}

func Del(key string) {
	gocache.mu.Lock()
	if Exist(key) {
		delete(gocache.str,key)
	}
	gocache.mu.Unlock()

}

func TTL(key string) float64 {
	if _, ok := gocache.str[key]; !ok {
		return -1
	}
	exp := time.Since(gocache.str[key].end).Seconds()
	if exp < 0 {
		return  exp * -1
	}
	return 0
}

func Exist(key string) bool {
	if _, ok := gocache.str[key]; ok {
		return true
	}
	return false
}

func expire(key string, d time.Duration) {
	select {
	case <- time.After(d):
		gocache.mu.Lock()
		if Exist(key) {
			delete(gocache.str,key)
		}
		gocache.mu.Unlock()
	}
}
