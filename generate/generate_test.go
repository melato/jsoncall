package generate

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
}

func TestTypeName(t *testing.T) {
	var a []A
	s := TypeName(reflect.TypeOf(a), "generate")
	if s != "[]A" {
		fmt.Println(s)
		t.Fail()
	}
	s = TypeName(reflect.TypeOf(a), "x")
	if s != "[]generate.A" {
		fmt.Println(s)
		t.Fail()
	}
}
