package cachetable

import (
	"reflect"
	"time"
)

type Table struct {
	Keys  []string                   // 保存key, 为了去重， 使用map
	Cache map[string]map[string]*Row // 保存field， 将所有值都转为string, 双层map， 第一层是key， 第二层是key的值
	Typ   interface{}                // 保存表结构

}

func (c *Table) Add(table interface{}, expire time.Duration) error {
	//table必须是指针
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return ErrorNotPointer
	}
	if len(c.Keys) == 0 {
		return ErrorNoKey
	}

	st := reflect.TypeOf(table)
	ct := reflect.TypeOf(c.Typ)

	// 必须是同一类型
	if reflect.DeepEqual(ct, st) {
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

func (c *Table) SetKeys(keys ...string) error {
	if c.Typ == nil {
		return ErrorNotInit
	}

	sv := reflect.TypeOf(c.Typ)

	// 判断key 是否有效
	for _, k := range keys {
		if _, ok := sv.Elem().FieldByName(k); ok {
			if !c.hasKey(k) {
				c.Keys = append(c.Keys, k)
			}

		}
	}

	return nil
}

//

//
func (c *Table) GetKeys() (ks []string) {
	return c.Keys
}

func (c *Table) hasKey(s string) bool {
	for _, v := range c.Keys {
		if v == s {
			return true
		}
	}
	return false
}

func (c *Table) clean(t time.Duration) {
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
				f, err := c.Filter(c.Keys[0], k)
				if err != nil {
					continue
				}
				f.Del()
			}
		}
	}
}

// 通过key 获取结构
func (c *Table) GetAllLine() []interface{} {
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

// 通过key 获取结构
func (c *Table) Columns(col string) []interface{} {
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

		lines = append(lines, reflect.ValueOf(v.Value).Elem().FieldByName(col).Interface())
	}
	return lines
}
