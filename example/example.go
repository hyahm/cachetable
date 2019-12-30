package main

import (
	"fmt"
	"github.com/hyahm/cachetable"
)

// 添加了key， 那么就无法删除了
type people struct {
	Name string
	Age  int
	Id   int
}

const (
	Name = "Name"
	Age  = "Age"
	Id   = "Id"
)

func main() {
	u := &people{
		Name: "2222",
		Age:  888,
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
	c, err := cachetable.NewTable(people{})
	if err != nil {
		fmt.Println(err)
	}

	if err := c.SetKeys(Id, Age); err != nil {
		panic(err)
	}

	//c.SetKey("Age")

	c.Add(u)
	c.Add(u1)
	c.Add(u2)
	c.Add(u3)
	c.Add(u4)

	c.Add(u5)
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
