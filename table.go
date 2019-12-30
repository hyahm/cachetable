package cachetable

import (
	"errors"
	"reflect"
	"sync"
)

type Cache struct {
	keys  map[string]int                         // 保存key, 为了去重， 使用map
	cache map[string]map[interface{}]interface{} // 保存field
	s     interface{}                            // 保存表结构
	mu    sync.RWMutex
}

func NewTable(table interface{}) *Cache {
	return &Cache{
		keys:  make(map[string]int),
		cache: make(map[string]map[interface{}]interface{}),
		mu:    sync.RWMutex{},
		s:     table,
	}

}

func (c *Cache) Add(table interface{}) error {
	//必须是指针
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return errors.New("must be a pointer")
	}
	if len(c.keys) == 0 {
		return errors.New("at least set one key")
	}
	var sv reflect.Type
	if reflect.TypeOf(c.s).Kind() == reflect.Ptr {
		sv = reflect.TypeOf(c.s).Elem()
	} else {
		sv = reflect.TypeOf(c.s)
	}
	// 必须是同一类型
	if sv == reflect.TypeOf(table).Elem() {
		c.mu.Lock()
		defer c.mu.Unlock()
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

func (c *Cache) SetKey(key string) error {
	if c == nil {
		return errors.New("init first")
	}
	var sv reflect.Type
	if reflect.TypeOf(c.s).Kind() == reflect.Ptr {
		sv = reflect.TypeOf(c.s).Elem()
	} else {
		sv = reflect.TypeOf(c.s)
	}
	// 判断key 是否有效
	if _, ok := sv.FieldByName(key); !ok {
		return errors.New("not a same struct")
	}
	c.keys[key] = 0

	return nil
}

func (c *Cache) Set(setKey string, setValue interface{}, searchKey string, searchValue interface{}) error {
	if c == nil {
		return errors.New("init first")
	}
	// tag
	if _, ok := c.keys[searchKey]; !ok {
		return errors.New("key must be key")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if f, ok := c.cache[searchKey]; ok {
		if v, ok := f[searchValue]; ok {
			reflect.ValueOf(v).Elem().FieldByName(setKey).Set(reflect.ValueOf(setValue))
			if value, ok := c.keys[setKey]; ok {
				if value == setValue {
					return errors.New("Duplicate key value")
				}
				// 如果是主键， 更新map
				f[setValue] = v
				delete(f, value)
			}
		}
	}
	return nil
}

func (c *Cache) Get(key string, field string, value interface{}) (interface{}, error) {
	if c == nil {
		return nil, errors.New("init first")
	}
	// 如果是索引， 直接返回即可
	if f, ok := c.cache[field]; ok {
		if v, ok := f[value]; ok {
			return reflect.ValueOf(v).Elem().FieldByName(key).Interface(), nil
		}
	} else {
		//遍历

		return nil, errors.New("field not a key")
	}
	//

	return nil, nil
}

func (c *Cache) del(key string, field string, value interface{}) {
	// 如果是索引， 直接返回即可
	v, _ := c.Get(key, field, value)
	if _, ok := c.cache[key]; ok {
		delete(c.cache[key], v)
	}

}

func (c *Cache) Del(field string, value interface{}) error {
	if c == nil {
		return errors.New("init first")
	}
	// 如果是索引， 直接返回即可
	if _, ok := c.cache[field]; ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, _ := range c.keys {
			// 删掉其他的key
			if k != field {
				c.del(k, field, value)
			}

		}
		delete(c.cache[field], value)
		//if v, ok := f[value]; ok {
		//	return  nil
		//}
	}
	//

	return nil
}

func (c *Cache) GetKeys() (ks []string) {
	for k, _ := range c.keys {
		ks = append(ks, k)
	}
	return
}
