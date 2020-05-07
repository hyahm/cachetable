package cachetable

import (
	"reflect"
	"time"
)

type Filter struct {
	Row *Row
	c   *Cache
	Err error
}

func (c *Cache) Filter(field string, value interface{}) *Filter {
	if c.S == nil {
		return &Filter{
			Row: nil,
			Err: ErrorNotInit,
		}

	}
	if len(c.Keys) == 0 {
		return &Filter{
			Row: nil,
			Err: ErrorNoKey,
		}

	}
	// 找到所有索引， 删除,   必须是key
	if vms, ok := c.Cache[field]; ok {
		//找到所有所有的keys 的值
		key := asString(value)
		if vms[key].CanExpire && time.Now().Sub(vms[key].Expire) >= 0 {
			// 说明过期了
			//直接先删掉

			f := &Filter{
				Row: vms[key],
				Err: ErrorExpired,
				c:   c,
			}
			cmu.Lock()
			f.Del()
			cmu.Unlock()
			f.Row = nil
			return f
		}
		return &Filter{
			Row: vms[key],
			Err: nil,
			c:   c,
		}
		// return nil
	} else {
		return &Filter{
			Row: nil,
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
	if f.Row == nil {
		return 0
	}
	if f.Row.CanExpire {
		if time.Now().Sub(f.Row.Expire) >= 0 {
			return 0
		} else {
			return f.Row.Expire.Sub(time.Now())
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
		val := reflect.ValueOf(f.Row.Value).Elem().FieldByName(v).Interface()
		rl.values[i] = val
	}
	return rl
}

func (f *Filter) Del() error {
	if f.Row == nil {
		return f.Err
	}
	if f.expired() {
		return ErrorExpired
	}
	cmu.Lock()
	defer cmu.Unlock()
	//找到所有所有的keys 的值
	for _, k := range f.c.Keys {
		//v := c.Get(k)
		ft := reflect.ValueOf(f.Row.Value).FieldByName(k).Interface()
		value := asString(ft)
		delete(f.c.Cache[k], value)
	}

	return f.Err

}

func (f *Filter) SetTTL(t time.Duration) error {
	// 某一条数据设置过期时间
	if f.Row == nil {
		return f.Err
	}

	if f.expired() {
		return ErrorExpired
	}
	if t <= 0 {
		f.Row.CanExpire = false
		return nil
	}
	if t > 0 {
		f.Row.CanExpire = true
		f.Row.Expire = time.Now().Add(t)
		return nil
	}

	return f.Err
}

func (f *Filter) Set(field string, value interface{}) error {
	if f.Row == nil {
		return f.Err
	}

	if f.expired() {
		return ErrorExpired
	}
	cmu.Lock()
	defer cmu.Unlock()

	if ok := f.c.hasKey(field); ok {
		// 如果是key, value不能重复
		newvalue_str := asString(value)
		if _, ok := f.c.Cache[field][newvalue_str]; ok {
			return ErrorDuplicate
		}

		// 如果是设置的是key是主键， 重新生成
		oldvalue_str := asString(value)

		f.c.Cache[field][newvalue_str] = &Row{
			Value: f.Row,
		}
		// 删掉老的键值
		delete(f.c.Cache[field], oldvalue_str)
	}

	// 更新v
	newv := reflect.ValueOf(value)
	setv := reflect.ValueOf(f.Row.Value).Elem().FieldByName(field)
	if newv.Type().String() != setv.Type().String() {
		return ErrorTypeNoMatch
	}
	setv.Set(newv)

	return f.Err
}
