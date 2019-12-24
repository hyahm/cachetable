package main

import (
	"fmt"
	"github.com/hyahm/cache"
)

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
	c := cache.NewCache()
	c.Table(&people{})

	if err := c.Key("Id"); err != nil {
		panic(err)
	}
	if err := c.Add(u); err != nil {
		panic(err)
	}
	if err := c.Add(u1); err != nil {
		panic(err)
	}
	if err := c.Add(u2); err != nil {
		panic(err)
	}
	if err := c.Add(u3); err != nil {
		panic(err)
	}
	if err := c.Add(u4); err != nil {
		panic(err)
	}
	if err := c.Add(u5); err != nil {
		panic(err)
	}

	if err := c.Set("Id", 1, "Age", 222); err != nil {
		panic(err)
	}
	value, err := c.GetValue("Age", "Id", 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

}

