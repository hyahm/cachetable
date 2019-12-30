# cache 
go 缓存表   多个key 对应一行
如其名， 缓存的一张表， 当然可以多张， New就好    
老方法缓存， 建立很多map， 为了反向找到通常会这样  
map[int64]string    // 用户id 对应用户名  
map[string]int64    // 用户名对应id  
修改其中一个map， 另外一个map也要修改  

现在的话， 直接使用 struct保存此类数据， 设置key， 应为这2个都要对应， 所以要设置这2个，  
后面不管是修改还是查找， 使用set或get即可， 使用到reflect， 效率肯定没多map快  
初始包， 结构体只支持int, string, uint64, int64, bool,fload64 的值
# demo 

```
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

```
输出
```
[2222]
[hello world]
[6666]

```
