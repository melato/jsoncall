package jsoncall

import (
	"fmt"
	"reflect"
	"testing"
)

type TestInterface interface {
	SetA(int32)
	GetA() int32
	Div(a, b int32) (int32, error)
}

type TestImpl struct {
	A int32
}

func (t *TestImpl) SetA(i int32) {
	fmt.Printf("SetA(%d)\n", i)
	t.A = i
}

func (t *TestImpl) GetA() int32 {
	fmt.Printf("GetA()=%d\n", t.A)
	return t.A
}

func (t *TestImpl) Div(a, b int32) (int32, error) {
	fmt.Printf("Div(%d,%d)\n", a, b)
	if b == 0 {
		return 0, fmt.Errorf("division by 0")
	}
	return a / b, nil
}

func newCaller() *Caller {
	TraceCalls = false
	TraceInit = false
	TraceData = false
	var api *TestImpl
	c, _ := NewJsonCallerP(api)
	return c
}

func TestJsonCallSet(t *testing.T) {
	c := newCaller()
	i := &TestImpl{}
	data, err := c.Call("SetA", i, []byte(`{"P1": 7}`))
	if err != nil {
		t.Errorf("error %v", err)
	}
	if i.A != 7 {
		t.Errorf("expected 7")
	}
	fmt.Printf("%v\n", string(data))
}

func TestJsonCallGet(t *testing.T) {
	c := newCaller()
	i := &TestImpl{}
	i.A = 13
	data, err := c.Call("GetA", i, []byte(`{}`))
	if err != nil {
		t.Errorf("error %v", err)
	}
	s := string(data)
	if s != `{"P1":13}` {
		fmt.Printf("get: %s\n", s)
		t.Fail()
	}
}

func TestJsonCallDiv(t *testing.T) {
	c := newCaller()
	TraceData = true
	i := &TestImpl{}
	data, err := Marshal(3, 2)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	data, err = c.Call("Div", i, data)
	if err != nil {
		t.Fatalf("call error: %v", err)
	}
	m := c.Methods["Div"]
	if m == nil {
		t.Fatalf("method not found")
	}
	n := m.OutType.NumField()
	for i := 0; i < n; i++ {
		fmt.Printf("field[%d] type: %v\n", i, m.OutType.Field(i).Type)
	}
	fmt.Printf("%v\n", string(data))
	out, err := m.unmarshalOutputs(data)
	result := out[0].(int32)
	if result != 1 {
		t.Fatalf("incorrect int result")
	}
	var divError error
	if out[1] != nil {
		divError = out[1].(error)
	}
	fmt.Printf("error type: %v\n", reflect.TypeOf(divError))
	if divError != nil {
		t.Fatalf("error is not nil")
	}
}

func TestJsonCallError(t *testing.T) {
	c := newCaller()
	i := &TestImpl{}
	data, err := c.Call("Div", i, []byte(`{"P1":3,"P2":0}`))
	if err != nil {
		t.Errorf("error %v", err)
	}
	s := string(data)
	if s != `{"P1":0,"P2":{"error":"division by 0"}}` {
		fmt.Println(s)
		t.Fail()
	}
}
