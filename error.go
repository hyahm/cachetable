package cachetable

import (
	"errors"
	"strconv"
)

var (
	ErrorNotInit       = errors.New("init first")
	ErrorNotPointer    = errors.New("table must be a pointer")
	ErrorNoKey         = errors.New("at least set one key")
	ErrorStruct        = errors.New("not a same struct")
	ErrorDuplicate     = errors.New("Duplicate key ")
	ErrorExpired       = errors.New("row expired")
	ErrorNoFeildKey    = errors.New("field not a key")
	ErrorNotFoundValue = errors.New("not found value")
	ErrorNoRows        = errors.New("not rows")
	ErrorTypeNoMatch   = errors.New("type mismatch")
	ErrorLengthNoMatch = errors.New("length mismatch")
	errNilPtr          = errors.New("destination pointer is nil")
)

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}
