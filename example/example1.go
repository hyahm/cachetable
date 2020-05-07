package main

import (
	"fmt"
	"github.com/hyahm/cachetable/example/aaa"
	"github.com/hyahm/cachetable"
	"log"
)

func main() {
	ct := cachetable.NewCT()
	err := ct.Load("aa.txt", aaa.People{}, aaa.Teacher{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ct.ShowTables())
}