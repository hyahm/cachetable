package cachetable

import (
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
	Age = "Age"
	Id = "Id"
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
	table := NewTable(cat{})
	table.SetKeys(Name, Id)

	table.Add(t1, 0*time.Second)
	table.Add(t2, 0*time.Second)

	f := table.Filter(Id , 1)
	err := f.Set(Age, "asdf")
	if err == nil {
		t.Error("need error")
	}
}
