package cachetable

import (
	"sync"
	"time"
)

type row struct {
	mu     sync.RWMutex // 行锁
	value  interface{}  // 值
	expire time.Time
	canExpire bool
}