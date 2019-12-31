package cachetable

import "errors"

var (
	ErrorNotInit    = errors.New("init first")
	ErrorNotPointer = errors.New("table must be a pointer")
	ErrorNoKey      = errors.New("at least set one key")
	ErrorStruct     = errors.New("not a same struct")
	ErrorDuplicate  = errors.New("Duplicate key ")
	ErrorExpired = errors.New("row expired")
	ErrorNoFeildKey = errors.New("field not a key")
	ErrorNoRows     = errors.New("not rows")
	ErrorTypeNoMatch     = errors.New("type mismatch")

)
