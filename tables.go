package cachetable

import (
	"fmt"
	"runtime"
)

type CT map[string]*Cache

// func NewTable(table interface{}) *Cache {
// 	// 表字段不能是指针或结构
// 	return &Cache{
// 		keys:  make([]string, 0),
// 		cache: make(map[string]map[string]*row),
// 		s:     table,
// 	}

// }

func NewCT() CT {
	return make(map[string]*Cache)
}

func log(name *string) {
	fmt.Printf("key %v have already exsit \n", name)
}

func (ct CT) Add(name string, table Table) {
	
	if _, ok := ct[name]; ok {
		runtime.SetFinalizer(&name, log)
		return
	}
	ct[name] = &Cache{
		keys:  make([]string, 0),
		cache: make(map[string]map[string]*row),
		s:     table,
	}

}

func (ct CT) Delete(name string) {
	
	if _, ok := ct[name]; ok {
		delete(ct, name)
	}
}

func (ct CT) Exsit(name string) (ok bool) {
	
	_, ok = ct[name]
	return 
}

func (ct CT) Table(name string) *Cache {
	return ct[name] 
}


// 清除过期的key
func (ct CT) Clean(t time.Duration) {
	for k, v := range ct {
		v.clean()
	}
}