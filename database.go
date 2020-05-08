package cachetable

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var ctmu sync.RWMutex
var cmu sync.RWMutex

type CT map[string]*Table

func NewCT() CT {
	ctmu = sync.RWMutex{}
	cmu = sync.RWMutex{}
	return make(map[string]*Table)

}

func (ct CT) Add(name string, table interface{}) error {
	ctmu.Lock()
	defer ctmu.Unlock()
	if _, ok := ct[name]; ok {
		return errors.New("key " + name + " have already exsit \n")
	}

	ct[name] = &Table{
		Keys:  make([]string, 0),
		Cache: make(map[string]map[string]*Row),
		S:     table,
	}
	return nil
}

func (ct CT) Delete(name string) {
	ctmu.Lock()
	defer ctmu.Unlock()
	if _, ok := ct[name]; ok {
		delete(ct, name)
	}
}

func (ct CT) Exsit(name string) (ok bool) {
	ctmu.Lock()
	defer ctmu.Unlock()
	_, ok = ct[name]
	return
}
func (ct CT) ShowTables() (tablesname []string) {
	ctmu.Lock()
	defer ctmu.Unlock()
	for name, _ := range ct {
		tablesname = append(tablesname, name)
	}
	return
}

func (ct CT) Table(name string) (*Table, error) {
	ctmu.Lock()
	defer ctmu.Unlock()
	if v, ok := ct[name]; ok {
		return v, nil
	} else {
		return nil, ErrorNoFeildKey
	}
}

func (ct CT) Clean(expire time.Duration) {
	ctmu.Lock()
	defer ctmu.Unlock()
	for _, v := range ct {
		v.clean(expire)
	}
}

func (ct CT) Save(filename string) error {

	var w bytes.Buffer
	enc := gob.NewEncoder(&w)

	gob.Register(&ct)
	for _, v := range ct {
		gob.Register(v.S)
	}

	err := enc.Encode(ct)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, w.Bytes(), 0644)
}

func (ct CT) Load(filename string, table ...interface{}) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(r)

	for _, v := range table {
		gob.Register(v)
	}
	err = dec.Decode(&ct)
	if err != nil {
		return err
	}

	return nil
}
