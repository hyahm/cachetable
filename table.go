package cachetable

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
	"time"
)





type Cache struct {
	keys  []string                   // 保存key, 为了去重， 使用map
	cache map[string]map[string]*row // 保存field， 将所有值都转为string
	s     interface{}                // 保存表结构

}

func NewTable(table interface{}) *Cache {
	// 表字段不能是指针或结构
	return &Cache{
		keys:  make([]string, 0),
		cache: make(map[string]map[string]*row),
		s:     table,
	}

}

func (c *Cache) Add(table interface{}, expire time.Duration) error {
	//必须是指针
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return ErrorNotPointer
	}
	if len(c.keys) == 0 {
		return ErrorNoKey
	}
	var st reflect.Type
	if reflect.TypeOf(c.s).Kind() == reflect.Ptr {
		st = reflect.TypeOf(c.s).Elem()
	} else {
		st = reflect.TypeOf(c.s)
	}
	// 必须是同一类型
	if st == reflect.TypeOf(table).Elem() {
		//遍历 添加key
		for _, k := range c.keys {
			// 将字段的值全部转化为string

			if _, ok := c.cache[k]; !ok {
				// 没有字段， 初始化
				c.cache[k] = make(map[string]*row)

			}

			kv, err := c.toString(reflect.ValueOf(table).Elem().FieldByName(k).Interface())
			if err != nil {
				return err
			}
			r := &row{
				mu:     sync.RWMutex{},
				value:  table,

			}
			if expire > 0 {
				r.Expire = time.Now().Add(expire)
				r.CanExpire = true
			}

			r.mu.Lock()
			c.cache[k][kv] = r
			r.mu.Unlock()
		}

	} else {
		return ErrorStruct
	}
	return nil
}

func (c *Cache) toString(value interface{}) (string, error) {
	t := reflect.TypeOf(value).String()
	fv := reflect.ValueOf(value)
	//fmt.Println("type:", ft.String())
	switch t {
	case "string":
		return fv.Interface().(string), nil
	case "int":
		return strconv.Itoa(fv.Interface().(int)), nil
	case "int64":
		return strconv.FormatInt(fv.Interface().(int64), 10), nil
	case "uint64":
		return strconv.FormatUint(fv.Interface().(uint64), 10), nil
	case "bool":
		return strconv.FormatBool(fv.Interface().(bool)), nil
	case "float64":
		return strconv.FormatFloat(fv.Interface().(float64), 'f', -1, 64), nil

	default:
		return "", errors.New("not support type")
	}
}

func (c *Cache) SetKeys(keys ...string) error {
	if c.s == nil {
		return ErrorNotInit
	}
	var sv reflect.Type
	if reflect.TypeOf(c.s).Kind() == reflect.Ptr {
		sv = reflect.TypeOf(c.s).Elem()
	} else {
		sv = reflect.TypeOf(c.s)
	}
	// 判断key 是否有效
	for _, k := range keys {
		if _, ok := sv.FieldByName(k); ok {
			if !c.hasKey(k) {
				c.keys = append(c.keys, k)
			}

		}
	}

	return nil
}

//

//
func (c *Cache) GetKeys() (ks []string) {
	return c.keys
}

func (c *Cache) hasKey(s string) bool {
	for _, v := range c.keys {
		if v == s {
			return true
		}
	}
	return false
}

func (c *Cache) Clean(t time.Duration) {
	// 清除过期table
	if len(c.keys) == 0 {
		panic(ErrorNoKey)
	}
	for {
		// 第一个字段就行了
		time.Sleep(t)
		allmap := c.cache[c.keys[0]]
		for k, v := range allmap {
			if !v.CanExpire && time.Now().Sub(v.Expire) >= 0 {
				c.Filter(c.keys[0], k).Del()
			}
		}
	}
}
