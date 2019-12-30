package main

import (
	"fmt"
	"github.com/hyahm/cachetable"
	"log"
	"time"
)

// 添加了key， 那么就无法删除了
type people struct {
	Name    string
	Age     int
	Id      int
	Kecheng []string
	Data    []byte
}

const (
	Name = "Name"
	Age  = "Age"
	Id   = "Id"
)

func main() {
	var t time.Time
	fmt.Println(t.Unix())
	u := &people{
		Name: "2222",
		Age:  888,
		Id:   0,
		Data: []byte("hello world"),
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
	c := cachetable.NewTable(people{})

	if err := c.SetKeys(Id, Age); err != nil {
		panic(err)
	}

	err := c.Add(u, 0)
	if err != nil {
		log.Fatal(err)
	}
	c.Add(u1, 0)
	c.Add(u2, 0)
	c.Add(u3, 0)
	c.Add(u4, 0)

	c.Add(u5, 0)
	// 获取值
	filter := c.Filter(Id, 1)
	value := filter.Get(Name)
	fmt.Println(value)
	// 设置非 key的value
	err = filter.Set(Name, "hello world")
	if err != nil {
		panic(err)
	}
	value = filter.Get(Name)
	fmt.Println(value)

	// 设置 key的value
	err = filter.Set(Age, 6666)
	if err != nil {
		panic(err)
	}
	value = filter.Get(Age)
	fmt.Println(value)

}
