package cachetable

import "testing"

type cat struct {
	Name string
	Age  int
	Id   int
}

func TestTable(t *testing.T) {
	t1 := &cat{
		Id:   0,
		Name: "lucky",
		Age:  1,
	}
	tests := make([]*cat, 0)
	tests = append(tests, t1)
}
