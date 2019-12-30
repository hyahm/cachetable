package cachetable

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
)

type row struct {
	mu    sync.RWMutex // 行锁
	value interface{}  // 值
}

var (
	ErrorNotInit    = errors.New("init first")
	ErrorNotPointer = errors.New("table must be a pointer")
	ErrorNoKey      = errors.New("at least set one key")
	ErrorStruct     = errors.New("not a same struct")
	ErrorDuplicate  = errors.New("Duplicate key ")
	//ErrorStructFeild = errors.New("Struct need Not Have ptr struct")
	ErrorNoFeildKey = errors.New("field not a key")
	ErrorNoRows     = errors.New("not rows")
)

type Cache struct {
	keys  map[string]bool            // 保存key, 为了去重， 使用map
	cache map[string]map[string]*row // 保存field， 将所有值都转为string
	s     interface{}                // 保存表结构
}

func NewTable(table interface{}) *Cache {
	// 表字段不能是指针或结构
	return &Cache{
		keys:  make(map[string]bool),
		cache: make(map[string]map[string]*row),
		s:     table,
	}

}

func (c *Cache) Add(table interface{}) error {
	//必须是指针
	if reflect.TypeOf(table).Kind() != reflect.Ptr {
		return ErrorNotPointer
	}
	if len(c.keys) == 0 {
		return ErrorNoKey
	}
	var st reflect.Type
	if reflect.TypeOf(c.s).Kind() == reflect.Ptr {
		st = reflect.TypeOf(c.s).Elem()
	} else {
		st = reflect.TypeOf(c.s)
	}
	// 必须是同一类型
	if st == reflect.TypeOf(table).Elem() {
		//遍历 添加key
		for k, _ := range c.keys {
			// 将字段的值全部转化为string

			if _, ok := c.cache[k]; !ok {
				// 没有字段， 初始化

				c.cache[k] = make(map[string]*row)

			}

			kv, err := c.toString(reflect.ValueOf(table).Elem().FieldByName(k).Interface())
			if err != nil {
				return err
			}
			r := &row{
				mu:    sync.RWMutex{},
				value: table,
			}
			r.mu.Lock()
			c.cache[k][kv] = r
			r.mu.Unlock()
		}

	} else {
		return ErrorStruct
	}
	return nil
}

func (c *Cache) toString(value interface{}) (string, error) {
	t := reflect.TypeOf(value).String()
	fv := reflect.ValueOf(value)
	//fmt.Println("type:", ft.String())
	switch t {
	case "string":
		return fv.Interface().(string), nil
	case "int":
		return strconv.Itoa(fv.Interface().(int)), nil
	case "int64":
		return strconv.FormatInt(fv.Interface().(int64), 10), nil
	case "uint64":
		return strconv.FormatUint(fv.Interface().(uint64), 10), nil
	case "bool":
		return strconv.FormatBool(fv.Interface().(bool)), nil
	case "float64":
		return strconv.FormatFloat(fv.Interface().(float64), 'f', -1, 64), nil

	default:
		return "", errors.New("not support type")
	}
}

func (c *Cache) SetKeys(keys ...string) error {
	if c.s == nil {
		return ErrorNotInit
	}
	var sv reflect.Type
	if reflect.TypeOf(c.s).Kind() == reflect.Ptr {
		sv = reflect.TypeOf(c.s).Elem()
	} else {
		sv = reflect.TypeOf(c.s)
	}
	// 判断key 是否有效
	for _, k := range keys {
		if _, ok := sv.FieldByName(k); ok {
			c.keys[k] = true
		}
	}

	return nil
}

//

//
func (c *Cache) GetKeys() (ks []string) {
	for k, _ := range c.keys {
		ks = append(ks, k)
	}
	return
}
