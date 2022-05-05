package jsoncall

import (
	"encoding/json"
	"fmt"
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
	data, code, err := c.Call("SetA", i, []byte(`{"P1": 7}`))
	if err != nil || code != ErrNone {
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
	data, _, err := c.Call("GetA", i, []byte(`{}`))
	if err != nil {
		t.Errorf("error %v", err)
	}
	s := string(data)
	if s != `{"P1":13}` {
		fmt.Printf("get: %s\n", s)
		t.Fail()
	}
}

type rDiv struct {
	P1 int32
}

func TestJsonCallDiv(t *testing.T) {
	c := newCaller()
	TraceData = true
	i := &TestImpl{}
	data, err := Marshal(3, 2)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	data, code, err := c.Call("Div", i, data)
	if err != nil || code != ErrNone {
		t.Fatalf("call error code: %v, err: %v", code, err)
	}
	m := c.Methods["Div"]
	if m == nil {
		t.Fatalf("method not found")
	}
	fmt.Printf("%v\n", string(data))
	var out rDiv
	err = json.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.P1 != 1 {
		t.Fatalf("incorrect int result")
	}
}

func TestJsonCallError(t *testing.T) {
	c := newCaller()
	i := &TestImpl{}
	data, _, err := c.Call("Div", i, []byte(`{"P1":3,"P2":0}`))
	if err != nil {
		t.Errorf("error %v", err)
	}
	s := string(data)
	if s != `{"P1":0,"P2":{"error":"division by 0"}}` {
		fmt.Println(s)
		t.Fail()
	}
}
