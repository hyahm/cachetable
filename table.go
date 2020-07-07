package cachetable

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Table struct {
	Keys  []string                        // 保存key, 为了去重， 使用map
	Cache map[string]map[interface{}]*Row // 保存field， 将所有值都转为string, 双层map， 第一层是key， 第二层是key的值
	Typ   interface{}                     // 保存表结构
	Name  string
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
				c.Cache[k] = make(map[interface{}]*Row)

			}
			kv := asString(reflect.ValueOf(table).Elem().FieldByName(k).Interface())

			r := &Row{
				Value: table,
				Mu:    &sync.RWMutex{},
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

func (c *Table) clean() {
	// 清除过期table
	if len(c.Keys) == 0 {
		panic(ErrorNoKey)
	}
	for _, key := range c.Keys {

		for k, v := range c.Cache[key] {
			fmt.Println("key:", key)
			fmt.Println("value:", k)
			fmt.Printf("vvvvv%+v\n", v)
			fmt.Printf("vvvvv%+v\n", v.Value)
			now := time.Now()
			fmt.Println("now", now)
			fmt.Printf("vvvvv%v\n", now.Sub(v.Expire).Seconds())
			if v.CanExpire && time.Now().Sub(v.Expire).Seconds() >= float64(0) {
				if v.Mu != nil {
					v.Mu.Lock()
					delete(c.Cache[key], k)
					v.Mu.Unlock()
				} else {
					delete(c.Cache[key], k)
				}

				// f, err := c.Filter(key, k)
				// fmt.Printf("ffffff%+v\n", f)
				// if err != nil {
				// 	continue
				// }
				fmt.Println("deleleleleleleleletetete")
				// f.Del()
			}
		}
	}
	// 第一个字段就行了

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
