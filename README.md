# cache 
go 缓存表   多个key 对应一行
如其名， 缓存的一张表， 当然可以多张， New就好    
老方法缓存， 建立很多map， 为了反向找到通常会这样  
map[int64]string    // 用户id 对应用户名  
map[string]int64    // 用户名对应id  
修改其中一个map， 另外一个map也要修改  

现在的话， 直接使用 struct保存此类数据， 设置key， 应为这2个都要对应， 所以要设置这2个，  
后面不管是修改还是查找， 使用set或get即可， 使用到reflect， 效率肯定没多map快  
初始包， 结构体的key只支持int, string, uint64, int64, bool,fload64 的值  


增加过期时间， 代替简单的缓存

支持load 和 save
# demo 

```
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
	fmt.Println(filter.Row().(*aaa.People))
	fmt.Println(c.Columns("Age"))
}

```

建议新开一个goroutine 删除过期的row
```
go func(){
	func (c *Cache) Clean(t time.Duration) {   // t表示检查的时间间隔
}()
```
