package main

import (
	"fmt"
	"reflect"
)

type Table struct {
	Id     int64
	Name   string
	Age    int
	Weight int
}

type Cacher interface {

	//cache 一张表
}

type Cache struct {
	keys  map[string]interface{}
	table map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{
		table: make(map[string]interface{}),
		keys:  make(map[string]interface{}),
	}
}

func (c *Cache) Load(name string, table interface{}) *Cache {
	// 加载数据
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		panic("load value is not a pointer")
	}
	menu := reflect.TypeOf(table)
	mem := menu.Elem().NumField()
	for i := 0; i < mem; i++ {
		fmt.Println(menu.Field(i))
	}
	fmt.Println(menu.Field(mem))
	//fmt.Println(reflect.ValueOf(&table).FieldByName("Name"))
	//判断table的值， 必须是struct 指针
	c.table[name] = table
	return c
}

func (c *Cache) Key(interface{}) *Cache {
	// 加载数据
	return c
}

func main() {
	c := NewCache()
	c.Load("zth", &Table{
		Id:     345,
		Name:   "asdf",
		Weight: 3224,
		Age:    13,
	})
}
