package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hyahm/cachetable"
	"github.com/hyahm/cachetable/example/aaa"
)

// 添加了key， 那么就无法删除了

const (
	Name = "Name"
	Age  = "Age"
	Id   = "Id"
)

func main() {
	u := &aaa.People{
		Name: "2222",
		Age:  888,
		Id:   0,
		Data: []byte("hello world"),
	}
	u1 := &aaa.Teacher{
		Name: "2222",
		Age:  111,
		Id:   1,
	}
	u2 := &aaa.Teacher{
		Name: "2",
		Age:  222,
		Id:   2,
	}
	u3 := &aaa.People{
		Name: "2222",
		Age:  333,
		Id:   3,
	}
	u4 := &aaa.People{
		Name: "2222",
		Age:  444,
		Id:   4,
	}
	u5 := &aaa.People{
		Name: "2222",
		Age:  555,
		Id:   5,
	}

	ct := cachetable.NewCT()
	ct.CreateTable("me", &aaa.People{})
	ct.CreateTable("teacher", &aaa.Teacher{})
	c, _ := ct.Use("me")
	err := c.SetKeys(Id, Age)
	if err != nil {
		panic(err)
	}
	t, _ := ct.Use("teacher")
	err = t.SetKeys(Id, Age)
	if err != nil {
		panic(err)
	}

	err = c.Add(u, 0)
	if err != nil {
		log.Fatal(err)
	}
	t.Add(u1, 10*time.Second)
	t.Add(u2, 0)
	c.Add(u3, 0)
	c.Add(u4, 0)

	c.Add(u5, 0)
	// 获取值

	filter, err := c.Filter(Id, 3)
	if err != nil {
		panic(err)
	}
	value := filter.Get(Name)
	fmt.Println(value)
	// 设置非 key的value
	err = filter.Set(Name, "hello world")
	if err != nil {
		panic(err)
	}
	var a string
	value = filter.Get(Name)
	err = value.Scan(&a)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)

	// 设置 key的value
	err = filter.Set(Age, 6666)
	if err != nil {
		panic(err)
	}

	var age string
	err = filter.Get(Age).Scan(&age)

	fmt.Println(age)
	fmt.Println(filter.TTL())
	fmt.Println(filter.Table().(*aaa.People))
	fmt.Println(c.Columns("Age"))
}
