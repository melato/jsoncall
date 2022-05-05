package jsoncall

import (
	"encoding/json"
	"fmt"
	"testing"
)

type A struct {
	B string
	C *Error
}

func TestMarshalError(t *testing.T) {
	a := &A{B: "b", C: &Error{"c"}}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%s\n", string(data))
	a = nil
	err = json.Unmarshal(data, &a)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("b=%s c=%v\n", a.B, a.C)
}

func TestUnmarshalNilError(t *testing.T) {
	var a A
	err := json.Unmarshal([]byte(`{"B":"x","C":null}`), &a)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("b=%s c=%v\n", a.B, a.C)
	if a.C != nil {
		t.Fail()
	}
}

func TestCastError(t *testing.T) {
	err := &Error{"x"}
	var v interface{}
	v = err
	_ = v.(error)
}
