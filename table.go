package cache

import (
	"errors"
	"reflect"
)

func NewCache() *Cache {
	return &Cache{
		keys:  make(map[string]int),
		cache: make(map[string]map[interface{}]interface{}),
		s:     make(map[string]interface{}),
	}

}

type Cache struct {
	keys  map[string]int                         // 保存key, 为了去重， 使用map
	cache map[string]map[interface{}]interface{} // 保存field
	s     interface{}                            // 保存表结构
}

func (c *Cache) Table(table interface{}) error {
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return errors.New("must be a pointer")
	}
	c.s = table
	return nil
}

func (c *Cache) Add(table interface{}) error {
	//必须是指针
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return errors.New("must be a pointer")
	}
	if len(c.keys) == 0 {
		return errors.New("at least set one key")
	}

	// 必须是同一类型
	if reflect.TypeOf(c.s).Elem() == reflect.TypeOf(table).Elem() {
		for v, _ := range c.keys {
			if _, ok := c.cache[v]; !ok {
				c.cache[v] = make(map[interface{}]interface{})
			}
			key := reflect.ValueOf(table).Elem().FieldByName(v).Interface() // 获取tag 的值
			if _, ok := c.cache[v][key]; ok {
				return errors.New("Duplicate key value")
			}
			c.cache[v][key] = table
		}

	} else {
		return errors.New("not a same struct")
	}
	return nil
}

func (c *Cache) Key(key string) error {
	if c == nil {
		return errors.New("init first")
	}
	// 判断key 是否有效
	if _, ok := reflect.TypeOf(c.s).Elem().FieldByName(key); !ok {
		return errors.New("not a same struct")
	}
	c.keys[key] = 0

	return nil
}

func (c *Cache) Set(key string, value interface{}, setkey string, setvalue interface{}) error {
	if c == nil {
		return errors.New("init first")
	}
	// tag
	if _, ok := c.keys[key];!ok {
		return errors.New("key must be key")
	}

	if f, ok := c.cache[key]; ok {
		if v, ok := f[value]; ok {
			reflect.ValueOf(v).Elem().FieldByName(setkey).Set(reflect.ValueOf(setvalue))
			if value, ok := c.keys[setkey];ok {
				if value == value {
					return errors.New("Duplicate key value")
				}
				// 如果是主键， 更新map
				f[setvalue] = v
				delete(f, value)
			}
		}
	}
	return nil
}

func (c *Cache) GetValue(key string, field string, value interface{}) (interface{}, error) {
	if c == nil {
		return nil, errors.New("init first")
	}
	// 如果是索引， 直接返回即可
	if f, ok := c.cache[field]; ok {
		if v, ok := f[value]; ok {
			return reflect.ValueOf(v).Elem().FieldByName(key).Interface(), nil
		}
	} else {
		return nil, errors.New("field not a key")
	}
	//

	return nil, nil
}
