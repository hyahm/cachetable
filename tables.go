package cachetable

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
	"time"
)

var mmu sync.RWMutex
var cmu sync.RWMutex

type CT struct {
	Data map[string]*Cache
	Ttl  time.Duration
}

// func NewTable(table interface{}) *Cache {
// 	// 表字段不能是指针或结构
// 	return &Cache{
// 		keys:  make([]string, 0),
// 		cache: make(map[string]map[string]*row),
// 		s:     table,
// 	}

// }

func NewCT() *CT {
	mmu = sync.RWMutex{}
	cmu = sync.RWMutex{}
	return &CT{
		Ttl:  5 * time.Second,
		Data: make(map[string]*Cache),
	}
}

func log(name *string) {
	fmt.Printf("key %v have already exsit \n", name)
}

func (ct *CT) Add(name string, table interface{}) {

	if _, ok := ct.Data[name]; ok {
		runtime.SetFinalizer(&name, log)
		return
	}
	ct.Data[name] = &Cache{
		Keys:  make([]string, 0),
		Cache: make(map[string]map[string]*Row),
		S:     table,
	}

}

func (ct *CT) Delete(name string) {

	if _, ok := ct.Data[name]; ok {
		delete(ct.Data, name)
	}
}

func (ct *CT) Exsit(name string) (ok bool) {

	_, ok = ct.Data[name]
	return
}

func (ct *CT) Table(name string) *Cache {
	return ct.Data[name]
}

// 清除过期的key
func (ct *CT) Clean(t time.Duration) {
	for _, v := range ct.Data {
		v.clean(ct.Ttl)
	}
}

func (ct *CT) Save(filename string) error {

	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	// defer func() {
	// 	if x := recover(); x != nil {
	// 		fmt.Println("Error registering item types with Gob library")
	// 	}
	// }()
	mmu.Lock()
	defer mmu.Unlock()

	gob.Register(&CT{})
	for _, v := range ct.Data {
		gob.Register(v.S)
	}

	err := enc.Encode(ct)
	if err != nil {
		return err
	}
	fmt.Println(w.String(), 111)

	return ioutil.WriteFile(filename, w.Bytes(), 0644)
}
