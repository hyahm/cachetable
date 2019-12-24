package main

import (
	"errors"
	"fmt"
	"log"
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
	_, iskey := c.keys[key]
	if f, ok := c.cache[key]; ok {
		//
		if v, ok := f[value]; ok {
			fmt.Println(v)
			reflect.ValueOf(v).Elem().FieldByName(setkey).Set(reflect.ValueOf(setvalue))
			if iskey {
				// 如果是主键， 更新map
				f[setvalue] = v
				delete(f, value)
			}
			fmt.Println(v)
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
		fmt.Println(f)
		if v, ok := f[value]; ok {
			fmt.Println(v)
			return reflect.ValueOf(v).Elem().FieldByName(key).Interface(), nil
		}
	} else {
		fmt.Println("1111")
		return nil, errors.New("field not a key")
	}
	//

	return nil, nil
}

// 添加了key， 那么就无法删除了
type people struct {
	Name string
	Age  int
	Id   int
}

func main() {
	u := &people{
		Name: "2222",
		Age:  111,
		Id:   0,
	}
	u1 := &people{
		Name: "2222",
		Age:  111,
		Id:   1,
	}
	u2 := &people{
		Name: "2",
		Age:  222,
		Id:   2,
	}
	u3 := &people{
		Name: "2222",
		Age:  333,
		Id:   3,
	}
	u4 := &people{
		Name: "2222",
		Age:  444,
		Id:   4,
	}
	u5 := &people{
		Name: "2222",
		Age:  555,
		Id:   5,
	}
	c := NewCache()
	c.Table(&people{})

	if err := c.Key("Name"); err != nil {
		log.Fatal(err)
	}
	if err := c.Key("Id"); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u1); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u2); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u3); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u4); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u5); err != nil {
		log.Fatal(err)
	}
	if err := c.Add(u5); err != nil {
		log.Fatal(err)
	}

	//if err := c.Set("Name", "111", "Age", 222); err != nil {
	//	log.Fatal(err)
	//}
	value, err := c.GetValue("Age", "Id", "2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)

}
