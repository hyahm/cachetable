package cachetable

import (
	"fmt"
	"log"
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
	ct.CreateTable("t1", cat{})

	T1, err := ct.Use("t1")
	if err != nil {
		log.Fatal(err)
	}
	T1.SetKeys(Name, Id)

	T1.Add(t1, 0*time.Second)
	T1.Add(t2, 0*time.Second)

	f, err := T1.Filter(Id, 1)
	if err == nil {
		t.Error("need error")
	}
	err = f.Set(Age, "asdf")
	if err == nil {
		t.Error("need error")
	}
	fmt.Println(f.Get(Age))
}

func TestSave(t *testing.T) {
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
	ct.CreateTable("t1", cat{})
	T1, err := ct.Use("t1")
	if err != nil {
		log.Fatal(err)
	}
	T1.SetKeys(Name, Id)

	T1.Add(t1, 0*time.Second)
	T1.Add(t2, 0*time.Second)

	f, err := T1.Filter(Id, 1)
	if err == nil {
		t.Error("need error")
	}
	err = f.Set(Age, "asdf")

	err = ct.Save("aa.txt")
	if err != nil {
		t.Log(err)
	}
}
