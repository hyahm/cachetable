package cachetable

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type Filter struct {
	Row    interface{}
	c      *Cache
	Err    error
	mu     sync.RWMutex
	expire time.Time
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
		if vms[key].expire.Unix() != -62135596800 && time.Now().Unix() >= vms[key].expire.Unix() {
			// 说明过期了

			return &Filter{
				Row: nil,
				Err: errors.New("table expired"),
				c:   c,
				mu:  vms[key].mu,
			}
		}
		return &Filter{
			Row:    vms[key].value,
			Err:    nil,
			c:      c,
			expire: vms[key].expire,
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

func (f *Filter) TTL() time.Duration {
	if f.Row == nil {
		return 0
	}
	if f.expire.Unix() == -62135596800 {
		return -1
	}
	if time.Now().Unix() <= f.expire.Unix() {
		return 0
	}
	return 0
}

func (f *Filter) Get(keys ...string) []interface{} {
	if f.Row == nil {
		return nil
	}
	vs := make([]interface{}, 0)
	for _, v := range keys {
		i := reflect.ValueOf(f.Row).Elem().FieldByName(v).Interface()
		vs = append(vs, i)
	}
	return vs
}

func (f *Filter) Del() error {
	if f.Row == nil {
		return f.Err
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

func (f *Filter) Set(field string, value interface{}) error {
	if f.Row == nil {
		return f.Err
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
	reflect.ValueOf(f.Row).Elem().FieldByName(field).Set(newv)

	return f.Err
}
