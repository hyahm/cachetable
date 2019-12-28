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

type student struct {
	Name string
	Age  int
	Id   int
}

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
	c := cachetable.NewTable(people{})

	if err := c.SetKey("Id"); err != nil {
		panic(err)
	}

	c.SetKey("Age")
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

	if err := c.Set("Age", 777, "Id", 1); err != nil {
		panic(err)
	}
	value, err := c.Get("Age", "Id", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

	c.Del("Id", 1)

	value, err = c.Get("Age", "Id", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

	value, err = c.Get("Name", "Age", 777)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

}
