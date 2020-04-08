package cachetable

import (
	"reflect"
	"sync"
	"time"
)

type Filter struct {
	row *row
	c   *Cache
	Err error
}

func (c *Cache) Filter(field string, value interface{}) *Filter {
	if c.s == nil {
		return &Filter{
			row: nil,
			Err: ErrorNotInit,
		}

	}
	if len(c.keys) == 0 {
		return &Filter{
			row: nil,
			Err: ErrorNoKey,
		}

	}
	// 找到所有索引， 删除,   必须是key
	if vms, ok := c.cache[field]; ok {
		//找到所有所有的keys 的值
		key := asString(value)
		if vms[key].canExpire && time.Now().Sub(vms[key].expire) >= 0 {
			// 说明过期了
			//直接先删掉

			f := &Filter{
				row: vms[key],
				Err: ErrorExpired,
				c:   c,
			}
			f.row.mu.Lock()
			f.Del()
			f.row.mu.Unlock()
			f.row = nil
			return f
		}
		return &Filter{
			row: vms[key],
			Err: nil,
			c:   c,
		}
		// return nil
	} else {
		return &Filter{
			row: nil,
			Err: ErrorNoFeildKey,
		}
	}

}

func (f *Filter) expired() bool {
	return f.TTL() == 0
}

func (f *Filter) Expired() bool {
	return f.expired()
}

func (f *Filter) TTL() time.Duration {
	if f.row == nil {
		return 0
	}
	if f.row.canExpire {
		if time.Now().Sub(f.row.expire) >= 0 {
			return 0
		} else {
			return f.row.expire.Sub(time.Now())
		}

	} else {
		return -1
	}
}

func (f *Filter) Get(keys ...string) *Result {
	rl := &Result{}
	if f.Err != nil {
		rl.err = f.Err
		return rl
	}
	l := len(keys)
	rl.values = make([]interface{}, l)
	if f.expired() {
		return rl
	}

	for i, v := range keys {
		val := reflect.ValueOf(f.row.value).Elem().FieldByName(v).Interface()
		rl.values[i] = val
	}
	return rl
}

func (f *Filter) Del() error {
	if f.row == nil {
		return f.Err
	}
	if f.expired() {
		return ErrorExpired
	}
	f.row.mu.Lock()
	defer f.row.mu.Unlock()
	//找到所有所有的keys 的值
	for _, k := range f.c.keys {
		//v := c.Get(k)
		ft := reflect.ValueOf(f.row.value).FieldByName(k).Interface()
		value := asString(ft)
		delete(f.c.cache[k], value)
	}

	return f.Err

}

func (f *Filter) SetTTL(t time.Duration) error {
	if f.row == nil {
		return f.Err
	}

	if f.expired() {
		return ErrorExpired
	}
	if t <= 0 {
		f.row.canExpire = false
		return nil
	}
	if t > 0 {
		f.row.canExpire = true
		f.row.expire = time.Now().Add(t)
		return nil
	}

	return f.Err
}

func (f *Filter) Set(field string, value interface{}) error {
	if f.row == nil {
		return f.Err
	}

	if f.expired() {
		return ErrorExpired
	}
	f.row.mu.Lock()
	defer f.row.mu.Unlock()

	if ok := f.c.hasKey(field); ok {
		// 如果是key, value不能重复
		newvalue_str := asString(value)
		if _, ok := f.c.cache[field][newvalue_str]; ok {
			return ErrorDuplicate
		}

		// 如果是设置的是key是主键， 重新生成
		oldvalue_str := asString(value)

		f.c.cache[field][newvalue_str] = &row{
			mu:    sync.RWMutex{},
			value: f.row,
		}
		// 删掉老的键值
		delete(f.c.cache[field], oldvalue_str)
	}

	// 更新v
	newv := reflect.ValueOf(value)
	setv := reflect.ValueOf(f.row.value).Elem().FieldByName(field)
	if newv.Type().String() != setv.Type().String() {
		return ErrorTypeNoMatch
	}
	setv.Set(newv)

	return f.Err
}
