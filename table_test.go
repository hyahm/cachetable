package cachetable

import (
	"fmt"
	"testing"
	"time"
)

type cat struct {
	Name string
	Age  int
	Id   int
}

const (
	Name = "Name"
	Age  = "Age"
	Id   = "Id"
)

func TestTable(t *testing.T) {
	t1 := &cat{
		Id:   0,
		Name: "lucky",
		Age:  1,
	}

	t2 := &cat{
		Name: "mimi",
		Age:  1,
		Id:   1,
	}
	ct := NewCT()
	ct.Add("t1", cat{})

	ct.Table("t1").SetKeys(Name, Id)

	ct.Table("t1").Add(t1, 0*time.Second)
	ct.Table("t1").Add(t2, 0*time.Second)

	f := ct.Table("t1").Filter(Id, 1)
	err := f.Set(Age, "asdf")
	if err == nil {
		t.Error("need error")
	}
	fmt.Println(f.Get(Age))
}
