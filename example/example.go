package main

import (
	"fmt"
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

	// u := &aaa.People{
	// 	Name: "2222",
	// 	Age:  888,
	// 	Id:   0,
	// 	Data: []byte("hello world"),
	// }
	u1 := &aaa.Teacher{
		Name: "1111",
		Age:  111,
		Id:   1,
	}
	u2 := &aaa.Teacher{
		Name: "2222",
		Age:  222,
		Id:   2,
	}
	// u3 := &aaa.People{
	// 	Name: "2222",
	// 	Age:  333,
	// 	Id:   3,
	// }
	// u4 := &aaa.People{
	// 	Name: "2222",
	// 	Age:  444,
	// 	Id:   4,
	// }
	// u5 := &aaa.People{
	// 	Name: "2222",
	// 	Age:  555,
	// 	Id:   5,
	// }

	ct := cachetable.NewCT()
	ct.CreateTable("teacher", &aaa.Teacher{})

	t, _ := ct.Use("teacher")
	err := t.SetKeys(Id, Age)
	if err != nil {
		panic(err)
	}
	// 清东西

	fmt.Println(time.Now())
	go ct.Clean(time.Second * 2)
	t.Add(u1, 2*time.Second)
	t.Add(u2, 0)

	// 获取值

	// filter, err := c.Filter(Id, 3)
	// if err != nil {
	// 	panic(err)
	// }
	// value := filter.Get(Name)
	// fmt.Println(value)
	// // 设置非 key的value
	// err = filter.Set(Name, "hello world")
	// if err != nil {
	// 	panic(err)
	// }
	// var a string
	// value = filter.Get(Name)
	// err = value.Scan(&a)
	// if err != nil {
	// 	panic(err)
	// }

	// // 设置 key的value
	// err = filter.Set(Age, 6666)
	// if err != nil {
	// 	panic(err)
	// }
	time.Sleep(10 * time.Second)
	// var age string
	// err = filter.Get(Age).Scan(&age)

	// fmt.Println(age)
	// fmt.Println(filter.TTL())
	// fmt.Println(filter.Row().(*aaa.People))
	// fmt.Println(c.Columns("Age"))
	// for key, v := range ct["teacher"].Cache["Id"] {
	// 	fmt.Println("-------------------")
	// 	fmt.Println(v.Expire)
	// 	fmt.Println(key)
	// 	fmt.Printf("%+v\n", v.Value)
	// 	fmt.Println(v.CanExpire)
	// 	fmt.Println("-------------------")
	// }

	fmt.Println(ct["teacher"].Cache["Age"])
	fmt.Println(ct["teacher"].Cache["Id"])
}
