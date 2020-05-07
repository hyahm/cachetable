package cachetable

import (
	"reflect"
	"time"
)

type Cache struct {
	Keys  []string                   // 保存key, 为了去重， 使用map
	Cache map[string]map[string]*Row // 保存field， 将所有值都转为string
	S     interface{}                // 保存表结构

}

func (c *Cache) Add(table interface{}, expire time.Duration) error {
	//必须是指针
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return ErrorNotPointer
	}
	if len(c.Keys) == 0 {
		return ErrorNoKey
	}
	var st reflect.Type
	if reflect.TypeOf(c.S).Kind() == reflect.Ptr {
		st = reflect.TypeOf(c.S).Elem()
	} else {
		st = reflect.TypeOf(c.S)
	}
	// 必须是同一类型
	if st == reflect.TypeOf(table).Elem() {
		//遍历 添加key
		for _, k := range c.Keys {
			// 将字段的值全部转化为string

			if _, ok := c.Cache[k]; !ok {
				// 没有字段， 初始化
				c.Cache[k] = make(map[string]*Row)

			}

			kv := asString(reflect.ValueOf(table).Elem().FieldByName(k).Interface())

			r := &Row{
				Value: table,
			}
			if expire > 0 {
				r.Expire = time.Now().Add(expire)
				r.CanExpire = true
			}

			cmu.Lock()
			c.Cache[k][kv] = r
			cmu.Unlock()
		}

	} else {
		return ErrorStruct
	}
	return nil
}

func (c *Cache) SetKeys(keys ...string) error {
	if c.S == nil {
		return ErrorNotInit
	}
	var sv reflect.Type
	if reflect.TypeOf(c.S).Kind() == reflect.Ptr {
		sv = reflect.TypeOf(c.S).Elem()
	} else {
		sv = reflect.TypeOf(c.S)
	}
	// 判断key 是否有效
	for _, k := range keys {
		if _, ok := sv.FieldByName(k); ok {
			if !c.hasKey(k) {
				c.Keys = append(c.Keys, k)
			}

		}
	}

	return nil
}

//

//
func (c *Cache) GetKeys() (ks []string) {
	return c.Keys
}

func (c *Cache) hasKey(s string) bool {
	for _, v := range c.Keys {
		if v == s {
			return true
		}
	}
	return false
}

func (c *Cache) clean(t time.Duration) {
	// 清除过期table
	if len(c.Keys) == 0 {
		panic(ErrorNoKey)
	}
	for {
		// 第一个字段就行了
		time.Sleep(t)
		allmap := c.Cache[c.Keys[0]]
		for k, v := range allmap {
			if !v.CanExpire && time.Now().Sub(v.Expire) >= 0 {
				c.Filter(c.Keys[0], k).Del()
			}
		}
	}
}

// 通过key 获取结构
func (c *Cache) GetAllLine() []interface{} {
	if len(c.Keys) == 0 {
		panic(ErrorNoKey)
	}
	l := len(c.Cache[c.Keys[0]])
	lines := make([]interface{}, 0, l)
	for _, v := range c.Cache[c.Keys[0]] {
		// 判断是否过期

		if v.CanExpire && v.Expire.Sub(time.Now()) <= 0 {
			continue
		}
		lines = append(lines, v.Value)
	}
	return lines
}
