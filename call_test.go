// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

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
	c, _ := NewCaller(api, nil)
	return c
}

func TestJsonCallSet(t *testing.T) {
	c := newCaller()
	i := &TestImpl{}
	m := c.MethodByName("SetA")
	if m == nil {
		t.Fail()
		return
	}
	data, code, err := m.Call(i, []byte(`{"P1": 7}`))
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
	m := c.MethodByName("GetA")
	if m == nil {
		t.Fail()
		return
	}
	data, _, err := m.Call(i, []byte(`{}`))
	if err != nil {
		t.Errorf("error %v", err)
	}
	s := string(data)
	if s != `{"result":13}` {
		fmt.Printf("get: %s\n", s)
		t.Fail()
	}
}

func TestJsonCallDiv(t *testing.T) {
	c := newCaller()
	TraceData = true
	i := &TestImpl{}
	m := c.MethodByName("Div")
	if m == nil {
		t.Fatalf("method not found")
	}
	data, err := m.MarshalInputs(3, 2)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	data, code, err := m.Call(i, data)
	if err != nil || code != ErrNone {
		t.Fatalf("call error code: %v, err: %v", code, err)
	}
	fmt.Printf("%v\n", string(data))
	var out map[string]int
	err = json.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	result := out["result"]
	if result != 1 {
		t.Fatalf("incorrect result: %v", result)
	}
}

func TestJsonCallError(t *testing.T) {
	c := newCaller()
	i := &TestImpl{}
	m := c.MethodByName("Div")
	if m == nil {
		t.Fatalf("method not found")
	}
	data, _, err := m.Call(i, []byte(`{"P1":3,"P2":0}`))
	if err != nil {
		t.Errorf("should not return error")
		return
	}
	if data == nil {
		t.Errorf("should return data")
		return
	}
	s := string(data)
	if s != `{"error":"division by 0","result":0}` {
		t.Errorf("%s", s)
	}
}
