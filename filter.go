package cachetable

import (
	"reflect"
	"sync"
)

type F struct {
	Row interface{}
	c   *Cache
	Err error
	mu  sync.RWMutex
}

func (c *Cache) Filter(field string, value interface{}) *F {
	if c.s == nil {
		return &F{
			Err: ErrorNotInit,
		}

	}
	if len(c.keys) == 0 {
		return &F{
			Err: ErrorNoKey,
		}

	}
	// 找到所有索引， 删除,   必须是key
	if vms, ok := c.cache[field]; ok {
		//找到所有所有的keys 的值
		key, _ := c.toString(reflect.ValueOf(value))

		return &F{
			Row: vms[key].value,
			Err: nil,
			c:   c,
			mu:  vms[key].mu,
		}
		return nil
	} else {
		return &F{
			Err: ErrorNoFeildKey,
		}
	}

}

type Filter interface {
	Get(keys ...string) []interface{}
	Set(field string, Value interface{}) error
	Del() error
}

func (c *F) Get(keys ...string) []interface{} {
	vs := make([]interface{}, 0)
	for _, v := range keys {
		i := reflect.ValueOf(c.Row).Elem().FieldByName(v).Interface()
		vs = append(vs, i)
	}
	return vs
}

func (c *F) Del() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	//找到所有所有的keys 的值
	for k, _ := range c.c.keys {
		v := c.Get(k)
		value, _ := c.c.toString(reflect.ValueOf(v))
		delete(c.c.cache[k], value)
	}

	return c.Err

}

func (c *F) Set(field string, Value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.c.keys[field]; ok {
		// 如果是key
		// 如果是设置的是key是主键， 重新生成
		oldvalue := c.Get(field)
		oldvalue_str, _ := c.c.toString(reflect.ValueOf(oldvalue))
		newvalue_str, _ := c.c.toString(reflect.ValueOf(Value))
		c.c.cache[field][newvalue_str] = &row{
			mu:    sync.RWMutex{},
			value: c.Row,
		}
		// 删掉老的键值
		delete(c.c.cache[field], oldvalue_str)
	}

	// 更新v
	newv := reflect.ValueOf(Value)
	reflect.ValueOf(c.Row).Elem().FieldByName(field).Set(newv)

	return c.Err
}
