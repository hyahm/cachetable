package cachetable

import (
	"reflect"
	"sync"
	"time"
)

type Filter struct {
	Row   *row
	c      *Cache
	Err    error
	mu     sync.RWMutex
}

func (c *Cache) Filter(field string, value interface{}) *Filter {
	if c.s == nil {
		return &Filter{
			Row: nil,
			Err: ErrorNotInit,
		}

	}
	if len(c.keys) == 0 {
		return &Filter{
			Row: nil,
			Err: ErrorNoKey,
		}

	}
	// 找到所有索引， 删除,   必须是key
	if vms, ok := c.cache[field]; ok {
		//找到所有所有的keys 的值
		key, _ := c.toString(value)
		if vms[key].CanExpire && time.Now().Sub(vms[key].Expire) >= 0 {
			// 说明过期了
			//直接先删掉

			f := &Filter{
				Row: vms[key],
				Err: ErrorExpired,
				c:   c,
				mu:  vms[key].mu,
			}
			f.mu.Lock()
			f.Del()
			f.mu.Unlock()
			f.Row = nil
			return f
		}
		return &Filter{
			Row:    vms[key],
			Err:    nil,
			c:      c,
			mu:     vms[key].mu,
		}
		return nil
	} else {
		return &Filter{
			Row: nil,
			Err: ErrorNoFeildKey,
		}
	}

}

//type Filter interface {
//	Get(keys ...string) []interface{}
//	Set(field string, Value interface{}) error
//	Del() error
//}

func (f *Filter) expired() bool {
	return f.TTL() == 0
}

func (f *Filter) Expired() bool {
	return f.expired()
}

func (f *Filter) TTL() int64 {
	if f.Row == nil {
		return 0
	}
	if f.Row.CanExpire {
		if time.Now().Sub(f.Row.Expire) >= 0 {
			return 0
		} else {
			return int64(f.Row.Expire.Sub(time.Now()).Seconds())
		}

	} else {
		return -1
	}
}

func (f *Filter) Get(keys ...string) []interface{} {
	if f.Row == nil {
		return nil
	}
	l := len(keys)
	vs := make([]interface{}, l)
	if f.expired() {
		return vs
	}

	for i, v := range keys {
		val := reflect.ValueOf(f.Row.value).Elem().FieldByName(v).Interface()
		vs[i] = val
	}
	return vs
}

func (f *Filter) Del() error {
	if f.Row == nil {
		return f.Err
	}
	if f.expired() {
		return ErrorExpired
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	//找到所有所有的keys 的值
	for _, k := range f.c.keys {
		//v := c.Get(k)
		ft := reflect.ValueOf(f.Row).FieldByName(k).Interface()
		value, _ := f.c.toString(ft)
		delete(f.c.cache[k], value)
	}

	return f.Err

}

func (f *Filter) SetTTL(t time.Duration) error {
	if f.Row == nil {
		return f.Err
	}

	if f.expired() {
		return ErrorExpired
	}
	if t <= 0 {
		f.Row.CanExpire = false
		return nil
	}
	if t > 0 {
		f.Row.CanExpire = true
		f.Row.Expire = time.Now().Add(t)
		return nil
	}

	return f.Err
}


func (f *Filter) Set(field string, value interface{}) error {
	if f.Row == nil {
		return f.Err
	}

	if f.expired() {
		return ErrorExpired
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if ok := f.c.hasKey(field); ok {
		// 如果是key, value不能重复
		newvalue_str, _ := f.c.toString(value)
		if _, ok := f.c.cache[field][newvalue_str]; ok {
			return ErrorDuplicate
		}

		// 如果是设置的是key是主键， 重新生成
		oldvalue_str, _ := f.c.toString(value)

		f.c.cache[field][newvalue_str] = &row{
			mu:    sync.RWMutex{},
			value: f.Row,
		}
		// 删掉老的键值
		delete(f.c.cache[field], oldvalue_str)
	}

	// 更新v
	newv := reflect.ValueOf(value)
	reflect.ValueOf(f.Row.value).Elem().FieldByName(field).Set(newv)

	return f.Err
}
