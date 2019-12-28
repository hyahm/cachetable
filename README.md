# cache 
go 缓存表   
如其名， 缓存的一张表， 当然可以多张， New就好    
老方法缓存， 建立很多map， 为了反向找到通常会这样  
map[int64]string    // 用户id 对应用户名  
map[string]int64    // 用户名对应id  
修改其中一个map， 另外一个map也要修改  

现在的话， 直接使用 struct保存此类数据， 设置key， 应为这2个都要对应， 所以要设置这2个，  
后面不管是修改还是查找， 使用set或get即可， 使用到reflect， 效率肯定没多map快  

# demo 

```
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
	c := cachetable.NewTable(people{})    // 导入表结构， 可以是指针也可以是结构

	if err := c.SetKey("Id"); err != nil {    // 设置主键的值， 必须设置， 相当于数据库的唯一索引， 名字是结构体的字段名
		panic(err)
	}
	if err := c.Add(u); err != nil {     // 添加数据， 如果不是一样的结构体会报错
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

	if err := c.Set("Age", 555, "Id", 1); err != nil {   // 根据某个索引键值修改某个字段的值
		panic(err)
	}
	value, err := c.Get("Age", "Id", 1)     // 根据某个索引键值获取某个键值
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

}


```
输出
```
555
```
